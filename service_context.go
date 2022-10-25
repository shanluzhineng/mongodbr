package mongodbr

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceContextOption func() IServiceContextProvider

// 使用默认的context提供者,超时并自动取消
func WithDefaultServiceContext() ServiceContextOption {
	return func() IServiceContextProvider {
		return NewDefaultServiceContextProvider()
	}
}

// 使用mongodb的事条上下文
func WithMongodbSessionContext(sessionContext mongo.SessionContext) ServiceContextOption {
	return func() IServiceContextProvider {
		return newMongodbSessionContext(sessionContext)
	}
}

// 用来提供service执行过程中所需的上下文
type IServiceContextProvider interface {
	GetContext() context.Context
	GetCancelFunc() context.CancelFunc
}

type DefaultServiceContextProvider struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// 创建一个默认的IServiceContextProvider实现
func NewDefaultServiceContextProvider() IServiceContextProvider {
	ctx, cancelFunc := context.WithTimeout(context.Background(), DefaultConfiguration.QueryTimeout)
	contextProvider := &DefaultServiceContextProvider{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	return contextProvider
}

func (s *DefaultServiceContextProvider) GetContext() context.Context {
	return s.ctx
}

func (s *DefaultServiceContextProvider) GetCancelFunc() context.CancelFunc {
	return s.cancelFunc
}

type MongodbSessionServiceContextProvider struct {
	sessionContext mongo.SessionContext
}

// 创建一个使用mongodb事务的IServiceContextProvider实现
func newMongodbSessionContext(sessionContext mongo.SessionContext) IServiceContextProvider {
	return &MongodbSessionServiceContextProvider{
		sessionContext: sessionContext,
	}
}

func (s *MongodbSessionServiceContextProvider) GetContext() context.Context {
	return s.sessionContext
}

func (s *MongodbSessionServiceContextProvider) GetCancelFunc() context.CancelFunc {
	return func() {}
}
