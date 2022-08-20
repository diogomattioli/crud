package data

type Authenticator interface {
	Authenticate(user string, pass string) bool
	Create() string
	Use(token string) bool
}
