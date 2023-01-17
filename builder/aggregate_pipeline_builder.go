package builder

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	op_match = "$match"
	op_group = "$group"
	op_sort  = "$sort"

	op_meta = "$meta"

	field_group_id = "_id"
)

type AggregatePipelineBuilder struct {
	pipeline mongo.Pipeline

	match bson.M
	group bson.M
	sort  bson.M
}

func NewAggregatePipelineBuilder() *AggregatePipelineBuilder {
	builder := &AggregatePipelineBuilder{
		pipeline: make(mongo.Pipeline, 0),

		match: bson.M{},
	}
	builder.pipeline = append(builder.pipeline, bson.D{{
		Key:   op_match,
		Value: builder.match,
	}})
	return builder
}

func (b *AggregatePipelineBuilder) MatchWith(filter bson.M) *AggregatePipelineBuilder {
	if len(filter) <= 0 {
		return b
	}

	//合并match
	for eachKey := range filter {
		b.match[eachKey] = filter[eachKey]
	}
	return b
}

func (b *AggregatePipelineBuilder) SetGroupId(_id string) *AggregatePipelineBuilder {
	b.ensureGroupSetup()
	b.group[field_group_id] = _id
	return b
}

// append group field
func (b *AggregatePipelineBuilder) WithGroupField(fieldName string, value bson.M) *AggregatePipelineBuilder {
	if len(fieldName) <= 0 {
		return b
	}
	b.ensureGroupSetup()
	b.group[fieldName] = value
	return b
}

// append sort field
func (b *AggregatePipelineBuilder) WithSortField(fieldName string, isSortAsc bool, metaDataKeyword string) *AggregatePipelineBuilder {
	b.ensureSortSetup()
	if isSortAsc {
		b.sort[fieldName] = 1
	} else {
		b.sort[fieldName] = -1
	}
	if len(metaDataKeyword) > 0 {
		b.sort[fieldName] = bson.E{
			Key:   op_meta,
			Value: metaDataKeyword,
		}
	}
	return b
}

func (b *AggregatePipelineBuilder) BuildAggregatePipeline() interface{} {
	return b.pipeline
}

func (b *AggregatePipelineBuilder) ensureGroupSetup() {
	if b.group != nil {
		return
	}
	b.group = bson.M{}
	b.pipeline = append(b.pipeline, bson.D{{
		Key:   op_group,
		Value: b.group,
	}})
}

func (b *AggregatePipelineBuilder) ensureSortSetup() {
	if b.sort != nil {
		return
	}
	b.sort = bson.M{}
	b.pipeline = append(b.pipeline, bson.D{{
		Key:   op_sort,
		Value: b.sort,
	}})
}
