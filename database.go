package mongodbr

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// get mongo.Database instance
func GetDatabase(databaseName string, opts ...*options.DatabaseOptions) *mongo.Database {
	if DefaultClient == nil {
		return nil
	}
	if len(databaseName) <= 0 {
		return nil
	}
	return DefaultClient.Database(databaseName, opts...)
}

// get mongo.Database instance
func GetDatabaseByKey(key string, databaseName string, opts ...*options.DatabaseOptions) *mongo.Database {
	client := GetClient(key)
	if client == nil {
		return nil
	}
	if len(databaseName) <= 0 {
		return nil
	}
	return client.Database(databaseName, opts...)
}

// get mongo.Collection instanc
func GetCollection(databaseName string, collectionName string, opts ...*options.CollectionOptions) *mongo.Collection {
	database := GetDatabase(databaseName)
	if database == nil {
		return nil
	}
	if len(collectionName) <= 0 {
		return nil
	}
	return database.Collection(collectionName, opts...)
}

func GetCollectionByKey(key string, databaseName string, collectionName string, opts ...*options.CollectionOptions) *mongo.Collection {
	database := GetDatabaseByKey(key, databaseName)
	if database == nil {
		return nil
	}
	if len(collectionName) <= 0 {
		return nil
	}
	return database.Collection(collectionName, opts...)
}

func Ping(client *mongo.Client) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	//测试ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("mongodb ping测试时出现异常,异常信息:%s", err.Error())
	}
	fmt.Println("mongodb ping测试正常")
	return nil
}
