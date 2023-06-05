package mongodbr

import (
	"github.com/abmpio/mongodbr/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// update
type IEntityUpdate interface {
	FindOneAndUpdate(entity IEntity, opts ...*options.FindOneAndUpdateOptions) error
	FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error
	UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (interface{}, error)
}

var _ IEntityUpdate = (*MongoCol)(nil)

// #region update members

func (r *MongoCol) FindOneAndUpdate(entity IEntity, opts ...*options.FindOneAndUpdateOptions) error {
	// if entity == nil {
	// 	return fmt.Errorf("在更新%s数据时item参数不能为nil", r.documentName)
	// }

	objectId := entity.GetObjectId()
	update := builder.NewBsonBuilder().NewOrUpdateSet(entity).ToValue()
	return r.FindOneAndUpdateWithId(objectId, update, opts...)
}

func (r *MongoCol) FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	// if objectId.IsZero() {
	// 	return fmt.Errorf("在保存%s数据时objectId不能为nil", r.documentName)
	// }
	//没有设置参数，使用默认的
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	if len(opts) <= 0 {
		opts = make([]*options.FindOneAndUpdateOptions, 0)
		opts = append(opts, options.FindOneAndUpdate().SetUpsert(false))
	}
	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectId},
		update,
		opts...,
	).Err(); err != nil {
		return err
	}

	return nil
}

func (r *MongoCol) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoCol) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (interface{}, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	result, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		if result != nil {
			return result.UpsertedID, err
		} else {
			return nil, err
		}
	}

	if result != nil {
		return result.UpsertedID, nil
	}
	return nil, nil
}

// #endregion
