package grapho

import (
	"log"
	"os"
)

type Config struct {
	Storage func() EventStorage
}

func getenv(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	} else {
		return value
	}
}

var (
	Configurations = map[string]*Config{
		"test": &Config{
			Storage: func() EventStorage {
				storage, _ := NewEventsInMemory()
				return storage
			},
		},
		"development": &Config{
			Storage: func() EventStorage {
				storage, err := NewEventsOnDisk("grapho.log")
				if err != nil {
					log.Fatal(err)
				}

				return storage
			},
		},
	}
)
