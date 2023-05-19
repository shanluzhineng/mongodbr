package mongodbr

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type IFindResult interface {
	One(val interface{}) (err error)
	ToOne() (interface{}, error)
	All(val interface{}) (err error)
	ToAll() ([]interface{}, error)

	GetSingleResult() (res *mongo.SingleResult)
	GetCursor() (cur *mongo.Cursor)
	GetError() (err error)
}

type findResult struct {
	res           *mongo.SingleResult
	cur           *mongo.Cursor
	err           error
	configuration *Configuration
}

// #IFindResult members

func (r *findResult) One(val interface{}) (err error) {
	if r.err != nil {
		return r.err
	}
	if r.cur == nil {
		return r.res.Decode(val)
	}

	//没有设置参数，使用默认的
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	if !r.cur.TryNext(ctx) {
		return mongo.ErrNoDocuments
	}
	return r.cur.Decode(val)
}

func (r *findResult) ToOne() (interface{}, error) {
	result := r.configuration.safeCreateItem()
	err := r.One(result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return result, nil
}

func (r *findResult) All(val interface{}) (err error) {
	if r.err != nil {
		return r.err
	}

	ctx, cancel := CreateContext(r.configuration)
	defer cancel()
	if r.cur == nil {
		return ErrNoCursor
	}
	if !r.cur.TryNext(ctx) {
		return ctx.Err()
	}
	return r.cur.All(ctx, val)
}

func (r *findResult) ToAll() ([]interface{}, error) {
	if r.err != nil {
		return nil, nil
	}
	if r.cur == nil {
		return nil, nil
	}
	//没有设置参数，使用默认的
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()
	defer r.cur.Close(ctx)

	var result []interface{}
	for r.cur.Next(ctx) {
		o := r.configuration.safeCreateItem()
		if err := r.cur.Decode(o); err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, nil
}

func (r *findResult) GetSingleResult() (res *mongo.SingleResult) {
	return r.res
}

func (r *findResult) GetCursor() (cur *mongo.Cursor) {
	return r.cur
}

func (r *findResult) GetError() (err error) {
	return r.err
}

// #endregion
