module github.com/kont1n/MSA_Rocket_Factory/inventory

go 1.24.4

require (
	github.com/kont1n/MSA_Rocket_Factory/shared v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.51.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kont1n/MSA_Rocket_Factory/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250715232539-7130f93afb79 // indirect
)
