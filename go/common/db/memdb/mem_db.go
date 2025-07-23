package memdb

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Hanasou/news_feed/go/common/commodels"
)

// Driver for data held in memory
type MemDb struct {
	Table string
	Data  map[string]commodels.Serializable
}

func Initialize(table string) (*MemDb, error) {
	db := &MemDb{
		Table: table,
		Data:  map[string]commodels.Serializable{},
	}
	return db, nil
}

func GetDataFromFile(filePath string) (map[string]commodels.Serializable, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Could not read file, %s: %v", filePath, err)
		return map[string]commodels.Serializable{}, err
	}

	data := map[string]commodels.Serializable{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Printf("Error unmarshalling json into serializable data: %v", err)
		return map[string]commodels.Serializable{}, err
	}
	return data, nil
}

func (db *MemDb) Upsert(item commodels.Serializable) error {
	id, _ := item.GetId()
	db.Data[id] = item
	return nil
}

func (db *MemDb) GetData() ([]commodels.Serializable, error) {
	data := []commodels.Serializable{}
	for _, value := range db.Data {
		data = append(data, value)
	}
	return data, nil
}
