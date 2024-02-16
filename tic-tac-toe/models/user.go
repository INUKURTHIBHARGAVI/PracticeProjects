package models

type User struct {
	Username string
	Email    string
	Photo    string
}

type UserBuilder struct {
	Username string
	email    string
	photo    []byte
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}
func (u *UserBuilder) WithUsername(username string) *UserBuilder {
	u.Username = username
	return u
}

func (u *UserBuilder) WithEmail(email string) *UserBuilder {
	u.email = email
	return u
}

func (u *UserBuilder) WithPhoto(photo []byte) *UserBuilder {
	u.photo = photo
	return u
}

func (u *UserBuilder) Build() *User {
	return &User{Username: u.Username, Email: u.email, Photo: string(u.photo)}
}
