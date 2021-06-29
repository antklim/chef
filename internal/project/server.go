package project

import "strings"

type Server string

const (
	// ServerUnknown represents unknown server option.
	ServerUnknown Server = "unknown"
	// ServerNone represents no server option.
	ServerNone Server = ""
	// ServerHTTP represents http server option.
	ServerHTTP Server = "http"
	// ServerGRPC represents grpc server option.
	ServerGRPC Server = "grpc"
)

func NewServer(v string) Server {
	switch strings.ToLower(v) {
	case "":
		return ServerNone
	case "http":
		return ServerHTTP
	case "grpc":
		return ServerGRPC
	default:
		return ServerUnknown
	}
}

func (s Server) IsUnknown() bool {
	return s == ServerUnknown
}
