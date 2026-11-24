package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)


type GoSilverDB struct {
	path string
}

type JSON map[string]any

func GoSilver(fileName string) (*GoSilverDB, error) {
	if err := os.MkdirAll("db", 0750); err != nil {
		return nil, err
	}
	fullFileName := fmt.Sprintf("%s.%s", fileName, "gs")
	path := filepath.Join("db", fullFileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	// Initialize empty file with empty object
	if err == nil {
		f.WriteString("{}")
		f.Close()
	} else if !os.IsExist(err) {
		return nil, err
	}
	return &GoSilverDB{path: path}, nil
}

func (gs GoSilverDB) createTable(tableName string) {

	f, err := os.OpenFile(gs.path, os.O_RDWR, 0644)

	if err != nil {
		log.Fatal("Error while opening file")
	}

	defer f.Close()

	// START. Check if table name already exists section
	content, err := io.ReadAll(f)

	if err != nil {
		log.Fatal("Error while reading the content of file")
	}

	var deserializedContent map[string]any
	err = json.Unmarshal([]byte(string(content)), &deserializedContent)
	if err != nil {
		log.Fatal("Error while parsing file content to JSON", err)
	}

	if _, ok := deserializedContent[tableName]; ok {
		fmt.Printf("Table '%s' already in Database", tableName)
		fmt.Println()
		return
	}
	// END. Check if table name already exists section

	// START. Adding new table name to database object
	deserializedContent[tableName] = map[string]any{}
	tableData, err := json.Marshal(deserializedContent)
	if err != nil {
		log.Fatal("Error whle encoding table name to JSON")
	}
	if err = os.WriteFile(gs.path, tableData, 0660); err != nil {
		log.Fatal("Error while trying to insert the table name", err)
	}
	//END. Adding new table name to database object
}

func (gs GoSilverDB) write(fileContent JSON) {

}

func main() {
	db, err := GoSilver("userdata")
	if err != nil {
		log.Fatal(err)
	}
	_ = db

	// Create table and return
	db.createTable("user")
	db.createTable("market")
	db.createTable("client")

	// Get table
	//table = db.getTable("users")

	// Write opperation
	//db.write()

	// Read opperation (by field)
	//db.read()

	fmt.Print(db.path)
}
