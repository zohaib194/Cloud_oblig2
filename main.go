package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Webhook mongodb stores the details of the DB connection.
type WebhookMongoDB struct {
	DatabaseURL       string
	DatabaseName      string
	WebhookCollection string
}

type Subscriber struct {
	//ID              bson.ObjectId `json:"_id, omitempty" bson:"_id"`
	WebhookURL      string  `json:"webhookurl" bson:"webhookurl"`
	BaseCurrency    string  `json:"basecurrency" bson:"basecurrency"`
	TargetCurrency  string  `json:"targetcurrency" bson:"targetcurrency"`
	MinTriggerValue float32 `json:"mintriggervalue" bson:"mintriggervalue"`
	MaxTriggerValue float32 `json:"maxtriggervalue" bson:"maxtriggervalue"`
}
type Id struct {
	ID bson.ObjectId `bson:"_id"`
}
type Invoked struct {
	BaseCurrency    string  `json:"basecurrency"`
	TargetCurrency  string  `json:"targetcurrency"`
	CurrentRate     float32 `json:"currentrate"`
	MinTriggerValue float32 `json:"mintriggervalue"`
	MaxTriggerValue float32 `json:"maxtriggervalue"`
}

type Fixer struct {
	WebhookID bson.ObjectId      `json:"_id, omitempty" bson:"_id"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float32 `json:"rates"`
}

type Latest struct {
	BaseCurrency   string `json:"basecurrency"`
	TargetCurrency string `json:"targetcurrency"`
}

type Payload_Latest struct {
	LatestRate float32
}

type SlackPayload struct {
	Text string `json:"text"`
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
	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Insert(p)
	session.DB(db.DatabaseName).C(db.WebhookCollection).Find(bson.M{"webhookurl": p.WebhookURL}).One(&id)
	l := id.ID.Hex()

	if err != nil {
		fmt.Printf("error in Insert(), %v", err.Error())
		return l, false
	}
	return l, true

}

/*
Get the unique id of a given webhook from the storage.
*/
func (db *WebhookMongoDB) Get(keyId string) (Subscriber, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	tempP := Subscriber{}

	//check the query
	id := bson.ObjectIdHex(keyId)
	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Find(bson.M{"_id": id}).One(&tempP)

	if err != nil {
		fmt.Printf("err in Get(), %v", err.Error())
		return tempP, false
	}
	return tempP, true
}

func (db *WebhookMongoDB) Delete(keyId string) bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	id := bson.ObjectIdHex(keyId)
	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Remove(bson.M{"_id": id})

	if err != nil {
		fmt.Printf("err in Delete(), %v", err.Error())
		return false
	}
	return true
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
	defer r.Body.Close()
	var p Subscriber

	err = json.Unmarshal(body, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if len(p.WebhookURL) != 0 && len(p.BaseCurrency) != 0 && len(p.TargetCurrency) != 0 && p.MaxTriggerValue > 0 && p.MinTriggerValue > 0 {
		db.Init()
		id, ok := db.Add(p)

		if !ok {
			http.Error(w, "Not found in database", http.StatusInternalServerError)
		}
		//fmt.Printf(string(body))
		//fmt.Println(s)
		fmt.Fprint(w, id)
	} else {
		http.Error(w, "Post request body is not correctly formed", http.StatusBadRequest)
	}
}

func registeredWebhook(w http.ResponseWriter, r *http.Request) {
	db := WebhookMongoDB{
		DatabaseURL:       "mongodb://localhost",
		DatabaseName:      "Webhook",
		WebhookCollection: "WebhookPayload",
	}

	id := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case "GET":
		p, ok := db.Get(id[2])

		if !ok {
			http.Error(w, "The id is incorrect", http.StatusBadRequest)
		}
		bytes, err := json.Marshal(p)
		if err != nil {
			http.Error(w, "Error during marshaling", http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(bytes))

	case "DELETE":
		//fmt.Printf(r.Method, id[2])
		ok := db.Delete(id[2])
		if !ok {
			http.Error(w, "The id is incorrect", http.StatusBadRequest)
		}
	}
}

func retrivingLatest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()
	var l Latest

	err = json.Unmarshal(body, &l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fixerURL := "http://api.fixer.io/latest?base=" + l.BaseCurrency

	f, ok := GetFixer(fixerURL)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
	}

	var pl Payload_Latest

	for key, value := range f.Rates {
		if key == l.TargetCurrency {
			pl.LatestRate = value
		}
	}

	fmt.Fprint(w, pl.LatestRate)
}

/*
This function runs once automatically every 24hours
InvokeWebhook take out all the payloads from WebhookCollection,
get the current rate according to a certain payloads base currency and target currency
and send a notification if current rate trigger min or max value of the payload
*/
func InvokeWebhook() {
	db := WebhookMongoDB{
		DatabaseURL:       "mongodb://localhost",
		DatabaseName:      "Webhook",
		WebhookCollection: "WebhookPayload",
	}

	var form Invoked
	var results []Subscriber
	var ids []Id

	//Connection to the database
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DatabaseName).C(db.WebhookCollection).Count()
	if err != nil {
		panic(err)
	}

	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Find(nil).All(&results)
	if err != nil {
		panic(err)
	}

	//Update the payload with currect rate
	for i := 0; i < count; i++ {
		fixerURL := "http://api.fixer.io/latest?base=" + results[i].BaseCurrency

		form.BaseCurrency = results[i].BaseCurrency
		form.TargetCurrency = results[i].TargetCurrency
		form.MinTriggerValue = results[i].MinTriggerValue
		form.MaxTriggerValue = results[i].MaxTriggerValue

		f, ok := GetFixer(fixerURL)
		err = session.DB(db.DatabaseName).C(db.WebhookCollection).Find(bson.M{"webhookurl": results[i].WebhookURL}).All(&ids)
		if err != nil {
			panic(err)
		}
		SaveFixer(&f, ids[i].ID)

		if !ok {
			panic(err)
		}
		// Run through all the rates
		for key, value := range f.Rates {
			// Checks if key"currency" matches a target currency
			if key == form.TargetCurrency {
				form.CurrentRate = value

				if form.CurrentRate > form.MaxTriggerValue || form.CurrentRate < form.MinTriggerValue {

					if strings.Contains(results[i].WebhookURL, "slack") {

						var slack SlackPayload
						cr := strconv.FormatFloat(float64(form.CurrentRate), 'f', 3, 32)
						min := strconv.FormatFloat(float64(form.MinTriggerValue), 'f', 1, 32)
						max := strconv.FormatFloat(float64(form.MaxTriggerValue), 'f', 1, 32)

						slack.Text = "\nbaseCurrency: " + form.BaseCurrency + ",\ntargetCurrency: " + form.TargetCurrency + ",\ncurrentRate: " + cr + ",\nminTriggerValue: " + min + ",\nmaxTriggerValue: " + max

						postJSON, err := json.Marshal(slack)
						if err != nil {
							panic(err)
						}
						postContent := bytes.NewBuffer(postJSON)

						//Send notification to the webhookurl
						res, err := http.Post(results[i].WebhookURL, "application/json", postContent)
						if err != nil {
							panic(err)

						}
						//if recieved the 200 or 204 status code
						fmt.Printf("status: %s", res.Status)
					} else {
						//Trigger and send the notification
						postJSON, err := json.Marshal(form)
						if err != nil {
							panic(err)
						}
						postContent := bytes.NewBuffer(postJSON)

						//Send notification to the webhookurl
						res, err := http.Post(results[i].WebhookURL, "application/x-www-form-urlencoded", postContent)
						if err != nil {
							panic(err)
						}
						//if recieved the 200 or 204 status code
						fmt.Printf("status: %s", res.Status)
					}
				}
			}
		}
	}
}

func GetFixer(url string) (Fixer, bool) {
	var f Fixer
	/*	res, err := http.Get(url)
		if err != nil {
			fmt.Printf(err.Error(), http.StatusBadRequest)
			return f, false
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
	*/
	body, err := ioutil.ReadFile("base.json")
	if err != nil {
		fmt.Printf(err.Error(), http.StatusNotFound)
		return f, false
	}
	err = json.Unmarshal(body, &f)
	if err != nil {
		fmt.Printf(err.Error(), http.StatusBadRequest)
		return f, false
	}
	return f, true
}

func SaveFixer(f *Fixer, webhookID bson.ObjectId) {
	db := WebhookMongoDB{
		DatabaseURL:       "mongodb://localhost",
		DatabaseName:      "Webhook",
		WebhookCollection: "FixerPayload",
	}

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	err = session.DB(db.DatabaseName).C(db.WebhookCollection).Insert(f)
	if err != nil {
		fmt.Printf("error in SaveFixer(), %v", err.Error())
	}
}

func main() {
	os.Chdir("/home/zohaib/Desktop/Go/projects/cloud_oblig2")

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for t := range ticker.C {
			//call functions
			fmt.Print(t)
			InvokeWebhook()
		}
	}()
	http.HandleFunc("/root", postReqHandler)
	http.HandleFunc("/root/", registeredWebhook)
	http.HandleFunc("/root/latest", retrivingLatest)
	http.ListenAndServe(":8080", nil)
	ticker.Stop()
}

//invoking part
/*
//Get the webhookUrl
		url := id[2] + "//" + id[3] + "/" + id[4] + "/" + id[5] + "/" + id[6] + "/" + id[7]
		var form Invoked
		fmt.Print(url + "\n")
		//Check if it exist in db

		temp, ok := db.Get(url)

		if !ok {
			http.Error(w, "Not found in database", http.StatusInternalServerError)
		}

		//Update the Subscriber with currect rate

		fixerUrl := "http://api.fixer.io/latest?base=" + temp.BaseCurrency

		form.BaseCurrency = temp.BaseCurrency
		form.TargetCurrency = temp.TargetCurrency
		form.MinTriggerValue = temp.MinTriggerValue
		form.MaxTriggerValue = temp.MaxTriggerValue
		f, ok := GetFixer(fixerUrl)

		if !ok {
			http.Error(w, "Not found in database", http.StatusInternalServerError)
		}

		for key, value := range f.Rates {
			if key == form.TargetCurrency {
				form.CurrentRate = value
			}
		}

		//Send notification to the webhookurl
		postJSON, _ := json.Marshal(form)
		postContent := bytes.NewBuffer(postJSON)
		fmt.Println(postContent)
		res, err := http.Post(temp.WebhookURL, "application/x-www-form-urlencoded", postContent)

		if err != nil {
			http.Error(w, "Post request body is not correctly formed", http.StatusBadRequest)
		}
		//if recieved the 200 or 204 status code
		fmt.Printf("status: %s", res.Status)

*/
