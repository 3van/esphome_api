package model

import (
	"errors"
	"fmt"

	"github.com/mycontroller-org/esphome_api/pkg/api"
	"google.golang.org/protobuf/proto"
)

// Error types
var (
	ErrCommunicationTimeout  = errors.New("esphome_api: communication timeout")
	ErrConnRequireEncryption = errors.New("esphome_api: connection requires encryption")
)

// call back function used to report received messages
type CallBackFunc func(proto.Message)

// DeviceInfo struct
type DeviceInfo struct {
	Name            string
	Model           string
	MacAddress      string
	EsphomeVersion  string
	CompilationTime string
	HasDeepSleep    bool
}

func (di *DeviceInfo) String() string {
	return fmt.Sprintf("{name: %v, model:%v, mac_address:%v, esphome_version:%v, compilation_time:%v, has_deep_sleep:%v}",
		di.Name, di.Model, di.MacAddress, di.EsphomeVersion, di.CompilationTime, di.HasDeepSleep)
}

// LogLevel type
type LogLevel int32

// log levels
const (
	LogLevelNone LogLevel = iota
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug // default
	LogLevelVerbose
	LogLevelVeryVerbose
)

// LogEntry of a message
type LogEntry struct {
	Level   LogLevel
	Tag     string
	Message string
}

func (le *LogEntry) String() string {
	return fmt.Sprintf("{level: %v, tag:%v, message:[%v]}",
		le.Level, le.Tag, le.Message)
}

func GetLogEntry(msg proto.Message) (*LogEntry, error) {
	entry, ok := msg.(*api.SubscribeLogsResponse)
	if !ok {
		return nil, fmt.Errorf("received invalid data type:%T", msg)
	}
	log := LogEntry{
		Level:   LogLevel(entry.Level),
		Message: string(entry.Message),
	}
	return &log, nil
}

type HelloResponse struct {
	ApiVersionMajor uint32
	ApiVersionMinor uint32
	ServerInfo      string
	Name            string
}

func (hr *HelloResponse) String() string {
	return fmt.Sprintf("{name: %v, api_version_major: %v, api_version_minor:%v, server_info:%v}",
		hr.Name, hr.ApiVersionMajor, hr.ApiVersionMinor, hr.ServerInfo)
}
