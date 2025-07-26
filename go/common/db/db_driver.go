package db

import "github.com/Hanasou/news_feed/go/common"

// Common interface database drivers can implement
type DbDriver interface {
	Upsert(common.Serializable) error
	GetData() ([]common.Serializable, error)
}
