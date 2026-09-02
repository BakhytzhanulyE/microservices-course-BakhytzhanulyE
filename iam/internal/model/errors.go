package model

import "errors"

var (
	// ErrUserNotFound — пользователя с такими данными нет.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists — логин или почта уже заняты.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidCredentials — неверная пара логин/пароль.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidToken — токен просрочен, подделан или отозван.
	ErrInvalidToken = errors.New("invalid token")
	// ErrSessionNotFound — сессии больше нет.
	ErrSessionNotFound = errors.New("session not found")
	// ErrWeakPassword — пароль короче минимальной длины.
	ErrWeakPassword = errors.New("password is too short")
)
