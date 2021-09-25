package app

import (
	"github.com/sirupsen/logrus"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id       string            // 应用ID
	name     string            // 应用名称
	version  string            // 应用版本
	metadata map[string]string // 元数据信息

	log    *logrus.Logger // 日志
	server []Serve
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Logger with service logger.
func Logger(logger *logrus.Logger) Option {
	return func(o *options) { o.log = logger }
}

func Server(srv ...Serve) Option {
	return func(o *options) {
		o.server = srv
	}
}
