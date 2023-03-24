package err

import "go.mongodb.org/mongo-driver/mongo"

func IsDuplicateKeyError(err error) bool {
	// TODO: maybe there is (or will be) a better way of checking duplicate key error
	// this one is based on https://github.com/mongodb/mongo-go-driver/blob/master/mongo/integration/collection_test.go#L54-L65
	we, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}

	return len(we.WriteErrors) > 0 && we.WriteErrors[0].Code == 11000
}
