module github.com/kont1n/MSA_Rocket_Factory/payment

go 1.24

replace github.com/kont1n/MSA_Rocket_Factory/shared => ../shared

replace github.com/kont1n/MSA_Rocket_Factory/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	github.com/joho/godotenv v1.5.1
	github.com/kont1n/MSA_Rocket_Factory/platform v0.0.0-00010101000000-000000000000
	github.com/kont1n/MSA_Rocket_Factory/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.74.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
