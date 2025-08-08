package env

import (
	"os"
	"strconv"
)

type httpConfig struct {
	address           string
	readHeaderTimeout int
	shutdownTimeout   int
}

func NewHTTPConfig() (*httpConfig, error) {
	address := os.Getenv("HTTP_ADDRESS")
	if address == "" {
		address = "localhost:8080"
	}

	readHeaderTimeoutStr := os.Getenv("HTTP_READ_HEADER_TIMEOUT")
	readHeaderTimeout, err := strconv.Atoi(readHeaderTimeoutStr)
	if err != nil {
		readHeaderTimeout = 5 // секунды
	}

	shutdownTimeoutStr := os.Getenv("HTTP_SHUTDOWN_TIMEOUT")
	shutdownTimeout, err := strconv.Atoi(shutdownTimeoutStr)
	if err != nil {
		shutdownTimeout = 10 // секунды
	}

	return &httpConfig{
		address:           address,
		readHeaderTimeout: readHeaderTimeout,
		shutdownTimeout:   shutdownTimeout,
	}, nil
}

func (c *httpConfig) Address() string {
	return c.address
}

func (c *httpConfig) ReadHeaderTimeout() int {
	return c.readHeaderTimeout
}

func (c *httpConfig) ShutdownTimeout() int {
	return c.shutdownTimeout
}
