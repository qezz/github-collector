package models

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	Objid    bson.ObjectId `bson:"_id,omitempty"`
	Id       int64
	Username string
	Fullname string
	Location string
}

func NewUser(id int64, username, fullname, location string) User {
	return User{
		Id:       id,
		Username: username,
		Fullname: fullname,
		Location: location,
	}
}

func (u User) String() string {
	return fmt.Sprintf("User { \"%v\", \"%v\", \"%v\", \"%v\" }",
		u.Id, u.Username, u.Fullname, u.Location)
}
