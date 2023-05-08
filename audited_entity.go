package mongodbr

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// can audit creation entity
type CreationAuditedEntity struct {
	Entity `bson:",inline"`
	//create time
	CreationTime time.Time `json:"creationTime,omitempty" bson:"creationTime" `
	//create user
	CreatorId string `json:"creatorId,omitempty" bson:"creatorId"`
}

// auditable entity
type AuditedEntity struct {
	CreationAuditedEntity `bson:",inline"`
	//last modification time
	LastModificationTime *time.Time `json:"lastModificationTime,omitempty" bson:"lastModificationTime"`
	//last modification user
	LastModifierId string `json:"lastModifierId,omitempty" bson:"lastModifierId"`
}

func (e AuditedEntity) GetObjectId() primitive.ObjectID {
	return e.ObjectId
}

func (entity *CreationAuditedEntity) BeforeCreate() {
	entity.Entity.BeforeCreate()
	if entity.CreationTime.IsZero() {
		entity.CreationTime = time.Now()
	}
}

func (entity *AuditedEntity) BeforeUpdate() {
	now := time.Now()
	entity.LastModificationTime = &now
}
