package mongodbr

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RepositoryBase represents a mongodb repository
type RepositoryBase struct {
	configuration *Configuration
	documentName  string
	collection    *mongo.Collection
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
		collection:    coll,
		documentName:  coll.Name(),
		configuration: NewConfiguration(),
	}
	for _, eachItem := range opts {
		eachItem(repository.configuration)
	}
	return repository, nil
}

func (r *RepositoryBase) FindAll() ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()

	findOptions := options.Find()
	if r.configuration.setDefaultSort != nil {
		r.configuration.setDefaultSort(findOptions)
	}
	cur, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []interface{}
	for cur.Next(ctx) {
		o := r.configuration.createItemFunc()
		if err := cur.Decode(o); err != nil {
			return nil, err
		}

		result = append(result, o)
	}

	return result, cur.Err()
}

func (r *RepositoryBase) CountByFilter(filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// 查找一条记录
func (r *RepositoryBase) FindOne(filter interface{}, opts ...FindOneOption) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()

	//设置默认搜索参数
	findOneOptions := options.FindOne()
	for _, o := range opts {
		o(findOneOptions)
	}

	result := r.configuration.createItemFunc()
	err := r.collection.FindOne(ctx, filter, findOneOptions).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, err

}

// 根据条件来筛选
func (r *RepositoryBase) FindByFilter(filter interface{}, opts ...FindOption) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()

	//设置默认搜索参数
	findOptions := options.Find()
	if r.configuration.setDefaultSort != nil {
		r.configuration.setDefaultSort(findOptions)
	}
	for _, o := range opts {
		o(findOptions)
	}
	cur, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []interface{}
	for cur.Next(ctx) {
		o := r.configuration.createItemFunc()
		if err := cur.Decode(o); err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	return result, cur.Err()
}

// 根据_id来查找，返回的是对象的指针
func (r *RepositoryBase) FindByObjectId(id primitive.ObjectID) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()

	result := r.configuration.createItemFunc()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return result, err
}

// aggregate
func (r *RepositoryBase) Aggregate(pipeline interface{}, dataList interface{}, opts ...AggregateOption) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
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

func (r *RepositoryBase) Create(item interface{}, opts ...*options.InsertOneOptions) error {
	if item == nil {
		return fmt.Errorf("在插入%s数据时item参数不能为nil", r.documentName)
	}
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	r.onBeforeCreate(item)
	_, err := r.collection.InsertOne(ctx, item, opts...)
	if err != nil {
		if r.isDuplicateKeyError(err) {
			return fmt.Errorf("%s中已经存在着相同的记录", r.documentName)
		}
		return err
	}
	return nil
}

// create many
func (r *RepositoryBase) CreateMany(itemList []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if len(itemList) <= 0 {
		return &mongo.InsertManyResult{}, nil
	}
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	for index := range itemList {
		r.onBeforeCreate(itemList[index])
	}
	return r.collection.InsertMany(ctx, itemList, opts...)
}

func (r *RepositoryBase) FindOneAndUpdateEntityWithId(entity interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	if entity == nil {
		return fmt.Errorf("在更新%s数据时item参数不能为nil", r.documentName)
	}

	value, ok := entity.(IEntity)
	if !ok {
		return fmt.Errorf("entity必须实现IEntity接口")
	}
	objectId := value.GetObjectId()
	update := NewBsonBuilder().NewOrUpdateSet(entity).ToValue()
	return r.FindOneAndUpdateWithId(objectId, update, opts...)
}

func (r *RepositoryBase) FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	if objectId.IsZero() {
		return fmt.Errorf("在保存%s数据时objectId不能为nil", r.documentName)
	}
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	if len(opts) <= 0 {
		opts = make([]*options.FindOneAndUpdateOptions, 0)
		opts = append(opts, options.FindOneAndUpdate().SetUpsert(true))
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

func (r *RepositoryBase) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	if update == nil {
		return fmt.Errorf("在保存%s数据时update参数不能为nil", r.documentName)
	}
	contextProvider := NewDefaultServiceContextProvider()
	ctx := contextProvider.GetContext()
	cancel := contextProvider.GetCancelFunc()
	defer cancel()

	// updateValue := bson.M{"$set": update}
	_, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

// 删除指定id的记录
func (r *RepositoryBase) DeleteOne(id primitive.ObjectID, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return result, err
	}

	return result, nil
}

// 删除指定条件的一条记录
func (r *RepositoryBase) DeleteOneByFilter(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if filter == nil {
		err := fmt.Errorf("filter参数不能为null")
		return nil, err
	}
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
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
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	result, err := r.collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *RepositoryBase) CreateIndexIfNotExist(indexDefine EntityIndexDefine, indexOptions *options.IndexOptions) (string, error) {
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	indexModel := indexDefine.ToIndexModel()
	indexModel.Options = indexOptions
	return r.collection.Indexes().CreateOne(ctx, *indexModel)
}

func (r *RepositoryBase) onBeforeCreate(item interface{}) {
	entityHookable, ok := item.(IEntityBeforeCreate)
	if !ok {
		return
	}
	entityHookable.BeforeCreate()
}

func (r *RepositoryBase) onBeforeUpdate(item interface{}) {
	entityHookable, ok := item.(IEntityBeforeUpdate)
	if !ok {
		return
	}
	entityHookable.BeforeUpdate()
}

func (r *RepositoryBase) isDuplicateKeyError(err error) bool {
	// TODO: maybe there is (or will be) a better way of checking duplicate key error
	// this one is based on https://github.com/mongodb/mongo-go-driver/blob/master/mongo/integration/collection_test.go#L54-L65
	we, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}

	return len(we.WriteErrors) > 0 && we.WriteErrors[0].Code == 11000
}
