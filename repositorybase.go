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

func (r *RepositoryBase) Create(item interface{}, contextOpts ...ServiceContextOption) error {
	if item == nil {
		return fmt.Errorf("在插入%s数据时item参数不能为nil", r.documentName)
	}
	if len(contextOpts) <= 0 {
		//没有设置参数，使用默认的
		contextOpts = []ServiceContextOption{WithDefaultServiceContext()}
	}
	ctx := contextOpts[0]().GetContext()
	cancel := contextOpts[0]().GetCancelFunc()
	defer cancel()

	r.onBeforeCreate(item)
	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		if r.isDuplicateKeyError(err) {
			return fmt.Errorf("%s中已经存在着相同的记录", r.documentName)
		}
		return err
	}
	return nil
}

func (r *RepositoryBase) Update(item interface{}, contextOpts ...ServiceContextOption) error {
	if item == nil {
		return fmt.Errorf("在更新%s数据时item参数不能为nil", r.documentName)
	}

	if len(contextOpts) <= 0 {
		//没有设置参数，使用默认的
		contextOpts = []ServiceContextOption{WithDefaultServiceContext()}
	}
	ctx := contextOpts[0]().GetContext()
	cancel := contextOpts[0]().GetCancelFunc()
	defer cancel()

	objectId := item.(IEntity).GetObjectId()
	r.onBeforeUpdate(item)
	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectId},
		bson.M{"$set": item},
		options.FindOneAndUpdate().SetUpsert(true),
	).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryBase) UpdateFields(objectId primitive.ObjectID, update interface{}, contextOpts ...ServiceContextOption) error {
	if objectId.IsZero() {
		return fmt.Errorf("在保存%s数据时objectId不能为nil", r.documentName)
	}
	if len(contextOpts) <= 0 {
		//没有设置参数，使用默认的
		contextOpts = []ServiceContextOption{WithDefaultServiceContext()}
	}
	ctx := contextOpts[0]().GetContext()
	cancel := contextOpts[0]().GetCancelFunc()
	defer cancel()

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectId},
		bson.M{"$set": update},
	).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RepositoryBase) UpdateMany(filter interface{}, update interface{}) error {
	if update == nil {
		return fmt.Errorf("在保存%s数据时update参数不能为nil", r.documentName)
	}
	contextProvider := NewDefaultServiceContextProvider()
	ctx := contextProvider.GetContext()
	cancel := contextProvider.GetCancelFunc()
	defer cancel()

	updateValue := bson.M{"$set": update}
	_, err := r.collection.UpdateMany(ctx, filter, updateValue)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryBase) Delete(id primitive.ObjectID, contextOpts ...ServiceContextOption) error {
	if len(contextOpts) <= 0 {
		//没有设置参数，使用默认的
		contextOpts = []ServiceContextOption{WithDefaultServiceContext()}
	}
	ctx := contextOpts[0]().GetContext()
	cancel := contextOpts[0]().GetCancelFunc()
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryBase) DeleteMany(filter interface{}, contextOpts ...ServiceContextOption) error {
	if filter == nil {
		err := fmt.Errorf("无法删除多条%s记录,filter参数不能为null", r.documentName)
		return err
	}
	if len(contextOpts) <= 0 {
		//没有设置参数，使用默认的
		contextOpts = []ServiceContextOption{WithDefaultServiceContext()}
	}
	ctx := contextOpts[0]().GetContext()
	cancel := contextOpts[0]().GetCancelFunc()
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
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
