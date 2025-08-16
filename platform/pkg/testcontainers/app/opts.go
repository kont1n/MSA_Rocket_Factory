package app

import (
	"io"

	"github.com/testcontainers/testcontainers-go/wait"
)

type Option func(*Config)

func WithName(name string) Option {
	return func(c *Config) {
		c.Name = name
	}
}

func WithDockerfile(dir, file string) Option {
	return func(c *Config) {
		c.DockerfileDir = dir
		c.Dockerfile = file
	}
}

func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithNetwork(name string) Option {
	return func(c *Config) {
		c.Networks = append(c.Networks, name)
	}
}

func WithEnv(env map[string]string) Option {
	return func(c *Config) {
		for k, v := range env {
			c.Env[k] = v
		}
	}
}

func WithLogOutput(out io.Writer) Option {
	return func(c *Config) {
		c.LogOutput = out
	}
}

func WithStartupWait(strategy wait.Strategy) Option {
	return func(c *Config) {
		c.StartupWait = strategy
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
