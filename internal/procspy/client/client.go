package client

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/domain"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	_ "github.com/mitchellh/go-ps"
)

type Spy struct {
	Config     *config.Client
	enabled    bool
	currentDay int
	targets    *domain.TargetList
	commandBuf chan *domain.Command
	matchBuf   chan *domain.Match
}

func NewSpy(config *config.Client) *Spy {
	ret := &Spy{
		Config:     config,
		enabled:    false,
		currentDay: time.Now().Day(),
		targets:    domain.NewTargetList(),
		commandBuf: make(chan *domain.Command, 1000),
		matchBuf:   make(chan *domain.Match, 1000),
	}
	/*
		s.router.GET("/targets/:user", s.targetHandler.GetTargets)
		s.router.POST("/match/:user", s.matchHandler.InsertMatch)
		s.router.GET("/match/:user", s.matchHandler.GetMatches)
		s.router.POST("/command/:user/:name", s.commandHandler.InsertCommand)
	*/
	return ret
}

func (s *Spy) httpGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[client] Error getting url: %s", err)
		return "", http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[client] Error reading body: %s", err)
		return "", res.StatusCode, err
	}

	if s.Config.Debug {
		log.Printf("[HTTP-GET] %d Response: %s", res.StatusCode, body)
	}

	return string(body), res.StatusCode, nil
}

func (s *Spy) httpPost(url string, data string) (string, int, error) {
	res, err := http.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		log.Printf("[client] Error posting url: %s", err)
		return "", http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[client] Error reading body: %s", err)
		return "", res.StatusCode, err
	}

	if s.Config.Debug {
		log.Printf("[HTTP-POST] %d \nRequest: %s\nResponse: %s", res.StatusCode, data, body)
	}

	return string(body), res.StatusCode, nil
}

