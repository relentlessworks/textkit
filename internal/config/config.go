package config

import (
	"flag"
	"os"
)

// Config holds application configuration.
type Config struct {
	Addr   string
	NoAuth bool
}

// Load loads config from defaults < env < flags.
func Load() *Config {
	c := &Config{
		Addr:   ":8080",
		NoAuth: false,
	}
	if v := os.Getenv("TEXTKIT_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("TEXTKIT_NO_AUTH"); v == "1" || v == "true" {
		c.NoAuth = true
	}
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.BoolVar(&c.NoAuth, "no-auth", c.NoAuth, "disable auth (dev mode)")
	flag.Parse()
	return c
}
