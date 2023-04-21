package mongodbr

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewRepositoryOption struct {
	databaseName     string
	collectionName   string
	DefaultSortField string
}

func newDefaultRepositoryOption() *NewRepositoryOption {
	o := &NewRepositoryOption{}
	return o
}

func NewRepository(databaseName string, collectionName string, opts ...func(*NewRepositoryOption)) (*RepositoryBase, error) {
	if len(databaseName) <= 0 {
		err := fmt.Errorf("database参数不能为nil")
		return nil, err
	}
	if len(collectionName) <= 0 {
		err := fmt.Errorf("collectionName参数不能为nil")
		return nil, err
	}
	o := newDefaultRepositoryOption()
	o.databaseName = databaseName
	o.collectionName = collectionName
	for _, eachOpt := range opts {
		eachOpt(o)
	}
	collection := GetCollection(databaseName, collectionName)
	mongodbrOpts := make([]RepositoryOption, 0)
	if len(o.DefaultSortField) > 0 {
		mongodbrOpts = append(mongodbrOpts, WithDefaultSort(func(fo *options.FindOptions) *options.FindOptions {
			return fo.SetSort(bson.D{{Key: o.DefaultSortField, Value: -1}})
		}))
	}
	repositoryBase, err := NewRepositoryBase(func() *mongo.Collection {
		return collection
	}, mongodbrOpts...)
	if err != nil {
		return nil, err
	}
	return repositoryBase, nil
}
