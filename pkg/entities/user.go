package entities

type User struct {
	Authenticated bool

	ID          string
	Name        string
	Email       string
	PhoneNumber string
}

// NewGuestUser returns unauthenticated "guest" user
func NewGuestUser() *User {
	return &User{
		Authenticated: false,
	}
}
