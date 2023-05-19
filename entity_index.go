package mongodbr

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IEntityIndex interface {
	// index
	CreateIndex(indexModel mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
	CreateIndexes(indexModelList []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)
	MustCreateIndex(indexModel mongo.IndexModel, opts ...*options.CreateIndexesOptions)
	MustCreateIndexes(indexModelList []mongo.IndexModel, opts ...*options.CreateIndexesOptions)
	DeleteIndex(name string) (err error)
	DeleteAllIndexes() (err error)
	ListIndexes() (indexes []map[string]interface{}, err error)
}

var _ IEntityIndex = (*MongoCol)(nil)

// #region indexes members

func (r *MongoCol) CreateIndex(indexModel mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()
	name, err := r.collection.Indexes().CreateOne(ctx, indexModel, opts...)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *MongoCol) CreateIndexes(indexModelList []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	return r.collection.Indexes().CreateMany(ctx, indexModelList, opts...)
}

func (r *MongoCol) MustCreateIndex(indexModel mongo.IndexModel, opts ...*options.CreateIndexesOptions) {
	r.CreateIndex(indexModel, opts...)
}

func (r *MongoCol) MustCreateIndexes(indexModelList []mongo.IndexModel, opts ...*options.CreateIndexesOptions) {
	r.CreateIndexes(indexModelList, opts...)
}

func (r *MongoCol) DeleteIndex(name string) (err error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	_, err = r.collection.Indexes().DropOne(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoCol) DeleteAllIndexes() (err error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	_, err = r.collection.Indexes().DropAll(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoCol) ListIndexes() (indexes []map[string]interface{}, err error) {
	ctx, cancel := CreateContext(r.configuration)
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
