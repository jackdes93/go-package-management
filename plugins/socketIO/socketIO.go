package socketIO

//
//import (
//	"flag"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	socketio "github.com/googollee/go-socket.io"
//	go_package_manager "github.com/jackdes93/go-package-management"
//	"github.com/jackdes93/go-package-management/common"
//	"github.com/jackdes93/go-package-management/common/logger"
//	"net/http"
//)
//
//type Socket interface {
//	Id() string
//	Rooms() []string
//	Request() *http.Request
//	On(event string, f interface{}) error
//	Emit(event string, args ...interface{}) error
//	Join(room string) error
//	Leave(room string) error
//	Disconnect()
//	BroadcastTo(room, event string, args ...interface{}) error
//}
//
//type AppSocket interface {
//	ServiceContext() go_package_manager.ServiceContext
//	Logger() logger.Logger
//	CurrentUser() common.Requester
//	SetCurrentUser(requester common.Requester)
//	BroadcastToRoom(room, event string, args ...interface{})
//	String() string
//	Socket
//}
//
//type Config struct {
//	Name          string
//	MaxConnection int
//}
//
//type socketServer struct {
//	Config
//	io     *socketio.Server
//	logger logger.Logger
//}
//
//func NewSocketServer(name string) *socketServer {
//	return &socketServer{
//		Config: Config{Name: name},
//	}
//}
//
//type ObserverProvider interface {
//	AddObservers(server *socketio.Server, sc go_package_manager.ServiceContext, l logger.Logger) func(socket socketio.)
//}
//
//func (s *socketServer) StartRealtimeServer(engine *gin.Engine, sc go_package_manager.ServiceContext, op ObserverProvider) {
//	server := socketio.NewServer(nil)
//	s.io = server
//	_ = s.io.OnConnect("/", func(conn socketio.Conn) error {
//		op.AddObservers(server, sc, s.logger)
//		return nil
//	})
//	engine.GET("/socket.io/", gin.WrapH(server))
//	engine.POST("/socket.io/", gin.WrapH(server))
//}
//
//func (s *socketServer) GetPrefix() string {
//	return s.Config.Name
//}
//
//func (s *socketServer) Get() interface{} {
//	return s
//}
//
//func (s *socketServer) Name() string {
//	return s.Config.Name
//}
//
//func (s *socketServer) InitFlags() {
//	pre := s.GetPrefix()
//	flag.IntVar(&s.MaxConnection, fmt.Sprintf("%s-max-connection", pre), 2000, "socket max connection")
//}
//
//func (s *socketServer) Configure() error {
//	s.logger = logger.GetCurrent().GetLogger("io.socket")
//	return nil
//}
//
//func (s *socketServer) Run() error {
//	return s.Configure()
//}
//
//func (s *socketServer) Stop() <-chan bool {
//	c := make(chan bool)
//	go func() { c <- true }()
//	return c
//}
