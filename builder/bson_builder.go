package builder

import "go.mongodb.org/mongo-driver/bson"

const (
	setKey string = "$set"
)

type BsonBuilder struct {
	bson bson.M
}

func NewBsonBuilder() *BsonBuilder {
	return &BsonBuilder{}
}

// 增加$set类型的值
func (b *BsonBuilder) NewOrUpdateSet(v interface{}) *BsonBuilder {
	b.ensureBson()
	setValue := b.bson[setKey]
	if setValue == nil {
		b.bson = bson.M{"$set": v}
	}
	return b
}

func (b *BsonBuilder) ensureBson() *BsonBuilder {
	if b.bson == nil {
		b.bson = bson.M{}
	}
	return b
}

func (b *BsonBuilder) ToValue() bson.M {
	return b.bson
}
