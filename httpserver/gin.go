package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/jackdes/go-package-management/common/logger"
	"sync"
)

var (
	ginMode      string
	ginOffLogger bool
	defaultPort  = 19593
)

type Config struct {
	Port         int    `json:"http_port"`
	BindAddr     string `json:"http_bind_addr"`
	GinNoDefault bool   `json:"http_no_default"`
}

type GinService interface {
	Port() int
	isGinService()
}

type ginService struct {
	Config
	isEnabled bool
	name      string
	logger    logger.Logger
	srv       *myHttpServer
	router    *gin.Engine
	mu        *sync.Mutex
	handlers  []func(*gin.Engine)
}

func NewService(name string) *ginService {
	return &ginService{
		name:     name,
		mu:       &sync.Mutex{},
		handlers: []func(*gin.Engine){},
	}
}
