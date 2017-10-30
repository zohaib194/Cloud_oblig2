package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// Fixer (To save data of fixer.io)
type Fixer struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

/*
	Get the json from Fixer.io
*/
func (f *Fixer) GetFixer(url string) (*Fixer, bool) {

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf(err.Error(), http.StatusBadRequest)
		return f, false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

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

/*
	Save Fixer payload in the collection
*/
func (f *Fixer) SaveFixer() bool {
	db := WebhookMongoDB{
		DatabaseURL:  "mongodb://localhost",
		DatabaseName: "Webhook",
		Collection:   "FixerPayload",
	}

	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	f.Date = time.Now().Format("2006-01-02")
	err = session.DB(db.DatabaseName).C(db.Collection).Insert(&f)
	if err != nil {
		fmt.Printf("error in SaveFixer(), %v", err.Error())
		return false
	}

	return true
}

func (f *Fixer) LatestFixer() {
	//Send request to Fixer.io
	fixerURL := "http://api.fixer.io/latest?base=EUR"
	result, ok := f.GetFixer(fixerURL)
	if !ok {
		fmt.Print("Error occured during Get req to Fixer.io")
	}
	f.SaveFixer()
}

func main() {
	ticker := time.NewTicker(time.Second * 120)
	go func() {
		for t := range ticker.C {
			//call functions
			fmt.Printf("\n", t)
			f.LatestFixer()
			//GetFixerSevenDays(time.Now().AddDate(0, 0, -7), time.Now())
		}
	}()
	ticker.Stop()
}
