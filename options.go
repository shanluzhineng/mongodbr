package mongodbr

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultConfiguration = NewConfiguration()
)

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
