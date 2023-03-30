package mongodbr

import (
	"context"
	"fmt"

	"github.com/abmpio/mongodbr/builder"
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

func (r *RepositoryBase) CountByFilter(filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// #region find members

func (r *RepositoryBase) FindAll(opts ...FindOption) IFindResult {
	return r.FindByFilter(bson.M{}, opts...)
}

// 根据_id来查找，返回的是对象的指针
func (r *RepositoryBase) FindByObjectId(id primitive.ObjectID) IFindResult {
	return r.FindOne(bson.M{"_id": id})
}

// 查找一条记录
func (r *RepositoryBase) FindOne(filter interface{}, opts ...FindOneOption) IFindResult {
	ctx, cancel := context.WithTimeout(context.Background(), r.configuration.QueryTimeout)
	defer cancel()

	//设置默认搜索参数
	findOneOptions := options.FindOne()
	for _, o := range opts {
		o(findOneOptions)
	}

	res := r.collection.FindOne(ctx, filter, findOneOptions)
	if res.Err() != nil {
		return &findResult{
			configuration: r.configuration,
			err:           res.Err(),
		}
	}
	return &findResult{
		configuration: r.configuration,
		res:           res,
	}
}

// 根据条件来筛选
func (r *RepositoryBase) FindByFilter(filter interface{}, opts ...FindOption) IFindResult {
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
		return &findResult{
			configuration: r.configuration,
			err:           err,
		}
	}
	return &findResult{
		configuration: r.configuration,
		cur:           cur,
	}
}

// #endregion

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

// #region create members

func (r *RepositoryBase) Create(item interface{}, opts ...*options.InsertOneOptions) (id primitive.ObjectID, err error) {
	if item == nil {
		return primitive.NilObjectID, fmt.Errorf("item is nil,col:%s", r.documentName)
	}
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
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
	//没有设置参数，使用默认的
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
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
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
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
	contextProvider := NewDefaultServiceContextProvider()
	ctx := contextProvider.GetContext()
	cancel := contextProvider.GetCancelFunc()
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryBase) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	contextProvider := NewDefaultServiceContextProvider()
	ctx := contextProvider.GetContext()
	cancel := contextProvider.GetCancelFunc()
	defer cancel()

	_, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

// #endregion

func (r *RepositoryBase) ReplaceById(id primitive.ObjectID, doc interface{}, opts ...*options.ReplaceOptions) (err error) {
	return r.Replace(bson.M{"_id": id}, doc, opts...)
}

func (r *RepositoryBase) Replace(filter interface{}, doc interface{}, opts ...*options.ReplaceOptions) (err error) {
	contextProvider := NewDefaultServiceContextProvider()
	ctx := contextProvider.GetContext()
	cancel := contextProvider.GetCancelFunc()
	defer cancel()

	_, err = r.collection.ReplaceOne(ctx, filter, doc, opts...)
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

// #region indexes members

func (r *RepositoryBase) CreateIndex(indexDefine EntityIndexDefine, indexOptions *options.IndexOptions) (string, error) {
	res, err := r.CreateIndexes([]EntityIndexDefine{indexDefine}, indexOptions)
	if err != nil {
		return "", err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return "", nil
}

func (r *RepositoryBase) CreateIndexes(indexDefineList []EntityIndexDefine, indexOptions *options.IndexOptions) ([]string, error) {
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	indexModels := make([]mongo.IndexModel, 0)
	for _, eachIndexDefine := range indexDefineList {
		indexModel := eachIndexDefine.ToIndexModel()
		indexModel.Options = indexOptions
	}
	return r.collection.Indexes().CreateMany(ctx, indexModels)
}

func (r *RepositoryBase) MustCreateIndex(indexDefine EntityIndexDefine, indexOptions *options.IndexOptions) {
	r.CreateIndex(indexDefine, indexOptions)
}

func (r *RepositoryBase) MustCreateIndexes(indexDefineList []EntityIndexDefine, indexOptions *options.IndexOptions) {
	r.CreateIndexes(indexDefineList, indexOptions)
}

func (r *RepositoryBase) DeleteIndex(name string) (err error) {
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	_, err = r.collection.Indexes().DropOne(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryBase) DeleteAllIndexes() (err error) {
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	_, err = r.collection.Indexes().DropAll(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryBase) ListIndexes() (indexes []map[string]interface{}, err error) {
	contextOpts := WithDefaultServiceContext()
	ctx := contextOpts().GetContext()
	cancel := contextOpts().GetCancelFunc()
	defer cancel()

	cur, err := r.collection.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &indexes); err != nil {
		return nil, err
	}
	return indexes, nil
}

// #endregion

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
