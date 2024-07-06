package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/shanluzhineng/mongodbr/builder"
)

type EntityUpdate struct {
}
type IEntityBulkWrite interface {
	BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	BulkWriteEntityList(entityList []IEntity, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
}

var _ IEntityBulkWrite = (*MongoCol)(nil)

func _buildWriteModelForUpdate(list []IEntity) []mongo.WriteModel {
	modelList := make([]mongo.WriteModel, 0)
	if len(list) <= 0 {
		return modelList
	}
	for _, eachEntity := range list {
		currentModel := mongo.NewUpdateOneModel()
		currentModel.SetFilter(bson.M{"_id": eachEntity.GetObjectId()})
		currentModel.SetUpdate(builder.NewBsonBuilder().NewOrUpdateSet(eachEntity))
		modelList = append(modelList, currentModel)
	}
	return modelList
}

// build mongo.WriteModel list with ObjectId list
func BuildWriteModelListWithObjectId(dataList map[primitive.ObjectID]interface{}) []mongo.WriteModel {
	modelList := make([]mongo.WriteModel, 0)
	if len(dataList) <= 0 {
		return modelList
	}
	for eachObjectId, eachValue := range dataList {
		currentModel := mongo.NewUpdateOneModel()
		currentModel.SetFilter(bson.M{"_id": eachObjectId})
		currentModel.SetUpdate(builder.NewBsonBuilder().NewOrUpdateSet(eachValue))
		modelList = append(modelList, currentModel)
	}
	return modelList
}

// build mongo.WriteModel list with ObjectId list
func BuildWriteModelList(filterList []interface{}, getUpdateFn func(filter interface{}) interface{}) []mongo.WriteModel {
	modelList := make([]mongo.WriteModel, 0)
	if len(filterList) <= 0 {
		return modelList
	}
	for _, eachFilter := range filterList {
		currentModel := mongo.NewUpdateOneModel()
		currentModel.SetFilter(eachFilter)
		currentModel.SetUpdate(builder.NewBsonBuilder().NewOrUpdateSet(getUpdateFn(eachFilter)))
		modelList = append(modelList, currentModel)
	}
	return modelList
}

// #region update members

func (c *MongoCol) BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (
	*mongo.BulkWriteResult, error) {
	if len(models) <= 0 {
		return nil, nil
	}
	//没有设置参数，使用默认的
	ctx, cancel := CreateContext(c.configuration)
	defer cancel()

	res, err := c.collection.BulkWrite(
		ctx,
		models,
		opts...,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *MongoCol) BulkWriteEntityList(entityList []IEntity, opts ...*options.BulkWriteOptions) (
	*mongo.BulkWriteResult, error) {
	modelList := _buildWriteModelForUpdate(entityList)
	return c.BulkWrite(modelList, opts...)
}

// #endregion
