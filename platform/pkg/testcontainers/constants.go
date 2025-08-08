package testcontainers

// MongoDB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)

// PostgreSQL constants
const (
	// PostgreSQL container constants
	PostgresContainerName = "postgres"
	PostgresPort          = "5432"
	PostgresImageName     = "postgres:15-alpine"
	PostgresDatabase      = "order_service"
	PostgresUsername      = "postgres"
	PostgresPassword      = "postgres" //nolint:gosec

	// PostgreSQL environment variables
	PostgresImageNameKey = "POSTGRES_IMAGE_NAME"
	PostgresHostKey      = "POSTGRES_HOST"
	PostgresPortKey      = "POSTGRES_PORT"
	PostgresDatabaseKey  = "POSTGRES_DATABASE"
	PostgresUsernameKey  = "POSTGRES_USER"
	PostgresPasswordKey  = "POSTGRES_PASSWORD" //nolint:gosec
)
