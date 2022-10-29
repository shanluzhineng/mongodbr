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
	FindAll() ([]interface{}, error)
	CountByFilter(interface{}) (int64, error)
	FindByFilter(interface{}, ...FindOption) ([]interface{}, error)
	FindByObjectId(primitive.ObjectID) (interface{}, error)

	Create(interface{}, ...ServiceContextOption) error
	Update(interface{}, ...ServiceContextOption) error
	UpdateFields(objectId primitive.ObjectID, update interface{}, contextOpts ...ServiceContextOption) error
	UpdateMany(filter interface{}, update interface{}) error
	Delete(primitive.ObjectID, ...ServiceContextOption) error
	DeleteMany(interface{}, ...ServiceContextOption) error
}
