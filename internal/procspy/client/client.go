package client

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"procspy/internal/procspy"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/service"
	"time"

	"github.com/mitchellh/go-ps"
	_ "github.com/mitchellh/go-ps"
)

type Spy struct {
	Config     *config.Client
	enabled    bool
	currentDay int
	targets    []*domain.Target
	token      string
	pubKey     string
}

func NewSpy(config *config.Client) *Spy {
	ret := &Spy{
		Config:     config,
		enabled:    false,
		currentDay: time.Now().Day(),
		targets:    make([]*domain.Target, 0),
		token:      "",
	}
	/*
		s.router.GET("/key", s.authHandler.GetPubKey)
		s.router.POST("/user/:user", s.userHandler.CreateUser)
		s.router.POST("/auth/", s.authHandler.Authenticate)

		s.router.GET("/targets/:user", s.targetHandler.GetTargets)
		s.router.POST("/match/:user", s.matchHandler.InsertMatch)
		s.router.GET("/match/:user", s.matchHandler.GetMatches)
		s.router.POST("/command/:user/:name", s.commandHandler.InsertCommand)
	*/
	return ret
}

func (s *Spy) Auth() error {
	keyUrl := fmt.Sprintf("%s/key/", s.Config.ServerURL)
	data, status, err := procspy.HttpGet(keyUrl, "")

	if err != nil {
		log.Fatalf("[Auth] Error getting public key, http status code: %s from %s -> error: %s", status, keyUrl, err)
		return err
	}

	if status != 200 {
		log.Fatalf("[Auth] Error getting public key, http status code: %s from %s", status, keyUrl)
		return fmt.Errorf("http get pub key error, http status code: %s", status)
	}

	pubKey, ok := data["key"]

	if !ok {
		log.Fatalf("[Auth] Error getting public key: %s -> bad format", err)
		return err
	}

	log.Printf("[Auth] Public key: %s", pubKey)
	s.pubKey = pubKey.(string)

	userInfo := fmt.Sprintf(`{ "user": "%s"}`, s.Config.User)

	payload, err := service.Cypher(userInfo, []byte(s.pubKey))

	authUrl := fmt.Sprintf("%s/auth/", s.Config.ServerURL)
	resp, status, err := procspy.HttpPost(authUrl, "", payload)

	if err != nil {
		log.Fatalf("[Auth] Error authenticating user: %s -> %s", s.Config.User, err)
		return err
	}

	if status != 200 {
		log.Fatalf("[Auth] Error authenticating user: %s -> http status code: %s", s.Config.User, status)
		return fmt.Errorf("http post auth error, http status code: %s", status)
	}

	token, ok := resp["token"]

	s.token = token.(string)

	log.Printf("[Auth] Token: %s", s.token)
}

func (s *Spy) GetTargets() error {
	if s.token == "" {
		err := s.Auth()
		if err != nil {
			return err
		}
	}

	targetUrl := fmt.Sprintf("%s/targets/%s", s.Config.ServerURL, s.Config.User)
	data, status, err := procspy.HttpGet(targetUrl, s.token)

	if err != nil {
		log.Fatalf("[GetTargets] Error getting targets, http status code: %s from %s -> error: %s", status, targetUrl, err)
		return err
	}

	if status != 200 {
		log.Fatalf("[GetTargets] Error getting targets, http status code: %s from %s", status, targetUrl)
		return fmt.Errorf("http get targets error, http status code: %s", status)
	}

	targets, ok := data["targets"]

	if !ok {
		log.Fatalf("[GetTargets] Error getting targets: %s -> bad format", err)
		return err
	}

	log.Printf("[GetTargets] Targets: %s", targets)
	s.targets = targets.([]*domain.Target)

	return nil
}

