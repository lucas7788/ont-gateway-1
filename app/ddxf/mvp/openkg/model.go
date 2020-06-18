package main

type UserInfo struct {
	UserId string  `bson:"user_id" json:"user_id"`
	OntId string `bson:"ont_id" json:"ont_id"`
}
