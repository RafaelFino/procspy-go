package procspy

import (
	"log"
	"time"

	"github.com/mitchellh/go-ps"
	_ "github.com/mitchellh/go-ps"
)

type Spy struct {
	Config  *Config
	enabled bool
}

func NewSpy(configFile string) *Spy {
	config := NewConfig()
	config.LoadFromFile(configFile)

	return &Spy{
		Config:  config,
		enabled: false,
	}
}

func (s *Spy) run(last time.Time) {
	storage := NewStorage(s.Config.Database)
	storage.Connect()
	storage.CreateTables()
	defer storage.Disconnect()

	for _, proc := range ps.Processes() {
		name := proc.Executable()

		if target, found := s.Config.Targets[name]; found {
			elapsed := time.Since(last).Seconds()

			target.AddElapsed(elapsed)
			storage.InsertProcess(name, elapsed)
			log.Printf("Process %s elapsed %f seconds", name, elapsed)

			if target.IsExpired() {
				log.Printf("Process %s exceeded limit of %f seconds", name, target.GetLimit())
				proc.Kill()
				log.Printf("Process %s killed", name)
			}
		}
	}
}

func (s *Spy) Start() {
	s.loadFromDatabase()

	last := time.Now()

	s.enabled = true

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
	storage := NewStorage(s.Config.Database)
	storage.Connect()
	defer storage.Disconnect()

	elapsed, err := storage.GetElapsed()
	if err != nil {
		log.Println(err)
		return
	}

	for name, limit := range s.Config.Targets {
		if elapsed, found := elapsed[name]; found {
			limit.AddElapsed(elapsed)
		}
	}
}
