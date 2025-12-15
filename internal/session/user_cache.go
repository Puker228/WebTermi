package session

type UserCache interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}
