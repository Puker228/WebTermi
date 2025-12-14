package session

type UserCache interface {
	Get()
	Set()
}
