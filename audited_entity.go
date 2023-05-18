package mongodbr

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICreationAuditedEntity interface {
	GetCreatorId() string
	GetCreationTime() time.Time
}

var _ ICreationAuditedEntity = (*CreationAuditedEntity)(nil)

// can audit creation entity
type CreationAuditedEntity struct {
	Entity `bson:",inline"`
	//create time
	CreationTime time.Time `json:"creationTime,omitempty" bson:"creationTime" `
	//create user
	CreatorId string `json:"creatorId,omitempty" bson:"creatorId"`
}

// #region ICreationAuditedEntity Members

func (e *CreationAuditedEntity) GetCreatorId() string {
	return e.CreatorId
}

func (e *CreationAuditedEntity) GetCreationTime() time.Time {
	return e.CreationTime
}

// #endregion

type IModificationEntity interface {
	GetLastModificationTime() *time.Time
	GetLastModifierId() string
}

var _ IModificationEntity = (*AuditedEntity)(nil)

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

// #region IModificationEntity Members

func (e *AuditedEntity) GetLastModificationTime() *time.Time {
	return e.LastModificationTime
}

func (e *AuditedEntity) GetLastModifierId() string {
	return e.CreatorId
}

// #endregion
