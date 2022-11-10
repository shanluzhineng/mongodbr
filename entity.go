package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Entity struct {
	ObjectId primitive.ObjectID `json:"objectId,omitempty" bson:"_id"`
}

type IEntity interface {
	GetObjectId() primitive.ObjectID
}

type IEntityBeforeCreate interface {
	BeforeCreate()
}

type IEntityBeforeUpdate interface {
	BeforeUpdate()
}

// 创建时设置对象的基本信息
func (entity *Entity) BeforeCreate() {
	if entity.ObjectId == primitive.NilObjectID {
		entity.ObjectId = primitive.NewObjectID()
	}
}

func (entity *Entity) GetObjectId() primitive.ObjectID {
	return entity.ObjectId
}

type FindOption func(*options.FindOptions)

// 一个抽象的用来处理任意类型的mongodb的仓储基类
type IRepository interface {
	FindAll() (dataList []interface{}, err error)
	CountByFilter(filter interface{}) (count int64, err error)
	FindByFilter(filter interface{}, opts ...FindOption) (dataList []interface{}, err error)
	FindByObjectId(id primitive.ObjectID) (dataList interface{}, err error)

	Create(data interface{}, contextOpts ...ServiceContextOption) error
	Update(data interface{}, contextOpts ...ServiceContextOption) error
	UpdateFields(objectId primitive.ObjectID, update interface{}, contextOpts ...ServiceContextOption) error
	UpdateMany(filter interface{}, update interface{}) error
	DeleteOne(id primitive.ObjectID, contextOpts ...ServiceContextOption) error
	DeleteOneByFilter(filter interface{}, contextOpts ...ServiceContextOption) error
	DeleteMany(filter interface{}, contextOpts ...ServiceContextOption) error
}
