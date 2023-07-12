package mongodbr

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultConfiguration = NewConfiguration()
	//默认的client
	DefaultClient *mongo.Client
	_cachedClient map[string]*mongo.Client = make(map[string]*mongo.Client)
)

// enable mongodb monitor
func EnableMongodbMonitor() func(*options.ClientOptions) {
	return func(co *options.ClientOptions) {
		monitor := &event.CommandMonitor{
			Started: func(_ context.Context, e *event.CommandStartedEvent) {
				log.Println(e.Command.String())
			},
			Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {
				log.Println(e.Reply.String())
			},
			Failed: func(ctx context.Context, e *event.CommandFailedEvent) {
				log.Println("mongodb error:", e.Failure)
			},
		}

		co.SetMonitor(monitor)
	}
}

// 构建默认的client
func SetupDefaultClient(uri string, opts ...func(*options.ClientOptions)) (*mongo.Client, error) {
	client, err := createClient(uri, opts...)
	if err != nil {
		return nil, err
	}
	DefaultClient = client
	return DefaultClient, nil
}

func RegistClient(key string, uri string, opts ...func(*options.ClientOptions)) (*mongo.Client, error) {
	client, err := createClient(uri, opts...)
	if err != nil {
		return nil, err
	}
	_cachedClient[key] = client
	return client, nil
}

// get client by key
func GetClient(key string) *mongo.Client {
	client, ok := _cachedClient[key]
	if !ok {
		return nil
	}
	return client
}

func createClient(uri string, opts ...func(*options.ClientOptions)) (*mongo.Client, error) {
	//测试能否连接
	clientOptions := options.Client().ApplyURI(uri)
	for _, eachOpt := range opts {
		eachOpt(clientOptions)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("无法初始化mongodb,在连接到mongodb时出现异常,异常信息:%s", err.Error())
	}
	return client, nil
}

type Configuration struct {
	QueryTimeout time.Duration

	//创建一条新的记录,并返回这条记录的指针地址
	createItemFunc func() interface{}
	//查询时设置默认的排序
	setDefaultSort func(*options.FindOptions) *options.FindOptions
}

func CreateContext(c *Configuration) (context.Context, context.CancelFunc) {
	if c == nil || c.QueryTimeout <= 0 {
		ctx := context.TODO()
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(context.Background(), c.QueryTimeout)
}

func (c *Configuration) safeCreateItem() interface{} {
	if c.createItemFunc == nil {
		return make(map[string]interface{})
	}
	return c.createItemFunc()
}

func NewConfiguration() *Configuration {
	return &Configuration{
		QueryTimeout: 120 * time.Second,
	}
}

type RepositoryOption func(*Configuration)

func WithDefaultSort(defaultSortFunc func(*options.FindOptions) *options.FindOptions) RepositoryOption {
	return func(configuration *Configuration) {
		configuration.setDefaultSort = defaultSortFunc
	}
}

func WithCreateItemFunc(createItemFunc func() interface{}) RepositoryOption {
	return func(configuration *Configuration) {
		configuration.createItemFunc = createItemFunc
	}
}
