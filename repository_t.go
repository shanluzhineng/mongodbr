package mongodbr

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// find all t
func FindAllT[T any](repository IRepository, opts ...FindOption) ([]T, error) {
	res := repository.FindAll()
	list := make([]T, 0)
	if err := res.All(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// find t by filter
func FindTByFilter[T any](repository IRepository, filter interface{}, opts ...FindOption) ([]T, error) {
	res := repository.FindByFilter(filter, opts...)
	list := make([]T, 0)
	if err := res.All(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// find t by _id
func FindTByObjectId[T any](repository IRepository, id primitive.ObjectID) (*T, error) {
	res := repository.FindByObjectId(id)
	result := new(T)
	if err := res.One(result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// find one by _id
func FindOneTByFilter[T any](repository IRepository, filter interface{}, opts ...FindOneOption) (*T, error) {
	res := repository.FindOne(filter, opts...)
	result := new(T)
	if err := res.One(result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
