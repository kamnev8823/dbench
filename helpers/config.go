package helpers

import (
	"errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

//name config name
const name = "config.yaml"

//absPath get absolute file path
func absPath() string {
	absPath, err := filepath.Abs(name)

	if err != nil {
		log.Fatal(err)
	}
	return absPath
}

//CheckConfig check file configuration
func CheckConfig() {
	absPath := absPath()
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatal("config.yaml is missing!")
	}
}

//GetValue get value from file configuration
func GetValue(key string) (interface{}, error) {
	config := all()

	value, ok := config[key]

	if !ok {
		return "", errors.New("Key not found! ")
	}

	return value, nil
}

//all get all config
func all() map[interface{}]interface{} {
	path := absPath()

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(file, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return m
}
