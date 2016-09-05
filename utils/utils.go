package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"github.com/dustin/go-jsonpointer"
	"github.com/jordan-wright/email"
	er "github.com/maliceio/malice/malice/errors"
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

// AskForConfirmation prompts user for yes/no response
func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "yes"}
	nokayResponses := []string{"n", "no"}
	if StringInSlice(strings.ToLower(response), okayResponses) {
		return true
	}
	if StringInSlice(strings.ToLower(response), nokayResponses) {
		return false
	}
	fmt.Println("Please type yes or no and then press enter:")
	return AskForConfirmation()
}

// StringInSlice returns whether or not a string exists in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetOpt returns Env var or default
func GetOpt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

// GetHashType returns the hash type (md5, sha1, sha256, sha512)
func GetHashType(hash string) (string, error) {
	var validMD5 = regexp.MustCompile(`^[a-fA-F\d]{32}$`)
	var validSHA1 = regexp.MustCompile(`^[a-fA-F\d]{40}$`)
	var validSHA256 = regexp.MustCompile(`^[a-fA-F\d]{64}$`)
	var validSHA512 = regexp.MustCompile(`^[a-fA-F\d]{128}$`)

	switch {
	case validMD5.MatchString(hash):
		return "md5", nil
	case validSHA1.MatchString(hash):
		return "sha1", nil
	case validSHA256.MatchString(hash):
		return "sha256", nil
	case validSHA512.MatchString(hash):
		return "sha512", nil
	default:
		return "", errors.New("This is not a valid hash.")
	}
}