func (s *Spy) run(last time.Time) error {
	// reload from web?
	err := s.Config.UpdateFromUrl()
	if err != nil {
		log.Printf("Error loading config from url: %s", err)
	}

	storage, err := NewStorage(s.Config.DBPath)

	if err != nil {
		log.Printf("Error opening database: %s", err)
		return err
	}

	defer storage.Close()

	processes, err := ps.Processes()
	if err != nil {
		log.Printf("Error getting processes: %s", err)
		return err
	}

	if s.currentDay != time.Now().Day() {
		log.Printf("Resetting elapsed time for all processes, day changed")
		s.currentDay = time.Now().Day()

		//reload from web?
		err = s.Config.UpdateFromUrl()
		if err != nil {
			log.Printf("Error loading config from url: %s", err)
		}

		for index, target := range s.Config.Targets {
			log.Printf(" # [%s] Resetting elapsed time", target.GetName())
			target.ResetElapsed()
			s.Config.Targets[index] = target
		}
	}

	elapsed := roundFloat(time.Since(last).Seconds(), 2)

	if elapsed < float64(s.Config.Interval) {
		return nil
	}

	for index, target := range s.Config.Targets {
		match := false
		pids := make([]int, 0)

		for _, proc := range processes {
			name := proc.Executable()

			if target.Match(name) {
				match = true
				pid := proc.Pid()
				pids = append(pids, pid)
			}
		}

		if match {
			log.Printf(" > [%s] Match process with pattern %s -> %v", target.GetName(), target.GetPattern(), pids)
			target.AddElapsed(elapsed)

			err = storage.InsertProcess(target.GetName(), elapsed, target.GetPattern(), target.GetCommand(), target.GetKill())
			if err != nil {
				log.Printf(" [%s] Error inserting process: %s", target.GetName(), err)
			}

			log.Printf(" > [%s] Add %.2fs -> Use %.2f from %.2fs", target.GetName(), elapsed, target.GetElapsed(), target.GetLimit())

			if target.IsExpired() {
				log.Printf(" >> [%s] Exceeded limit of %.2f seconds", target.GetName(), target.GetLimit())
				if target.GetKill() {
					log.Printf(" >> [%s] Killing processes: %v", target.GetName(), pids)
					s.kill(pids)
					log.Printf(" >> [%s] %d processes terminated", target.GetName(), len(pids))
				}

				err = storage.InsertMatch(target.GetName(), target.GetPattern(), target.GetCommand(), target.GetKill())
				if err != nil {
					log.Printf("[%s] Error inserting match: %s", target.GetName(), err)
				}

				if len(target.GetCommand()) > 0 {
					log.Printf(" >> [%s] Executing command: %s", target.GetName(), target.GetCommand())
					err = executeCommand(target.GetCommand())
					if err != nil {
						log.Printf(" [%s] Error executing command %s: %s", target.GetName(), target.GetCommand(), err)
					}
				}
			}
		}

		s.Config.Targets[index] = target
	}

	return err
}

func (s *Spy) kill(pids []int) {
	if pids == nil || len(pids) == 0 {
		return
	}

	for _, pid := range pids {
		p, err := os.FindProcess(pid)
		if err != nil {
			log.Printf(" >> Process %d not found: %s", pid, err)
		} else {
			err = p.Kill()
			if err != nil {
				log.Printf(" >> Warn: killing process %d: %s", pid, err)
			}
		}
	}
}

func (s *Spy) Start(onUpdate chan bool) {
	s.loadFromDatabase()

	last := time.Now()
	s.enabled = true

	log.Printf("Starting with config %s", s.Config.ToJson())

	go func() {
		for {
			select {
			case <-onUpdate:
				log.Printf("Spy: Config updated, reloading...")
				s.loadFromDatabase()
			}
		}
	}()

	for s.enabled {
		s.run(last)
		last = time.Now()
		time.Sleep(time.Duration(s.Config.Interval) * time.Second)
	}
}

func (s *Spy) Stop() {
	s.enabled = false
}

func (s *Spy) IsEnabled() bool {
	return s.enabled
}

func (s *Spy) loadFromDatabase() {
	storage, err := NewStorage(s.Config.DBPath)

	if err != nil {
		log.Printf("Error opening database: %s", err)
		return
	}

	defer storage.Close()

	elapsed, err := storage.GetElapsed()
	if err != nil {
		log.Printf("Error getting elapsed: %s", err)
		return
	}

	for index, target := range s.Config.Targets {
		if elapsed, found := elapsed[target.GetName()]; found && elapsed > 0 {
			target.ResetElapsed()
			target.AddElapsed(elapsed)
			log.Printf(" > [%s] Loaded %.2fs -> Use %.2f from %.2fs", target.GetName(), elapsed, target.GetElapsed(), target.GetLimit())
			s.Config.Targets[index] = target
		}
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func executeCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}
