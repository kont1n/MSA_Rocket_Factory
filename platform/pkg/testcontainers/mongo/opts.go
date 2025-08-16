package mongo

type Option func(*Config)

func WithNetworkName(network string) Option {
	return func(c *Config) {
		c.NetworkName = network
	}
}

func WithContainerName(containerName string) Option {
	return func(c *Config) {
		c.ContainerName = containerName
	}
}

func WithImageName(image string) Option {
	return func(c *Config) {
		c.ImageName = image
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

func WithAuth(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

func WithAuthDB(authDB string) Option {
	return func(c *Config) {
		c.AuthDB = authDB
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
