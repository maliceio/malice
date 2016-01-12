package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/dustin/go-jsonpointer"
	"github.com/jordan-wright/email"
	er "github.com/maliceio/malice/libmalice/errors"
)

// ParseJSON returns a JSON value for a given key
// NOTE: https://godoc.org/github.com/dustin/go-jsonpointer
func ParseJSON(data []byte, path string) (out string) {
	var o map[string]interface{}
	er.CheckError(json.Unmarshal(data, &o))
	out = jsonpointer.Get(o, string(data)).(string)
	return
}

// ParseMail takes in an HTTP Request and returns an Email object
// TODO: This function will likely be changed to take in a []byte
func ParseMail(r *http.Request) (email.Email, error) {
	e := email.Email{}
	m, err := mail.ReadMessage(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(m.Body)
	e.HTML = body
	return e, err
}
