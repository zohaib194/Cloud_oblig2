package main

/*
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Fixer struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

/*
	Get the json from Fixer.io

func (f *Fixer) GetFixer(url string) (*Fixer, bool) {
	//var f *Fixer
	/*
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf(err.Error(), http.StatusBadRequest)
			return f, false
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

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
*/
