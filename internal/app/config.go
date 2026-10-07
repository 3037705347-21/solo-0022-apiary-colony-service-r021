package app

import "os"

const DefaultAddr = "127.0.0.1:18093"

const addrEnv = "APIARY_ADDR"

type Config struct {
	Addr string
}

func ConfigFromEnv() Config {
	addr := os.Getenv(addrEnv)
	if addr == "" {
		addr = DefaultAddr
	}
	return Config{Addr: addr}
}
