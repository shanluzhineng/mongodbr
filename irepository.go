package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 一个抽象的用来处理任意类型的mongodb的仓储基类
type IRepository interface {
	IEntityFind
	IEntityCreate
	IEntityUpdate
	IEntityDelete
	IEntityIndex
	IEntityBulkWrite

	// aggregate
	Aggregate(pipeline interface{}, dataList interface{}, opts ...AggregateOption) (err error)

	// replace*
	ReplaceById(id primitive.ObjectID, doc interface{}, opts ...*options.ReplaceOptions) (err error)
	Replace(filter interface{}, doc interface{}, opts ...*options.ReplaceOptions) (err error)

	GetName() (name string)
	GetCollection() (c *mongo.Collection)
}

type IEntityCreate interface {
	// create
	Create(data interface{}, opts ...*options.InsertOneOptions) (id primitive.ObjectID, err error)
	CreateMany(itemList []interface{}, opts ...*options.InsertManyOptions) (ids []primitive.ObjectID, err error)
}

type IEntityDelete interface {
	// delete
	DeleteOne(id primitive.ObjectID, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOneByFilter(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}
