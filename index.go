package mongodbr

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// index model
type EntityIndexDefine struct {
	FieldList []IndexFieldDefine
}

func NewEntityIndexDefine() *EntityIndexDefine {
	return &EntityIndexDefine{
		FieldList: make([]IndexFieldDefine, 0),
	}
}

func (d *EntityIndexDefine) AddField(fieldName string, isAsc bool) *EntityIndexDefine {
	d.FieldList = append(d.FieldList, IndexFieldDefine{
		FieldName: fieldName,
		IsAsc:     isAsc,
	})
	return d
}

type IndexFieldDefine struct {
	FieldName string
	IsAsc     bool
}

func (d *EntityIndexDefine) ToIndexModel() *mongo.IndexModel {
	if len(d.FieldList) <= 0 {
		return nil
	}
	keys := bson.D{}
	for _, eachFieldDefine := range d.FieldList {
		keys = append(keys, bson.E{
			Key:   eachFieldDefine.FieldName,
			Value: isAscToIndexValue(eachFieldDefine.IsAsc),
		})
	}

	indexModel := &mongo.IndexModel{
		Keys: keys,
	}
	return indexModel
}

func isAscToIndexValue(isAsc bool) int32 {
	if isAsc {
		return 1
	}
	return -1
}
