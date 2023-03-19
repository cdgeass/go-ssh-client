package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/user"
)

type Config struct {
	Servers []Server `yaml:"servers"`
}

// Server server info
type Server struct {
	Name string `yaml:"name"`

	Host string `yaml:"host"`

	Port int `yaml:"port"`

	User string `yaml:"user"`

	Password string `yaml:"password"`
}

func configPath() string {
	u, err := user.Current()
	if err != nil {
		log.Fatalln("Failed to load config: ", err)
	}

	return u.HomeDir + "\\.go-ssh-client.yaml"
}

func Load() Config {
	configPath := configPath()

	b, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(configPath)
			if err != nil {
				log.Fatalln("Failed to create config file: ", err)
			}
		}
	}

	config := Config{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Fatalln("Failed to parse config file: ", err)
	}

	return config
}

func (conf Config) Save() error {
	b, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	configPath := configPath()
	err = os.WriteFile(configPath, b, 0666)
	if err != nil {
		return err
	}
	return nil
}
