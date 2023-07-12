package mongodbr

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewRepositoryOption struct {
	clientKey        string
	databaseName     string
	collectionName   string
	DefaultSortField string
}

func newDefaultRepositoryOption() *NewRepositoryOption {
	o := &NewRepositoryOption{}
	return o
}

// specifiy repository with client key
func RepositoryOptionWithClientKey(clientKey string) func(*NewRepositoryOption) {
	return func(nro *NewRepositoryOption) {
		nro.clientKey = clientKey
	}
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
	var collection *mongo.Collection
	if len(o.clientKey) <= 0 {
		collection = GetCollection(o.databaseName, o.collectionName)
	} else {
		collection = GetCollectionByKey(o.clientKey, o.databaseName, o.collectionName)
	}
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
