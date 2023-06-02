package mongodbr

import (
	"errors"
	"fmt"

	"github.com/abmpio/mongodbr/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCol struct {
	configuration *Configuration
	collection    *mongo.Collection
}

// new MongoCol instance, panic if col is nil
func NewMongoCol(col *mongo.Collection, opts ...*Configuration) *MongoCol {
	if col == nil {
		panic(errors.New("col cannot be nil"))
	}
	c := NewConfiguration()
	if len(opts) > 0 {
		c = opts[0]
	}
	mongoCol := &MongoCol{
		configuration: c,
		collection:    col,
	}
	return mongoCol
}

// RepositoryBase represents a mongodb repository
type RepositoryBase struct {
	documentName string
	*MongoCol
}

var _ IRepository = (*RepositoryBase)(nil)

// new一个新的实例
func NewRepositoryBase(getDbCollection func() *mongo.Collection, opts ...RepositoryOption) (*RepositoryBase, error) {
	if getDbCollection == nil {
		err := fmt.Errorf("getDbCollection参数不能为nil")
		return nil, err
	}
	coll := getDbCollection()
	repository := &RepositoryBase{
		MongoCol:     NewMongoCol(coll),
		documentName: coll.Name(),
	}
	for _, eachItem := range opts {
		eachItem(repository.configuration)
	}
	return repository, nil
}

// aggregate
func (r *RepositoryBase) Aggregate(pipeline interface{}, dataList interface{}, opts ...AggregateOption) (err error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	//设置默认搜索参数
	aggregateOptions := options.Aggregate()
	for _, o := range opts {
		o(aggregateOptions)
	}
	cur, err := r.collection.Aggregate(ctx, pipeline, aggregateOptions)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, dataList)
}

// #region create members

func (r *RepositoryBase) Create(item interface{}, opts ...*options.InsertOneOptions) (id primitive.ObjectID, err error) {
	if item == nil {
		return primitive.NilObjectID, fmt.Errorf("item is nil,col:%s", r.documentName)
	}
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	r.onBeforeCreate(item)
	res, err := r.collection.InsertOne(ctx, item, opts...)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		return id, nil
	}
	return primitive.NilObjectID, ErrInvalidType
}

func (r *RepositoryBase) CreateMany(itemList []interface{}, opts ...*options.InsertManyOptions) (ids []primitive.ObjectID, err error) {
	if len(itemList) <= 0 {
		return nil, nil
	}
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	for index := range itemList {
		r.onBeforeCreate(itemList[index])
	}
	res, err := r.collection.InsertMany(ctx, itemList, opts...)
	if err != nil {
		return nil, err
	}
	for _, v := range res.InsertedIDs {
		switch v := v.(type) {
		case primitive.ObjectID:
			ids = append(ids, v)
		default:
			return nil, ErrInvalidType
		}
	}
	return ids, nil
}

// #endregion

// #region update members

func (r *RepositoryBase) FindOneAndUpdate(entity IEntity, opts ...*options.FindOneAndUpdateOptions) error {
	if entity == nil {
		return fmt.Errorf("在更新%s数据时item参数不能为nil", r.documentName)
	}

	objectId := entity.GetObjectId()
	update := builder.NewBsonBuilder().NewOrUpdateSet(entity).ToValue()
	return r.FindOneAndUpdateWithId(objectId, update, opts...)
}

func (r *RepositoryBase) FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	if objectId.IsZero() {
		return fmt.Errorf("在保存%s数据时objectId不能为nil", r.documentName)
	}
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

func (r *RepositoryBase) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryBase) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (interface{}, error) {
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

func (r *RepositoryBase) ReplaceById(id primitive.ObjectID, doc interface{}, opts ...*options.ReplaceOptions) (err error) {
	return r.Replace(bson.M{"_id": id}, doc, opts...)
}

func (r *RepositoryBase) Replace(filter interface{}, doc interface{}, opts ...*options.ReplaceOptions) (err error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	_, err = r.collection.ReplaceOne(ctx, filter, doc, opts...)
	if err != nil {
		return err
	}
	return nil
}

// 删除指定id的记录
func (r *RepositoryBase) DeleteOne(id primitive.ObjectID, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return result, err
	}

	return result, nil
}

// 删除指定条件的一条记录
func (r *RepositoryBase) DeleteOneByFilter(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return result, err
	}

	return result, nil
}

// 删除多条记录
func (r *RepositoryBase) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if filter == nil {
		err := fmt.Errorf("无法删除多条%s记录,filter参数不能为null", r.documentName)
		return nil, err
	}
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	result, err := r.collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *RepositoryBase) GetName() (name string) {
	return r.documentName
}

func (r *RepositoryBase) GetCollection() (c *mongo.Collection) {
	return r.collection
}

func (r *RepositoryBase) onBeforeCreate(item interface{}) {
	entityHookable, ok := item.(IEntityBeforeCreate)
	if !ok {
		return
	}
	entityHookable.BeforeCreate()
}

// func (r *RepositoryBase) onBeforeUpdate(item interface{}) {
// 	entityHookable, ok := item.(IEntityBeforeUpdate)
// 	if !ok {
// 		return
// 	}
// 	entityHookable.BeforeUpdate()
// }
