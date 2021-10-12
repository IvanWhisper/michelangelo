package generate

import (
	"github.com/go-redis/redis/v8"
)

func NewCache(db *redis.Client) func() (interface{}, error) {
	return func() (interface{}, error) {
		return db, nil
	}
}
