package data

type Authenticator interface {
	Authenticate(user string, pass string) bool
	Create(user string) string
	Use(token string) bool
}
