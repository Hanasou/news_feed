package db

import "github.com/Hanasou/news_feed/go/common"

type DbDriver[T common.Serializable] interface {
	Upsert(T) error
	GetData() ([]T, error)
	GetByID(id string) (T, error)
	Delete(id string) error
}
