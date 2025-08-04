package memdb

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/Hanasou/news_feed/go/common"
)

// Driver for data held in memory
type MemDb[T common.Serializable] struct {
	Table      string
	Data       map[string]T
	FilePath   string
	SaveToDisk bool
}

func (db *MemDb[T]) String() string {
	for key, value := range db.Data {
		valueJson, err := value.ToJson()
		if err != nil {
			log.Printf("Error converting value to JSON: %v", err)
			return "MemDb{" + db.Table + ", " + key + ": <error converting to JSON>}"
		}
		return "MemDb{" + db.Table + ", " + key + ": " + valueJson + "}"
	}
	return "MemDb{" + db.Table + "}"
}

func Initialize[T common.Serializable](table string, rootPath string, saveToDisk bool) (*MemDb[T], error) {
	data := map[string]T{}
	filePath := ""
	if saveToDisk {
		var err error
		filePath = rootPath + "/" + table + ".json"
		data, err = GetDataFromFile[T](rootPath + "/" + table + ".json")
		if err != nil {
			log.Printf("Could not read data from file for table %s: %v", table, err)
			return nil, err
		}
	}
	db := &MemDb[T]{
		Table:      table,
		Data:       data,
		FilePath:   filePath,
		SaveToDisk: saveToDisk,
	}
	return db, nil
}

func GetDataFromFile[T common.Serializable](filePath string) (map[string]T, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Could not read file, %s: %v", filePath, err)
		return map[string]T{}, err
	}

	data := map[string]T{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Printf("Error unmarshalling json into serializable data: %v", err)
		return make(map[string]T), err
	}
	return data, nil
}

func (db *MemDb[T]) Upsert(item T) error {
	id, _ := item.GetID()
	db.Data[id] = item
	if db.FilePath != "" {
		return db.AppendToFile(item, db.FilePath)
	}
	return nil
}

// UpsertAndSave upserts an item and immediately saves all data to the table file
func (db *MemDb[T]) UpsertAndSave(item T) error {
	// First upsert the item
	err := db.Upsert(item)
	if err != nil {
		return err
	}

	// Then save to file if RootPath is set
	if db.FilePath != "" {
		filePath := db.FilePath + "/" + db.Table + ".json"
		return db.SaveAllDataToFile(filePath)
	}

	return nil
}

func (db *MemDb[T]) GetData() ([]T, error) {
	data := []T{}
	for _, value := range db.Data {
		data = append(data, value)
	}
	return data, nil
}

// AppendToFile appends a single serializable item to a file as a JSON line
func (db *MemDb[T]) AppendToFile(item T, filePath string) error {
	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Could not open file for appending, %s: %v", filePath, err)
		return err
	}
	defer file.Close()

	// Marshal the item to JSON
	jsonData, err := json.Marshal(item)
	if err != nil {
		log.Printf("Error marshalling item to JSON: %v", err)
		return err
	}

	// Append the JSON line with a newline
	_, err = file.Write(append(jsonData, '\n'))
	if err != nil {
		log.Printf("Error writing to file %s: %v", filePath, err)
		return err
	}

	log.Printf("Successfully appended item to file: %s", filePath)
	return nil
}

// SaveAllDataToFile saves all current data to a file as a JSON array
func (db *MemDb[T]) SaveAllDataToFile(filePath string) error {
	// Convert map to slice
	data := []T{}
	for _, value := range db.Data {
		data = append(data, value)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshalling data to JSON: %v", err)
		return err
	}

	// Write to file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing to file %s: %v", filePath, err)
		return err
	}

	log.Printf("Successfully saved all data to file: %s", filePath)
	return nil
}

// GetByID implements the DbDriver interface
func (db *MemDb[T]) GetByID(id string) (T, error) {
	item, exists := db.Data[id]
	if !exists {
		return item, errors.New("item not found")
	}
	return item, nil
}

// Delete implements the DbDriver interface
func (db *MemDb[T]) Delete(id string) error {
	delete(db.Data, id)

	if db.SaveToDisk && db.FilePath != "" {
		return db.SaveAllDataToFile(db.FilePath)
	}
	return nil
}
