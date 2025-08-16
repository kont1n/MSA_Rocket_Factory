package postgres

// Option представляет функцию для настройки PostgreSQL контейнера
type Option func(*Config)

// WithContainerName устанавливает имя контейнера
func WithContainerName(name string) Option {
	return func(c *Config) {
		c.ContainerName = name
	}
}

// WithImageName устанавливает имя Docker образа
func WithImageName(imageName string) Option {
	return func(c *Config) {
		c.ImageName = imageName
	}
}

// WithDatabase устанавливает имя базы данных
func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

// WithAuth устанавливает имя пользователя и пароль
func WithAuth(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

// WithPort устанавливает порт для подключения
func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithNetworkName устанавливает имя Docker сети
func WithNetworkName(networkName string) Option {
	return func(c *Config) {
		c.NetworkName = networkName
	}
}

// WithLogger устанавливает логгер для контейнера
func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
