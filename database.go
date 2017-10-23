package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

/*
// Webhook mongodb stores the details of the DB connection.
type WebhookMongoDB struct {
	DatabaseURL       string
	DatabaseName      string
	WebhookCollection string
}

/*
Init initializes the mongo storage.

func (db *WebhookMongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//TODO put extra constraints on the webhook collection

}

/*
Add adds new payload to the storage.

func (db *WebhookMongoDB) Add(p Payload) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Insert(p)
	if err != nil {
		fmt.Printf("error in Insert(), %v", err.Error())
	}
}
*/
func (db *WebhookMongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.WebhookCollection).Count()
	if err != nil {
		fmt.Printf("err in Count(), %v", err.Error())
		return -1
	}
	return count
}

/*
func (db *WebhookMongoDB) Get(keyId string) (Payload, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	tempP := Payload{}

	//check the query

	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Find(bson.M{"webhookurl": keyId}).One(&tempP)

	if err != nil {
		fmt.Printf("err in Get(), %v", err.Error())
		return tempP, false
	}
	return tempP, true
}
*/
