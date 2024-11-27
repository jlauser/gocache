package db

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CsvDB struct {
	Tables   map[string][][]string
	Mappings map[string]string
	Ctx      context.Context
}

func InitializeCsvDB(value string) (*CsvDB, error) {
	// Format path if needed
	path := value
	if !strings.HasSuffix(path, "\\") {
		path = path + "\\"
	}

	// Initialize memory objects
	db := CsvDB{
		Tables:   make(map[string][][]string),
		Mappings: make(map[string]string),
	}

	// Verify and read data directory
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("data file directory not found. %w", err)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory. %w", err)
	}

	// Iterate over directory entries
	for _, file := range files {
		if !file.IsDir() {
			index := strings.Index(file.Name(), "_")
			if index > -1 {
				prefix := file.Name()[:index]
				switch prefix {
				case "db":
					// Create tables and load data
					table := strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".csv"), "db_")
					filePath := filepath.Join(path, file.Name())
					db.Mappings[table] = filePath
					rows, err := loadFile(filePath)
					if err != nil {
						return nil, fmt.Errorf("failed to read file %s. %w", filePath, err)
					}
					db.Tables[table] = rows
				}
			}
		}
	}
	log.Printf("CSV DB created with %d tables", len(db.Tables))
	return &db, nil
}

// Load all csv files in the data directory
func loadFile(filename string) ([][]string, error) {
	result := make([][]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", filename, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)
	rowId := 1
	for {
		// Read a single record (row)
		row, err := reader.Read()
		if err != nil {
			// Break the loop if EOF (end of file) is reached
			if err.Error() == "EOF" {
				break
			}
		}
		record := []string{strconv.Itoa(rowId)}
		record = append(record, row...)
		result = append(result, record)
		rowId++
	}
	return result, nil
}

func (db *CsvDB) appendToFile(table string, data []string) bool {
	filename := db.Mappings[table]
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("could not open file %s: %v", filename, err)
		return false
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush() // Flush ensures all data is written

	// Write the new row to the file. omit the index
	err = writer.Write(data[1:])
	if err != nil {
		log.Printf("error writing to CSV %w", err)
		return false
	}
	return true
}

func (db *CsvDB) getRowFromMap(table string, index string) (interface{}, bool) {
	// Check if the key exists in the map
	rows, ok := db.Tables[table]
	if !ok {
		return nil, false
	}

	// Search for a row where the first value matches valueToMatch
	for _, row := range rows {
		if len(row) > 0 && row[0] == index {
			return row, true
		}
	}
	return nil, false
}

func (db *CsvDB) findInTable(table string, search string) (interface{}, bool) {
	result := make([][]string, 0)
	foundRows := false
	// Check if the key exists in the map
	rows, ok := db.Tables[table]
	if !ok {
		return nil, false
	}

	// Search for a row where the first value matches valueToMatch
	for _, row := range rows {
		text := strings.ToLower(strings.Join(row, " "))
		if search == "*" || strings.Contains(text, strings.ToLower(search)) {
			result = append(result, row)
			foundRows = true
		}
	}
	return result, foundRows
}

func (db *CsvDB) overwriteFile(table string) bool {
	filename := db.Mappings[table]
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("could not open file %s: %v", filename, err)
		return false
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ','

	for _, row := range db.Tables[table] {
		if row[0] != "DEL" {
			err := writer.Write(row[1:])
			if err != nil {
				log.Printf("error writing to CSV %s. %w", filename, err)
				return false
			}
		}
	}
	return true
}

// Create an item to the CsvDB
func (db *CsvDB) Create(key string, value interface{}) (string, bool) {
	keyParts := strings.SplitN(key, ":", 2)
	table := keyParts[0]
	key = keyParts[1]
	row, ok := value.([]string)
	if ok {
		if key == "" {
			key = strconv.Itoa(len(db.Tables[table]) + 1)
		}
		record := []string{key}
		record = append(record, row...)
		db.Tables[table] = append(db.Tables[table], record)
		return key, db.appendToFile(table, record)
	}
	return "", false
}

// Get an item from the CsvDB
func (db *CsvDB) Read(key string) (interface{}, bool) {
	keyParts := strings.SplitN(key, ":", 2)
	table := keyParts[0]
	key = keyParts[1]
	rows, ok := db.Tables[table]
	if !ok {
		return nil, false
	}
	for _, row := range rows {
		if len(row) > 0 && row[0] == key {
			return row, true
		}
	}
	return nil, false
}

// Find an item from the CsvDB
func (db *CsvDB) Find(key string, value interface{}) (interface{}, bool) {
	search, ok := value.(string)
	if !ok {
		return nil, false
	}
	return db.findInTable(key, search)
}

// Update an item in the CsvDB
func (db *CsvDB) Update(key string, value interface{}) bool {
	keyParts := strings.SplitN(key, ":", 2)
	table := keyParts[0]
	key = keyParts[1]
	data, ok := value.([]string)
	if !ok {
		return false
	}
	rows, ok := db.Tables[table]
	if !ok {
		return false
	}
	// update row
	for _, row := range rows {
		if len(row) > 0 && row[0] == key {
			row = data
			break
		}
	}
	// save to file
	return db.overwriteFile(table)
}

// Delete and item from the CsvDB
func (db *CsvDB) Delete(key string) bool {
	keyParts := strings.SplitN(key, ":", 2)
	table := keyParts[0]
	key = keyParts[1]
	rows, ok := db.Tables[table]
	if !ok {
		return false
	}
	// mark for deletion
	for _, row := range rows {
		if len(row) > 0 && row[0] == key {
			row[0] = "DEL"
		}
	}
	// save to file
	return db.overwriteFile(table)
}
