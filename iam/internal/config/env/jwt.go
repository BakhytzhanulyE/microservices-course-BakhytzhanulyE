package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type jwtEnvConfig struct {
	SecretKey  string        `env:"JWT_SECRET_KEY,required"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL"  envDefault:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
}

type jwtConfig struct {
	raw jwtEnvConfig
}

// NewJWTConfig читает настройки токенов.
func NewJWTConfig() (*jwtConfig, error) {
	var raw jwtEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &jwtConfig{raw: raw}, nil
}

// SecretKey — ключ подписи токенов.
func (cfg *jwtConfig) SecretKey() string { return cfg.raw.SecretKey }

// AccessTTL — время жизни access-токена.
func (cfg *jwtConfig) AccessTTL() time.Duration { return cfg.raw.AccessTTL }

// RefreshTTL — время жизни сессии и refresh-токена.
func (cfg *jwtConfig) RefreshTTL() time.Duration { return cfg.raw.RefreshTTL }
