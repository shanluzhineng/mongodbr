package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IEntityFind interface {
	CountByFilter(filter interface{}) (count int64, err error)
	CountAll() (count int64, err error)

	// find
	FindAll(opts ...FindOption) IFindResult
	FindByObjectId(id primitive.ObjectID) IFindResult
	FindOne(filter interface{}, opts ...FindOneOption) IFindResult
	FindByFilter(filter interface{}, opts ...FindOption) IFindResult

	Distinct(fieldName string, filter interface{}) ([]interface{}, error)
}

var _ IEntityFind = (*MongoCol)(nil)

func (r *MongoCol) CountByFilter(filter interface{}) (int64, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *MongoCol) CountAll() (count int64, err error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()
	total, err := r.collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *MongoCol) FindAll(opts ...FindOption) IFindResult {
	return r.FindByFilter(bson.M{}, opts...)
}

// 根据_id来查找，返回的是对象的指针
func (r *MongoCol) FindByObjectId(id primitive.ObjectID) IFindResult {
	return r.FindOne(bson.M{"_id": id})
}

// 查找一条记录
func (r *MongoCol) FindOne(filter interface{}, opts ...FindOneOption) IFindResult {
	ctx, cancel := CreateContext(r.configuration)
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
func (r *MongoCol) FindByFilter(filter interface{}, opts ...FindOption) IFindResult {
	ctx, cancel := CreateContext(r.configuration)
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

func (r *MongoCol) Distinct(fieldName string, filter interface{}) ([]interface{}, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	return r.collection.Distinct(ctx, fieldName, filter)
}
