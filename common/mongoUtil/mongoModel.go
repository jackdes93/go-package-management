package mongoUtil

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoModel struct {
	PK        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status    int                `json:"status" bson:"status,omitempty"`
	CreatedAt *JSONTime          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt *JSONTime          `json:"updatedAt" bson:"updatedAt,omitempty"`
	DeletedAt *JSONTime          `json:"deletedAt" bson:"deletedAt,omitempty"`
}

func (mg *MongoModel) PrepareForInsert(status int) {
	jsTime := JSONTime(time.Now().UTC())
	mg.PK = primitive.NewObjectID()
	mg.CreatedAt = &jsTime
	mg.UpdatedAt = &jsTime
	mg.Status = status
}
