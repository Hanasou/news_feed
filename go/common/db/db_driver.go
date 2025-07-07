package db

import (
	"github.com/Hanasou/news_feed/go/common/commodels"
)

// Common interface database drivers can implement
type DbDriver interface {
	Upsert(commodels.Serializable) error
	GetData() ([]commodels.Serializable, error)
}
