package memdb

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Hanasou/news_feed/go/common"
)

// Driver for data held in memory
type MemDb struct {
	Table      string
	Data       map[string]common.Serializable
	FilePath   string
	SaveToDisk bool
}

func (db *MemDb) String() string {
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

func Initialize(table string, rootPath string, saveToDisk bool) (*MemDb, error) {
	data := map[string]common.Serializable{}
	filePath := ""
	if saveToDisk {
		var err error
		filePath = rootPath + "/" + table + ".json"
		data, err = GetDataFromFile(rootPath + "/" + table + ".json")
		if err != nil {
			log.Printf("Could not read data from file for table %s: %v", table, err)
			return nil, err
		}
	}
	db := &MemDb{
		Table:    table,
		Data:     data,
		FilePath: filePath,
	}
	return db, nil
}

func GetDataFromFile(filePath string) (map[string]common.Serializable, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Could not read file, %s: %v", filePath, err)
		return map[string]common.Serializable{}, err
	}

	data := map[string]common.Serializable{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Printf("Error unmarshalling json into serializable data: %v", err)
		return map[string]common.Serializable{}, err
	}
	return data, nil
}

func (db *MemDb) Upsert(item common.Serializable) error {
	id, _ := item.GetId()
	db.Data[id] = item
	if db.FilePath != "" {
		return db.AppendToFile(item, db.FilePath)
	}
	return nil
}

// UpsertAndSave upserts an item and immediately saves all data to the table file
func (db *MemDb) UpsertAndSave(item common.Serializable) error {
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

func (db *MemDb) GetData() ([]common.Serializable, error) {
	data := []common.Serializable{}
	for _, value := range db.Data {
		data = append(data, value)
	}
	return data, nil
}

// AppendToFile appends a single serializable item to a file as a JSON line
func (db *MemDb) AppendToFile(item common.Serializable, filePath string) error {
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
func (db *MemDb) SaveAllDataToFile(filePath string) error {
	// Convert map to slice
	data := []common.Serializable{}
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
