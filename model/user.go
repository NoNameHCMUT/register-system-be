package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleStudent   Role = "student"
	RoleSchool    Role = "school"
	RoleCommunity Role = "community"
)

type User struct {
	ID            uint        `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username      string      `gorm:"uniqueIndex;not null" json:"username"`
	FullName      string      `gorm:"not null" json:"full_name"`
	PasswordHash  string      `gorm:"not null" json:"-"`
	Role          Role        `gorm:"type:varchar(20);default:'student'" json:"role"`
	IsActive      bool        `gorm:"default:false" json:"is_active"`
	CreatedAt     time.Time   `json:"created_at"`
	StudentID     string      `json:"student_id"`
	Email         string      `gorm:"uniqueIndex;not null" json:"email"`
	AffiliationID uint        `gorm:"not null" json:"affiliation_id"`
	Affiliation   Affiliation `gorm:"foreignKey:AffiliationID" json:"affiliation,omitempty"`
	RefreshToken  string      `json:"-"`
}

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    uint      `json:"user_id"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}