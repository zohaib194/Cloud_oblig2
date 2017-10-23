package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Webhook mongodb stores the details of the DB connection.
type WebhookMongoDB struct {
	DatabaseURL       string
	DatabaseName      string
	WebhookCollection string
}

type Payload struct {
	ID              bson.ObjectId `json:"_id" bson:"_id"`
	WebhookURL      string        `json:"webhookurl" bson:"webhookurl"`
	BaseCurrency    string        `json:"basecurrency" bson:"basecurrency"`
	TargetCurrency  string        `json:"targetcurrency" bson:"targetcurrency"`
	MinTriggerValue float64       `json:"mintriggervalue" bson:"mintriggervalue"`
	MaxTriggerValue float64       `json:"maxtriggervalue" bson:"maxtriggervalue"`
}

type Payload_I struct {
	BaseCurrency    string
	TargetCurrency  string
	CurrentRate     string
	MinTriggerValue float64
	MaxTriggerValue float64
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
Add adds new payload to the storage.
*/
func (db *WebhookMongoDB) Add(p Payload) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	i := bson.NewObjectId()
	//fmt.Println(l)
	p.ID = i

	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Insert(p)
	if err != nil {
		fmt.Printf("error in Insert(), %v", err.Error())
	}

}

/*
Get the unique id of a given webhook from the storage.
*/
func (db *WebhookMongoDB) Get(keyId string) (string, bool) {
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
		l := tempP.ID.Hex()
		return l, false
	}
	l := tempP.ID.Hex()
	return l, true
}

func postReqHandler(w http.ResponseWriter, r *http.Request) {
	db := WebhookMongoDB{
		DatabaseURL:       "mongodb://localhost",
		DatabaseName:      "Webhook",
		WebhookCollection: "WebhookPayload",
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var p Payload

	err = json.Unmarshal(body, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if len(p.WebhookURL) != 0 && len(p.BaseCurrency) != 0 && len(p.TargetCurrency) != 0 && p.MaxTriggerValue > 0 && p.MinTriggerValue > 0 {
		db.Init()
		db.Add(p)
		var s string
		s, ok := db.Get(p.WebhookURL)

		if !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//fmt.Printf(string(body))
		//fmt.Println(s)
		fmt.Fprint(w, s)
	} else {
		http.Error(w, "Post request body is not correctly formed", http.StatusBadRequest)
	}
}

func registeredWebhook(w http.ResponseWriter, r *http.Request) {
	//Get the webhookUrl

	//Check if it exist in db

	//Update the payload with currect rate

	//Send notification to the webhookurl

	//if recieved the 200 or 204 status code

	//else check errors
}

func main() {
	http.HandleFunc("/root", postReqHandler)
	http.HandleFunc("/root/{id}", registeredWebhook)

	http.ListenAndServe(":8080", nil)
}
