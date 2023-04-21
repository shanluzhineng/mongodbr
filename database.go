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

func Ping() error {
	if DefaultClient == nil {
		return fmt.Errorf("先调用SetupDefaultClient方法创建好一个Client对象后,再调用此方法")
	}
	//测试ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := DefaultClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("mongodb ping测试时出现异常,异常信息:%s", err.Error())
	}
	fmt.Println("mongodb ping测试正常")
	return nil
}
