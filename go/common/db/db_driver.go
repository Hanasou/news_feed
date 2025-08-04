package db

import "github.com/Hanasou/news_feed/go/common"

// Common interface database drivers can implement
type LegacyDbDriver interface {
	Upsert(common.Serializable) error
	GetData() ([]common.Serializable, error)
}

type DbDriver[T common.Serializable] interface {
	Upsert(T) error
	GetData() ([]T, error)
	GetByID(id string) (T, error)
	Delete(id string) error
}
