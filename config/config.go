package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var config Config

type Target struct {
	WorkerId string   `json:"workerId"`
	Threads  []string `json:"threads"`
}

type User struct {
	Email string `json:"email"`
	Alias string `json:"alias"`
}

type Config struct {
	Targets []Target `json:"targets"`
	Users   []User   `json:"users"`
	Matches []string `json:"matches"`
}

func init() {
	LoadConfig()
}

func LoadConfig() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fail to get working directory. Error=%v\n", err)
	}
	if data, err := ioutil.ReadFile(path.Join(dir, "config.json")); err != nil {
		log.Fatalf("Fail to read config file. Error=%v\n", err)
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatalf("Fail to parse config. Error=%v\n", err)
		}
	}
}

func GetConfig() Config {
	return config
}

func FindUser(email string) User {
	for _, user := range config.Users {
		if email == user.Email {
			return user
		}
	}
	return User{
		Email: email,
		Alias: "",
	}
}
