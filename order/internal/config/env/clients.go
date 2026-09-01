package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type clientsEnvConfig struct {
	InventoryHost string `env:"INVENTORY_GRPC_HOST,required"`
	InventoryPort string `env:"INVENTORY_GRPC_PORT,required"`
	PaymentHost   string `env:"PAYMENT_GRPC_HOST,required"`
	PaymentPort   string `env:"PAYMENT_GRPC_PORT,required"`
}

type clientsConfig struct {
	raw clientsEnvConfig
}

// NewClientsConfig читает адреса соседних сервисов.
func NewClientsConfig() (*clientsConfig, error) {
	var raw clientsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &clientsConfig{raw: raw}, nil
}

// InventoryAddress — адрес каталога деталей.
func (cfg *clientsConfig) InventoryAddress() string {
	return net.JoinHostPort(cfg.raw.InventoryHost, cfg.raw.InventoryPort)
}

// PaymentAddress — адрес платёжного сервиса.
func (cfg *clientsConfig) PaymentAddress() string {
	return net.JoinHostPort(cfg.raw.PaymentHost, cfg.raw.PaymentPort)
}
