package mongodbr

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DefaultConfiguration = NewConfiguration()
	//默认的client
	DefaultClient *mongo.Client
)

// 构建默认的client
func SetupDefaultClient(uri string, opts ...func(*options.ClientOptions)) (*mongo.Client, error) {
	//测试能否连接
	clientOptions := options.Client().ApplyURI(uri)
	for _, eachOpt := range opts {
		eachOpt(clientOptions)
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("无法初始化mongodb,在连接到mongodb时出现异常,异常信息:%s", err.Error())
	}

	DefaultClient = client
	return DefaultClient, nil
}

func Ping() error {
	if DefaultClient == nil {
		return fmt.Errorf("先调用SetupDefaultClient方法创建好一个Client对象后,再调用此方法")
	}
	//测试ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := DefaultClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("mongodb ping测试时出现异常,异常信息:%s", err.Error())
	}
	fmt.Println("mongodb ping测试正常")
	return nil
}

type Configuration struct {
	QueryTimeout time.Duration

	//创建一条新的记录,并返回这条记录的指针地址
	createItemFunc func() interface{}
	//查询时设置默认的排序
	setDefaultSort func(*options.FindOptions) *options.FindOptions
}

func NewConfiguration() *Configuration {
	return &Configuration{
		QueryTimeout: 30 * time.Second,
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
