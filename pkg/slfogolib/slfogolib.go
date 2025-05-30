package slfogolib

import (
	"fmt"
	"strconv"

	"gopkg.in/mcuadros/go-syslog.v2"
)

const (
	defaultPort = 9999
)

func GetPort(ps string) int64 {
	port, err := strconv.Atoi(ps)
	if err != nil {
		return defaultPort
	}

	return int64(port)
}

func BuildConnectString(ps string) string {
	return fmt.Sprintf("0.0.0.0:%d", GetPort(ps))
}

func BuildServer() (*syslog.Server, syslog.LogPartsChannel) {
	lpChan := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(lpChan)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)

	return server, lpChan
}
