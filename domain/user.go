package domain

type UserID string

type User struct {
	ID        UserID
	Name      string
	SecretKey string
	CreatedAt int
	UpdatedAt int
}
