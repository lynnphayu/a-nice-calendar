package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SpotifyCredentials struct {
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
	ExpiresIn    int    `bson:"expires_in" json:"expires_in"`
}

type User struct {
	ID                 string             `bson:"_id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Email              string             `bson:"email" json:"email"`
	Password           string             `bson:"password,omitempty" json:"password,omitempty"`
	SpotifyCredentials SpotifyCredentials `bson:"spotify_credentials" json:"spotify_credentials"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

func (u *User) HashedPassword(passwork string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwork), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func NewUser(email, name, password string) (*User, error) {
	user := &User{
		Email: email,
		Name:  name,
	}
	if password != "" {
		if err := user.HashedPassword(password); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
