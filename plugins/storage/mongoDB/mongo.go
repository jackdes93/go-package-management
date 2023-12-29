package mongoDB

import (
	"context"
	"flag"
	"github.com/jackdes/go-package-management/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math"
	"sync"
	"time"
)

var (
	defaultDBName  = "defaultMongoDB"
	DefaultMongoDB = getDefaultMongoDB()
)

const retryCount = 10

type MongoDBOpt struct {
	MgoUri       string
	Prefix       string
	PingInterval int // in seconds
}

type mongoDB struct {
	name      string
	logger    logger.Logger
	client    *mongo.Client
	isRunning bool
	once      *sync.Once
	*MongoDBOpt
}

func getDefaultMongoDB() *mongoDB {
	return NewMongoDB(defaultDBName, "")
}

func NewMongoDB(name, prefix string) *mongoDB {
	return &mongoDB{
		name:      name,
		isRunning: false,
		once:      new(sync.Once),
		MongoDBOpt: &MongoDBOpt{
			Prefix: prefix,
		},
	}
}

func (mgDB *mongoDB) GetPrefix() string {
	return mgDB.Prefix
}

func (mgDB *mongoDB) Name() string {
	return mgDB.name
}

func (mgDB *mongoDB) InitFlags() {
	prefix := mgDB.Prefix
	if mgDB.Prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&mgDB.MgoUri, prefix+"mgo-uri", "", "MongoDB connection-string. Ex: mongodb://...")
	flag.IntVar(&mgDB.PingInterval, prefix+"mgo-ping-interval", 5, "MongoDB ping check interval")
}

func (mgDB *mongoDB) isDisabled() bool {
	return mgDB.MgoUri == ""
}

func (mgDB *mongoDB) Configure() error {
	if mgDB.isDisabled() || mgDB.isRunning {
		return nil
	}

	mgDB.logger = logger.GetCurrent().GetLogger(mgDB.name)
	mgDB.logger.Info("Connect to Mongodb at ", mgDB.MgoUri, " ...")

	var err error
	mgDB.client, err = mgDB.getConnWithRetry(context.Background(), retryCount)
	if err != nil {
		mgDB.logger.Error("Error connect to mongodb at ", mgDB.MgoUri, ". ", err.Error())
		return err
	}
	mgDB.isRunning = true
	return nil
}

func (mgDB *mongoDB) Cleanup() {
	if mgDB.isDisabled() {
		return
	}

	if mgDB.client != nil {
		err := mgDB.client.Disconnect(context.Background())
		if err != nil {
			mgDB.logger.Errorf("error mongodb disconnect, Error: %v", err)
			return
		}
	}
}

func (mgDB *mongoDB) Run() error {
	return mgDB.Configure()
}

func (mgDB *mongoDB) Stop() <-chan bool {
	if mgDB.client != nil {
		err := mgDB.client.Disconnect(context.Background())
		if err != nil {
			mgDB.logger.Errorf("error mongodb disconnect, Error: %v", err)
		}
	}
	mgDB.isRunning = false

	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (mgDB *mongoDB) Get() interface{} {
	mgDB.once.Do(func() {
		if !mgDB.isRunning && !mgDB.isDisabled() {
			if client, err := mgDB.getConnWithRetry(context.Background(), math.MaxInt32); err == nil {
				mgDB.client = client
				mgDB.isRunning = true
			} else {
				mgDB.logger.Fatalf("%s connection cannot reconnect\n", mgDB.name)
			}
		}
	})

	if mgDB.client == nil {
		return nil
	}
	return mgDB.client
}

func (mgDB *mongoDB) getConnWithRetry(ctx context.Context, retryCount int) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(mgDB.MgoUri)
	authMongo := options.Credential{
		AuthSource: mgDB.name,
	}
	clientOpts.SetAuth(authMongo)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		for {
			time.Sleep(time.Second * 1)
			mgDB.logger.Errorf("Retry to connect %s.\n", mgDB.name)
			err = client.Ping(context.Background(), readpref.Primary())
			if err == nil {
				go mgDB.reconnectIfNeeded()
				break
			}
		}
	} else {
		go mgDB.reconnectIfNeeded()
	}
	return client, err
}

func (mgDB *mongoDB) reconnectIfNeeded() {
	conn := mgDB.client
	for {
		if err := conn.Ping(context.Background(), readpref.Primary()); err != nil {
			err := conn.Disconnect(context.Background())
			if err != nil {
				return
			}
			mgDB.logger.Errorf("%s connection is gone, try to reconnect\n", mgDB.name)
			mgDB.isRunning = false
			mgDB.once = new(sync.Once)

			if err := mgDB.Get().(*mongo.Client).Disconnect(context.Background()); err != nil {
				return
			}
			return
		}
		time.Sleep(time.Second * time.Duration(mgDB.PingInterval))
	}
}
