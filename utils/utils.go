package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type JSON map[string]any

type GoSilverDB struct {
	path string
}

func GetGosilverState(dbPath string) (JSON, error) {
	f, err := os.OpenFile(dbPath, os.O_RDWR, 0644)

	if err != nil {
		log.Fatal("Error while opening file")
		return nil, err
	}

	defer f.Close()

	rawState, err := io.ReadAll(f)

	if err != nil {
		log.Fatal("Error while reading the content of file")
		return nil, err
	}

	var parsedState map[string]any
	err = json.Unmarshal([]byte(string(rawState)), &parsedState)
	if err != nil {
		log.Fatal("Error while parsing file content to JSON", err)
		return nil, err
	}

	return parsedState, nil
}
