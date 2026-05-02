package db

import "github.com/Hanasou/news_feed/go/common"

type DbDriver[T common.Serializable] interface {
	Upsert(T) error
	GetAll() ([]T, error)
	GetByID(id string) (T, error)
	Delete(id string) error
	GetByField(field string, value any) (T, error)
	GetByFilter(filters map[string]any) ([]T, error)
}
