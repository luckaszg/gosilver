package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gosilver/utils"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)


type GoSilverDB struct {
	path string
}

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

	currentState, err := utils.GetGosilverState(gs.path)

	if err != nil {
		fmt.Printf("Error while getting database state\n")
	}
	
	// Check if the table name is already in database state
	if _, ok := currentState[tableName]; ok {
		fmt.Printf("Table '%s' already in Database\n", tableName)
		return
	}

	// If table name is not in the state, proceed to add it
	currentState[tableName] = map[string]any{}
	tableData, err := json.Marshal(currentState)
	if err != nil {
		log.Fatal("Error whle encoding table name to JSON")
	}
	if err = os.WriteFile(gs.path, tableData, 0660); err != nil {
		log.Fatal("Error while trying to insert the table name", err)
	}
}

func (gs GoSilverDB) write(tableName string, fileContent utils.JSON) {
	itemUuid, err := uuid.NewV7();
	if err != nil {
		log.Fatal("Error while generating UUID for table", err)
	}

	fileContent["id"] = itemUuid

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

	table, ok := deserializedContent[tableName].(map[string]any)
	if !ok {
		fmt.Printf("Table '%s' doesn't exist in database\n", tableName)
		return
	}
	table[itemUuid.String()] = fileContent

	// START. Adding new table name to database object
	tableData, err := json.Marshal(deserializedContent)
	if err != nil {
		log.Fatal("Error whle encoding table name to JSON")
	}
	if err = os.WriteFile(gs.path, tableData, 0660); err != nil {
		log.Fatal("Error while trying to insert the table data", err)
	}
	//END. Adding new table name to database object
}

func (gs GoSilverDB) read(tableName string, itemUuid string) (utils.JSON, error) {
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

	table, ok := deserializedContent[tableName].(map[string]any)
	if !ok {
		fmt.Printf("Table '%s' doesn't exist in database\n", tableName)
		return nil, errors.New("Table doesn't exist in database")
	}

	record, ok := table[itemUuid].(map[string]any)
	if !ok {
		fmt.Printf("Record with UUID '%s' doesn't exist in table '%s'", itemUuid, tableName)
		return nil, errors.New("Record doesn't exist in table")
	}
	return record, nil
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
	// db.createTable("client")

	// Get table
	//table = db.getTable("users")

	// Write opperation
	lucasGuerreroData := map[string]any{"name": "Lucas Guerrero", "year_joined": 2026, "title": "software engineer"}
	nelsonRojasData := map[string]any{"name": "Nelson Rojas", "year_joined": 2025, "title": "senior software engineer"}
	db.write("user", lucasGuerreroData)
	db.write("user", nelsonRojasData)

	nvidiaMarket := map[string]any{"name": "NVIDIA Corporation", "symbol": "NVDA", "active": true}
	anthropicMarket := map[string]any{"name": "Anthropic", "symbol": "ANT", "active": true}
	lehmanBrothersMarket := map[string]any{"name": "Lehman Brothers", "symbol": "LEHM", "active": false}
	db.write("market", nvidiaMarket)
	db.write("market", anthropicMarket)
	db.write("market", lehmanBrothersMarket)

	// Read opperation (by field)
	nvidia, err := db.read("market", "01a20fb0-6e11-72a7-b46f-189b7aeec495")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(nvidia["name"])

	lucas, err := db.read("user", "01a20fb0-6e10-7cd7-a93d-900d79001a06")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lucas["name"])

	// fmt.Print(db.path)
}
