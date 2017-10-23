package main

import "testing"
import "gopkg.in/mgo.v2"

func setupDB(t *testing.T) *WebhookMongoDB {
	db := WebhookMongoDB{
		DatabaseURL:       "mongodb://localhost",
		DatabaseName:      "testPayload",
		WebhookCollection: "payload",
	}

	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	if err != nil {
		t.Error(err)
	}
	return &db
}

func dropDB(t *testing.T, db *WebhookMongoDB) {
	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	if err != nil {
		t.Error(err)
	}
	err = session.DB(db.DatabaseName).DropDatabase()
	if err != nil {
		t.Error(err)
	}

}

func TestPayloadMongoDB_Add(t *testing.T) {
	db := setupDB(t)
	defer dropDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("Database not properly initialized, payload count should be 0")
	}
	payload := Payload{
		WebhookURL:      "http://remoteUrl:8080/randomWebhookPath",
		BaseCurrency:    "EUR",
		TargetCurrency:  "NOK",
		MinTriggerValue: 1.50,
		MaxTriggerValue: 2.55,
	}
	db.Add(payload)
	if db.Count() != 1 {
		t.Error("Adding new payload failed.")
	}
}

func TestPayloadMongoDB_Get(t *testing.T) {
	db := setupDB(t)
	defer dropDB(t, db)

	db.Init()
	if db.Count() != 0 {
		t.Error("Database not properly initialized, payload count should be 0")
	}
	payload := Payload{
		WebhookURL:      "http://remoteUrl:8080/randomWebhookPath",
		BaseCurrency:    "EUR",
		TargetCurrency:  "NOK",
		MinTriggerValue: 1.50,
		MaxTriggerValue: 2.55,
	}
	db.Add(payload)
	if db.Count() != 1 {
		t.Error("Adding new payload failed.")
	}

	newPayload, ok := db.Get(payload.WebhookURL)
	if !ok {
		t.Error("couldn't find " + payload.WebhookURL)
	}

	if newPayload.WebhookURL != payload.WebhookURL ||
		newPayload.BaseCurrency != payload.BaseCurrency ||
		newPayload.TargetCurrency != payload.TargetCurrency ||
		newPayload.MaxTriggerValue != payload.MaxTriggerValue ||
		newPayload.MinTriggerValue != payload.MinTriggerValue {
		t.Error("payload do not match")

	}
}
