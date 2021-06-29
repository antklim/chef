package project

import "strings"

const (
	serverUnknown = "unknown"
	serverNone    = ""
	serverHTTP    = "http"
	serverGRPC    = "grpc"
)

func server(v string) string {
	switch strings.ToLower(v) {
	case "":
		return serverNone
	case "http":
		return serverHTTP
	case "grpc":
		return serverGRPC
	default:
		return serverUnknown
	}
}
