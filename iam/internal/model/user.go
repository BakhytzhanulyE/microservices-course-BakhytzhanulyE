// Package model описывает доменные сущности IAM.
package model

import "time"

// User — пользователь системы. Пароль хранится только в виде хеша.
type User struct {
	UUID         string
	Login        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// Session — активная сессия пользователя, живёт в Redis.
type Session struct {
	UUID      string    `json:"uuid"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RegisterParams — данные для регистрации.
type RegisterParams struct {
	Login    string
	Email    string
	Password string
}

// LoginParams — данные для входа.
type LoginParams struct {
	Login    string
	Password string
}

// TokenPair — выданные пользователю токены.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
