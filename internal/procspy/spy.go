package procspy

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/mitchellh/go-ps"
	_ "github.com/mitchellh/go-ps"
)

type Spy struct {
	Config  *Config
	enabled bool
}

func NewSpy(config *Config) *Spy {
	fmt.Print("Starting spy...\n")

	return &Spy{
		Config:  config,
		enabled: false,
	}
}

func (s *Spy) run(last time.Time) error {
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

	targets := make(map[string][]int)
	elapsed := roundFloat(time.Since(last).Seconds(), 2)

	for _, proc := range processes {
		name := proc.Executable()

		if _, found := s.Config.Targets[name]; found {
			if _, found := targets[name]; !found {
				targets[name] = make([]int, 0)
			}
			pids := targets[name]
			pids = append(pids, proc.Pid())
			targets[name] = pids
		}
	}

	for name, pids := range targets {
		if target, found := s.Config.Targets[name]; found {
			target.AddElapsed(elapsed)
			s.Config.Targets[name] = target
			err = storage.InsertProcess(name, elapsed)
			if err != nil {
				log.Printf("Error inserting process %s: %s", name, err)
			}

			log.Printf(" > [%s] Add %.2fs -> Use %.2f from %.2fs", name, elapsed, target.GetElapsed(), target.GetLimit())

			if target.IsExpired() {
				log.Printf(" >> Process %s exceeded limit of %.2f seconds", name, target.GetLimit())
				s.kill(pids)
				log.Printf(" >> Killed %d processes from %s", len(pids), name)
			}
		}
	}

	return err
}

func (s *Spy) kill(pids []int) {
	if pids == nil || len(pids) == 0 {
		return
	}

	log.Printf(" >> Killing processes: %v", pids)

	for _, pid := range pids {
		p, err := os.FindProcess(pid)
		if err != nil {
			log.Printf("Error finding process %d: %s", pid, err)
		} else {
			err = p.Kill()
			if err != nil {
				log.Printf("Error killing process %d: %s", pid, err)
			}
		}
	}
}

func (s *Spy) Start() {
	s.loadFromDatabase()

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

	for name, limit := range s.Config.Targets {
		if elapsed, found := elapsed[name]; found && elapsed > 0 {
			limit.AddElapsed(elapsed)
			s.Config.Targets[name] = limit
		}
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
