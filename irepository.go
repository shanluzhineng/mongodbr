package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 一个抽象的用来处理任意类型的mongodb的仓储基类
type IRepository interface {
	CountByFilter(filter interface{}) (count int64, err error)

	// find
	FindAll(opts ...FindOption) (dataList []interface{}, err error)
	FindByObjectId(id primitive.ObjectID) (dataList interface{}, err error)
	FindOne(filter interface{}, opts ...FindOneOption) (data interface{}, err error)
	FindByFilter(filter interface{}, opts ...FindOption) (dataList []interface{}, err error)

	// aggregate
	Aggregate(pipeline interface{}, dataList interface{}, opts ...AggregateOption) (err error)

	// create
	Create(data interface{}, opts ...*options.InsertOneOptions) (id primitive.ObjectID, err error)
	CreateMany(itemList []interface{}, opts ...*options.InsertManyOptions) (ids []primitive.ObjectID, err error)

	// update
	// FindOneAndUpdateEntityWithId(entity interface{}, opts ...*options.FindOneAndUpdateOptions) error
	FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error

	// replace*
	ReplaceById(id primitive.ObjectID, doc interface{}, opts ...*options.ReplaceOptions) (err error)
	Replace(filter interface{}, doc interface{}, opts ...*options.ReplaceOptions) (err error)

	// delete
	DeleteOne(id primitive.ObjectID, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOneByFilter(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// index
	CreateIndex(indexDefine EntityIndexDefine, indexOptions *options.IndexOptions) (string, error)
	CreateIndexes(indexDefineList []EntityIndexDefine, indexOptions *options.IndexOptions) ([]string, error)
	MustCreateIndex(indexDefine EntityIndexDefine, indexOptions *options.IndexOptions)
	MustCreateIndexes(indexDefineList []EntityIndexDefine, indexOptions *options.IndexOptions)
	DeleteIndex(name string) (err error)
	DeleteAllIndexes() (err error)
	ListIndexes() (indexes []map[string]interface{}, err error)

	GetName() (name string)
	GetCollection() (c *mongo.Collection)
}
