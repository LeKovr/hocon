module github.com/LeKovr/hocon

go 1.19

replace SELF => ./

require (
	SELF v0.0.0-00010101000000-000000000000
	github.com/felixge/httpsnoop v1.0.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/soheilhy/cmux v0.1.5
	google.golang.org/genproto v0.0.0-20220909194730-69f6226f97e5
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
	golang.org/x/text v0.3.7 // indirect
)
