package dependencies

type ICache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
