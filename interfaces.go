package go_package_manager

import (
	"github.com/gin-gonic/gin"
	"github.com/jackdes93/go-package-management/common/logger"
)

type Option func(*service)

type HttpServerHandler = func(*gin.Engine)

type Function func(ctx ServiceContext) error

type Storage interface {
	Get(prefix string) (interface{}, error)
	MustGet(prefix string) interface{}
}

type PrefixRunnable interface {
	HashPrefix
	Runnable
}

type HashPrefix interface {
	GetPrefix() string
	Get() interface{}
}

// Service Interface -The heart of SDK, Service represents for a real micro service
// with its all components
type Service interface {
	// ServiceContext - A part of Service, it's passed to all handlers/functions
	ServiceContext
	// Name of the service
	Name()
	// Version of the service
	Version()
	// HttpServer - Gin HTTP Server wrapper
	HttpServer() HttpServer
	// Init with options, they can be db connections or
	// anything the service need handle before starting
	Init() error
	// IsRegistered - This method returns service if it is registered on discovery
	IsRegistered() bool
	// Start service and it's all component.
	// It will be stopped if any service return error
	Start() error
	// Stop service and it's all component.
	Stop()
	// OutEnv - Method export all flags to std/terminal
	// We might use: "> .env" to move its content .env file
	OutEnv()
}

type ServiceContext interface {
	// Logger for a specific service, usually it has a prefix to distinguish
	// with each others
	Logger(prefix string) logger.Logger
	// Get component with prefix
	Get(prefix string) (interface{}, bool)
	MustGet(prefix string) interface{}
	Env() string
}

type Runnable interface {
	Name() string
	InitFlags()
	Configure() error
	Run() error
	Stop() <-chan bool
}

type HttpServer interface {
	Runnable
	// AddHandler - Add handlers to GIN Server
	AddHandler(HttpServerHandler)
	// URI - That the server is listening
	URI() string
}
