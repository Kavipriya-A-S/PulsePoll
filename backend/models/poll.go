package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Poll struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string        `bson:"question" json:"question"`
	Options   []PollOption  `bson:"options" json:"options"`
	CreatedAt string        `bson:"createdAt" json:"createdAt"`
}

type PollOption struct {
	ID    string `bson:"id" json:"id"`
	Text  string `bson:"text" json:"text"`
	Votes int    `bson:"votes" json:"votes"`
}