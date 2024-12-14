package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              uint   `gorm:"primaryKey"`
	Username        string `gorm:"unique;not null" json:"username"`
	Email           string `gorm:"unique;not null" json:"email"`
	Password        string `json:"password,omitempty" gorm:"-"`
	PasswordHash    string `gorm:"not null" json:"-"`
	OAuthProvider   string `gorm:"default:''" json:"oauth_provider"` // e.g., google, facebook
	OAuthProviderID string `gorm:"unique" json:"oauth_provider_id"`  // e.g., Google ID
	OAuthToken      string `json:"-"`
	Name            string `json:"name"`
	Role            string `gorm:"default:'user'" json:"role"` // e.g., admin, user, editor
	ProfilePicture  string `json:"profile_picture"`
}
