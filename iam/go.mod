module github.com/kont1n/MSA_Rocket_Factory/iam

go 1.24.4

replace github.com/kont1n/MSA_Rocket_Factory/shared => ../shared

replace github.com/kont1n/MSA_Rocket_Factory/platform => ../platform

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/caarlos0/env/v11 v11.3.1
	github.com/gomodule/redigo v1.9.2
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/joho/godotenv v1.5.1
	github.com/kont1n/MSA_Rocket_Factory/platform v0.0.0-00010101000000-000000000000
	github.com/kont1n/MSA_Rocket_Factory/shared v0.0.0-00010101000000-000000000000
	github.com/pressly/goose v2.7.0+incompatible
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.41.0
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.8
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pkg/errors v0.9.1
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
)
