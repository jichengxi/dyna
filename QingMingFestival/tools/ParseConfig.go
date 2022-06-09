package tools

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	AppName  string   `json:"app_name"`
	AppMode  string   `json:"app_mode"`
	AppHost  string   `json:"app_host"`
	AppPort  string   `json:"app_port"`
	Database Database `json:"database"`
}

type Database struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Charset  string `json:"charset"`
	Debug    bool   `json:"debug"`
}

var _cfg *Config

func GetConfig() *Config {
	return _cfg
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&_cfg); err != nil {
		return nil, err
	}
	return _cfg, nil

}

func JsonDecode(io io.ReadCloser, v interface{}) error {
	return json.NewDecoder(io).Decode(v)
}