func (s *Spy) updateTargets() {
	if s.targets == nil {
		s.targets = domain.NewTargetList()
	}

	targetUrl := fmt.Sprintf("%s/targets/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpGet(targetUrl)

	if err != nil {
		log.Printf("[UpdateTargets] Error getting targets, http status code: %d from %s -> error: %s", status, targetUrl, err)
		return
	}

	if status != http.StatusOK {
		log.Printf("[UpdateTargets] Error getting targets, http status code: %d from %s", status, targetUrl)
		return
	}

	targets, err := domain.TargetListFromJson(data)

	if err != nil {
		log.Printf("[UpdateTargets] Error getting targets: %s -> bad format", err)
		return
	}

	if targets == nil {
		log.Printf("[UpdateTargets] Error getting targets: nil")
		return
	}

	if len(targets.Targets) == 0 {
		log.Printf("[UpdateTargets] No targets found")
	}

	s.targets = targets
}

func (s *Spy) postMatch(match *domain.Match) error {
	if match == nil {
		return fmt.Errorf("match is nil")
	}

	matchUrl := fmt.Sprintf("%s/match/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpPost(matchUrl, match.ToJson())

	if err != nil {
		log.Printf("[PostMatch] Error posting match, http status code: %d to %s -> error: %s", status, matchUrl, err)
		return err
	}

	if status != http.StatusCreated {
		log.Printf("[PostMatch] Error posting match, http status code: %d to %s", status, matchUrl)
		return fmt.Errorf("http post match error, http status code: %d", status)
	}

	if s.Config.Debug {
		log.Printf("[PostMatch] Match POST return: %s\n from \n%s", data, match.ToJson())
	}

	return nil
}

func (s *Spy) postCommand(cmd *domain.Command) error {
	if cmd == nil {
		return fmt.Errorf("command is nil")
	}

	commandUrl := fmt.Sprintf("%s/command/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpPost(commandUrl, cmd.ToJson())

	if err != nil {
		log.Printf("[PostCommand] Error posting command, http status code: %d to %s -> error: %s", status, commandUrl, err)
		return err
	}

	if status != http.StatusCreated {
		log.Printf("[PostCommand] Error posting command, http status code: %d to %s", status, commandUrl)
		return fmt.Errorf("http post command error, http status code: %d", status)
	}

	if s.Config.Debug {
		log.Printf("[PostCommand] Command POST return: %s\nfrom \n%s", data, cmd.ToJson())
	}

	return nil
}

func (s *Spy) consumeBuffers() {
	if s.matchBuf == nil {
		log.Printf("[Spy] Match buffer is nil")
		s.matchBuf = make(chan *domain.Match, 1000)
	}

	if s.commandBuf == nil {
		log.Printf("[Spy] Command buffer is nil")
		s.commandBuf = make(chan *domain.Command, 1000)
	}

	go func() {
		//Buffer
		matchDlq := make(chan *domain.Match, len(s.matchBuf))
		for len(s.matchBuf) > 0 {
			if s.Config.Debug {
				log.Printf("[Spy] %d matches in buffer", len(s.matchBuf))
			}

			match := <-s.matchBuf
			err := s.postMatch(match)
			if err != nil {
				log.Printf("[Spy] Error posting match: %s, waiting for next process", err)
				matchDlq <- match
			}
		}

		cmdDlq := make(chan *domain.Command, len(s.commandBuf))
		for len(s.commandBuf) > 0 {
			if s.Config.Debug {
				log.Printf("[Spy] %d commands in buffer", len(s.commandBuf))
			}

			cmd := <-s.commandBuf
			err := s.postCommand(cmd)
			if err != nil {
				log.Printf("[Spy] Error posting command: %s, waiting for next process", err)
				cmdDlq <- cmd
			}
		}

		//DLQ
		for len(matchDlq) > 0 {
			match := <-matchDlq
			log.Printf("[Spy] Add match to post dlq: %s", match.ToJson())
			s.matchBuf <- match
		}

		for len(cmdDlq) > 0 {
			cmd := <-cmdDlq
			log.Printf("[Spy] Add command to post dlq: %s", cmd.ToJson())
			s.commandBuf <- cmd
		}
	}()
}

func (s *Spy) run(last time.Time) error {
	log.Printf("[Spy] Starting...")
	defer log.Printf("[Spy] Finished")

	elapsed := roundFloat(time.Since(last).Seconds(), 2)

	defer s.consumeBuffers()
	s.updateTargets()

	processes, err := ps.Processes()
	if err != nil {
		log.Printf("[Spy] Error getting processes: %s", err)
		return err
	}

	for _, target := range s.targets.Targets {
		match := false
		pids := make([]int, 0)
		names := make(map[string]struct{})

		for _, proc := range processes {
			name := proc.Executable()

			if target.Match(name) {
				pid := proc.Pid()
				match = true
				pids = append(pids, pid)
				names[name] = struct{}{}
			}
		}

		if len(target.CheckCommand) > 0 {
			log.Printf("[Spy]  > [%s] Use %.2f from %.2fs", target.Name, target.Elapsed, target.Limit)
			cmdLog, err := executeCommand(target.CheckCommand)

			if err != nil {
				log.Printf("[Spy]  > [%s] Error executing check command [%s]: %s -> %s", target.Name, target.CheckCommand, err, cmdLog)
			} else {
				log.Printf("[Spy]  > [%s] Check command [%s] -> %s", target.Name, target.CheckCommand, cmdLog)
			}

			cmd := domain.NewCommand(s.Config.User, target.Name, target.LimitCommand, cmdLog)
			cmd.Source = "Check"
			s.commandBuf <- cmd
		}

		if match {
			log.Printf("[Spy]  > [%s] Found %d processes: %v", target.Name, len(pids), pids)

			matches := make([]string, 0)
			for k := range names {
				matches = append(matches, k)
			}

			log.Printf("[Spy]  > [%s] Match process with pattern %s (%s) -> %v", target.Name, target.Pattern, matches, pids)
			s.matchBuf <- domain.NewMatch(s.Config.User, target.Name, target.Pattern, strings.Join(matches, " / "), elapsed)

			target.AddElapsed(elapsed)
			log.Printf("[Spy]  > [%s] Add %.2fs -> Use %.2f from %.2fs", target.Name, elapsed, target.Elapsed, target.Limit)

			if target.CheckLimit() {
				log.Printf("[Spy]  >> [%s] Exceeded limit of %.2f seconds", target.Name, target.Limit)

				if len(target.LimitCommand) > 0 {
					cmdLog, err := executeCommand(target.LimitCommand)

					if err != nil {
						log.Printf("[Spy]  >> [%s] Error executing limit command [%s]: %s -> %s", target.Name, target.LimitCommand, err, cmdLog)
					} else {
						log.Printf("[Spy]  >> [%s] Limit command [%s] -> %s", target.Name, target.LimitCommand, cmdLog)
					}

					cmd := domain.NewCommand(s.Config.User, target.Name, target.LimitCommand, cmdLog)
					cmd.Source = "Limit"
					s.commandBuf <- cmd
				}

				if target.Kill {
					log.Printf("[Spy]  >> [%s] Killing processes: %v", target.Name, pids)
					s.kill(target.Name, pids)
					log.Printf("[Spy]  >> [%s] %d processes terminated", target.Name, len(pids))
				}
			} else {
				if target.CheckWarning() {
					log.Printf("[Spy]  >> [%s] Warning on %.2f seconds", target.Name, target.WarningOn)

					if len(target.WarningCommand) > 0 {
						cmdLog, err := executeCommand(target.WarningCommand)

						if err != nil {
							log.Printf("[Spy]  >> [%s] Error executing warning command [%s]: %s -> %s", target.Name, target.WarningCommand, err, cmdLog)
						} else {
							log.Printf("[Spy]  >> [%s] Warning command [%s] -> %s", target.Name, target.WarningCommand, cmdLog)
						}

						cmd := domain.NewCommand(s.Config.User, target.Name, target.WarningCommand, cmdLog)
						cmd.Source = "Warning"
						s.commandBuf <- cmd
					}
				}
			}
		}
	}

	return err
}

func (s *Spy) kill(name string, pids []int) {
	if len(pids) == 0 {
		return
	}

	for _, pid := range pids {
		p, err := os.FindProcess(pid)
		if err != nil {
			log.Printf("[Kill]  >> Process %d not found: %s", pid, err)
		} else {
			err = p.Kill()
			msg := ""
			if err != nil {
				log.Printf("[Kill]  >> Warn: killing process %d: %s", pid, err)
				msg = err.Error()
			}

			cmd := domain.NewCommand(s.Config.User, name, fmt.Sprintf("kill %d", pid), msg)
			cmd.Source = "Kill"
			s.commandBuf <- cmd
		}
	}
}

func (s *Spy) Start() {
	last := time.Now()
	s.enabled = true

	log.Printf("[StartSpy] Starting with config ->\n%s", s.Config.ToJson())

	for s.enabled {
		s.run(last)
		last = time.Now()
		time.Sleep(time.Duration(s.Config.Interval) * time.Second)
	}
}

func (s *Spy) Stop() {
	s.enabled = false
	log.Printf("[StopSpy] Stopping...")
}

func (s *Spy) IsEnabled() bool {
	return s.enabled
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)

	log.Printf("[ExecuteCommand] Executing command: %s", command)
	err := cmd.Run()

	if err != nil {
		log.Printf("[ExecuteCommand] Error executing command: %s -> %s", command, err)
	}

	buf, err := cmd.Output()

	if err != nil {
		log.Printf("[ExecuteCommand] Error to read command output: %s -> %s", command, err)
	}

	return string(buf), err
}
