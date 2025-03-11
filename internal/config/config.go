package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LoginNumberAttempts    int `json:"login_number_attempts"`
	PasswordNumberAttempts int `json:"password_number_attempts"`
	IPNumberAttempts       int `json:"ip_number_attempts"`
	GRPCPopt               int `json:"grpc_port"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, using default config")
	}

	path := os.Getenv("CONFIG_PATH")

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var conf Config

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Println(err)
	}

	return &conf
}
