package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 一个抽象的用来处理任意类型的mongodb的仓储基类
type IRepository interface {
	FindAll() (dataList []interface{}, err error)
	CountByFilter(filter interface{}) (count int64, err error)
	//查找一条记录
	FindOne(filter interface{}, opts ...FindOneOption) (data interface{}, err error)
	FindByFilter(filter interface{}, opts ...FindOption) (dataList []interface{}, err error)
	FindByObjectId(id primitive.ObjectID) (dataList interface{}, err error)

	Create(data interface{}, opts ...*options.InsertOneOptions) error
	//更新一个IEntity接口的对象
	// entity必须实现IEntity接口
	FindOneAndUpdateEntityWithId(entity interface{}, opts ...*options.FindOneAndUpdateOptions) error
	FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error

	DeleteOne(id primitive.ObjectID, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOneByFilter(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	CreateIndexIfNotExist(indexDefine EntityIndexDefine, indexOptions *options.IndexOptions) (string, error)
}
