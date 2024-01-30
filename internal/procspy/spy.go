package procspy

import (
	"log"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/mitchellh/go-ps"
	_ "github.com/mitchellh/go-ps"
)

type Spy struct {
	Config     *Config
	enabled    bool
	currentDay int
}

func NewSpy(config *Config) *Spy {
	return &Spy{
		Config:     config,
		enabled:    false,
		currentDay: time.Now().Day(),
	}
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
