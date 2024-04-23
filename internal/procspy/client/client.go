package client

import (
	"crypto/md5"
	"encoding/hex"
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
}

func NewSpy(config *config.Client) *Spy {
	ret := &Spy{
		Config:     config,
		enabled:    false,
		currentDay: time.Now().Day(),
		targets:    domain.NewTargetList(),
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

	//log.Printf("[client] %d Response: %s", res.StatusCode, body)

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

	//log.Printf("[client] %d Response: %s", res.StatusCode, body)

	return string(body), res.StatusCode, nil
}

func (s *Spy) getTargets() ([]*domain.Target, error) {
	targetUrl := fmt.Sprintf("%s/targets/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpGet(targetUrl)

	if err != nil {
		log.Fatalf("[GetTargets] Error getting targets, http status code: %d from %s -> error: %s", status, targetUrl, err)
		return nil, err
	}

	if status != http.StatusOK {
		log.Fatalf("[GetTargets] Error getting targets, http status code: %d from %s", status, targetUrl)
		return nil, fmt.Errorf("http get targets error, http status code: %d", status)
	}

	targets, err := domain.TargetListFromJson(data)

	if err != nil {
		log.Fatalf("[GetTargets] Error getting targets: %s -> bad format", err)
		return nil, err
	}

	if targets == nil {
		log.Fatalf("[GetTargets] Error getting targets: nil")
		return nil, fmt.Errorf("nil targets")
	}

	if len(targets.Targets) == 0 {
		log.Printf("[GetTargets] No targets found")
		return targets.Targets, nil
	}

	hash := getMD5(s.targets)
	newHash := getMD5(targets)

	if hash != newHash {
		log.Printf("[GetTargets] Targets changed, getting new -> %s", targets.ToLog())
		s.targets = targets
	}

	return s.targets.Targets, nil
}

func (s *Spy) getMatches() (map[string]float64, error) {
	matchUrl := fmt.Sprintf("%s/match/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpGet(matchUrl)

	if err != nil {
		log.Fatalf("[GetMatches] Error getting matches, http status code: %d from %s -> error: %s", status, matchUrl, err)
		return nil, err
	}

	if status != http.StatusOK {
		log.Fatalf("[GetMatches] Error getting matches, http status code: %d from %s", status, matchUrl)
		return nil, fmt.Errorf("http get matches error, http status code: %d", status)
	}

	matches, err := domain.MatchListFromJson(data)

	if err != nil {
		log.Fatalf("[GetMatches] Error getting matches: %s -> bad format", err)
		return nil, err
	}

	if matches == nil {
		log.Fatalf("[GetMatches] Error getting matches: nil")
		return nil, fmt.Errorf("nil matches")
	}

	if len(matches.Matches) == 0 {
		log.Printf("[GetMatches] No matches found")
		return matches.Matches, nil
	}

	//log.Printf("[GetMatches] Matches: %s", matches.ToLog())

	return matches.Matches, nil
}

func (s *Spy) postMatch(match *domain.Match) error {
	matchUrl := fmt.Sprintf("%s/match/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpPost(matchUrl, match.ToJson())

	if err != nil {
		log.Fatalf("[PostMatch] Error posting match, http status code: %d to %s -> error: %s", status, matchUrl, err)
		return err
	}

	if status != http.StatusCreated {
		log.Fatalf("[PostMatch] Error posting match, http status code: %d to %s", status, matchUrl)
		return fmt.Errorf("http post match error, http status code: %d", status)
	}

	log.Printf("[PostMatch] Match posted: %s", data)

	return nil
}

func (s *Spy) postCommand(cmd *domain.Command) error {
	commandUrl := fmt.Sprintf("%s/command/%s", s.Config.ServerURL, s.Config.User)

	data, status, err := s.httpPost(commandUrl, cmd.ToJson())

	if err != nil {
		log.Fatalf("[PostCommand] Error posting command, http status code: %d to %s -> error: %s", status, commandUrl, err)
		return err
	}

	if status != http.StatusCreated {
		log.Fatalf("[PostCommand] Error posting command, http status code: %d to %s", status, commandUrl)
		return fmt.Errorf("http post command error, http status code: %d", status)
	}

	log.Printf("[PostCommand] Command posted: %s", data)

	return nil
}

func (s *Spy) run(last time.Time) error {
	log.Printf("Running spy...")

	elapsed := roundFloat(time.Since(last).Seconds(), 2)

	targets, err := s.getTargets()

	if err != nil {
		log.Printf("Error getting targets: %s", err)
		return err
	}

	matches, err := s.getMatches()

	if err != nil {
		log.Printf("Error getting matches: %s", err)
		return err
	}

	processes, err := ps.Processes()
	if err != nil {
		log.Printf("Error getting processes: %s", err)
		return err
	}

	for _, target := range targets {
		if targetElapsed, found := matches[target.Name]; found {
			target.AddElapsed(targetElapsed)
			log.Printf(" > [%s] Use %.2f from %.2fs", target.Name, target.Elapsed, target.Limit)
		}

		match := false
		pids := make([]int, 0)
		names := make([]string, 0)

		for _, proc := range processes {
			name := proc.Executable()

			if target.Match(name) {
				pid := proc.Pid()
				match = true
				pids = append(pids, pid)
				names = append(names, name)
			}
		}

		if len(target.CheckCommand) > 0 {
			cmdLog, err := executeCommand(target.CheckCommand)

			if err != nil {
				log.Printf(" > [%s] Error executing check command [%s]: %s -> %s", target.Name, target.CheckCommand, err, cmdLog)
			} else {
				log.Printf(" > [%s] Check command [%s] -> %s", target.Name, target.CheckCommand, cmdLog)
			}
		}

		if match {
			log.Printf(" > [%s] Found %d processes: %v", target.Name, len(pids), pids)
			matches := strings.Join(names, "//")

			log.Printf(" > [%s] Match process with pattern %s (%s) -> %v", target.Name, target.Pattern, matches, pids)
			err = s.postMatch(domain.NewMatch(s.Config.User, target.Name, target.Pattern, matches, elapsed))

			if err != nil {
				log.Printf(" [%s] Error inserting process: %s", target.Name, err)
			}

			target.AddElapsed(elapsed)
			log.Printf(" > [%s] Add %.2fs -> Use %.2f from %.2fs", target.Name, elapsed, target.Elapsed, target.Limit)

			if target.CheckLimit() {
				log.Printf(" >> [%s] Exceeded limit of %.2f seconds", target.Name, target.Limit)

				if len(target.LimitCommand) > 0 {
					cmdLog, err := executeCommand(target.LimitCommand)

					if err != nil {
						log.Printf(" >> [%s] Error executing limit command [%s]: %s -> %s", target.Name, target.LimitCommand, err, cmdLog)
					} else {
						log.Printf(" >> [%s] Limit command [%s] -> %s", target.Name, target.LimitCommand, cmdLog)
					}
				}

				if target.Kill {
					log.Printf(" >> [%s] Killing processes: %v", target.Name, pids)
					s.kill(pids)
					log.Printf(" >> [%s] %d processes terminated", target.Name, len(pids))
				}
			} else {
				if target.CheckWarning() {
					log.Printf(" >> [%s] Warning on %.2f seconds", target.Name, target.WarningOn)

					if len(target.WarningCommand) > 0 {
						cmdLog, err := executeCommand(target.WarningCommand)

						if err != nil {
							log.Printf(" >> [%s] Error executing warning command [%s]: %s -> %s", target.Name, target.WarningCommand, err, cmdLog)
						} else {
							log.Printf(" >> [%s] Warning command [%s] -> %s", target.Name, target.WarningCommand, cmdLog)
						}
					}
				}
			}
		}
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

func (s *Spy) Start() {
	last := time.Now()
	s.enabled = true

	log.Printf("Starting with config %s", s.Config.ToJson())

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

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)

	err := cmd.Run()

	if err != nil {
		log.Printf("Error executing command: %s -> %s", command, err)
	}

	buf, err := cmd.Output()

	if err != nil {
		log.Printf("Error to read command output: %s -> %s", command, err)
	}

	return string(buf), err

}

func getMD5(t *domain.TargetList) string {
	text := t.ToLog()
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
