module github.com/fr0g-vibe/fr0g-ai-aip

go 1.23.0

toolchain go1.24.3

// No external dependencies - using only Go standard library
// gRPC dependencies can be added later when needed:
// google.golang.org/grpc
// google.golang.org/protobuf

require (
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
)
