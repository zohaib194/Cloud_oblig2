package database

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Id (used temporary various places in the code to store _id of a certain webhook payload)
type Id struct {
	ID bson.ObjectId `bson:"_id"`
}

// Webhook mongodb stores the details of the DB connection.
type WebhookMongoDB struct {
	DatabaseURL  string
	DatabaseName string
	Collection   string
}

/*
Init initializes the mongo storage.
*/
func (db *WebhookMongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//TODO put extra constraints on the webhook collection

}

/*
Add adds new Subscriber to the storage.
*/
func (db *WebhookMongoDB) Add(p Subscriber) (string, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var id Id
	err = session.DB(db.DatabaseName).C(db.Collection).Insert(p)
	session.DB(db.DatabaseName).C(db.Collection).Find(bson.M{"webhookurl": p.WebhookURL}).One(&id)
	l := id.ID.Hex()

	if err != nil {
		fmt.Printf("error in Insert(), %v", err.Error())
		return l, false
	}
	return l, true

}

/*
Get gets the unique id of a given webhook from the storage.
*/
func (db *WebhookMongoDB) Get(keyID string) (Subscriber, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	tempP := Subscriber{}

	id := bson.ObjectIdHex(keyID)
	err = session.DB(db.DatabaseName).C(db.Collection).Find(bson.M{"_id": id}).One(&tempP)

	if err != nil {
		fmt.Printf("err in Get(), %v", err.Error())
		return tempP, false
	}
	return tempP, true
}

/*
Delete remove a payload from the DB referenced to the keyID
*/
func (db *WebhookMongoDB) Delete(keyID string) bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	id := bson.ObjectIdHex(keyID)
	err = session.DB(db.DatabaseName).C(db.Collection).Remove(bson.M{"_id": id})

	if err != nil {
		fmt.Printf("err in Delete(), %v", err.Error())
		return false
	}
	return true
}

func (db *WebhookMongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.Collection).Count()
	if err != nil {
		fmt.Printf("err in Count(), %v", err.Error())
		return -1
	}
	return count
}
