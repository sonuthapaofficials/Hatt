package main

import (
	"encoding/json"
	"fmt"
	"hatt/configuration"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Credentials = map[string]map[string]map[string]string

func login(website string) {
	var conf configuration.Config = deserializeWebsiteConf(website + ".json")

	savedCredentials := deserializeCredentials()
	websiteCredentials := savedCredentials[website]["credentials"] // savedCredentials[website]["credentials"]

	hc := http.Client{}

	form := url.Values{}
	for field, value := range websiteCredentials {
		form.Add(field, value)
	}

	req, _ := http.NewRequest("POST", conf.Login.Url, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	for _, cookie := range resp.Cookies() {
		for _, token := range conf.Login.Tokens {
			if token == cookie.Name {
				fmt.Println(cookie.Name, cookie.Value)
				savedCredentials[website]["tokens"][cookie.Name] = cookie.Value
			}
		}
	}

	updatedCredentialsJson, _ := json.Marshal(savedCredentials)
	_ = ioutil.WriteFile(CREDENTIALS_PATH, updatedCredentialsJson, 0644)

}
