package utils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dustin/go-jsonpointer"
	"github.com/jordan-wright/email"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/parnurzeal/gorequest"
)

// Getopt reads environment variables.
// If not found will return a supplied default value
func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

// Assert asserts there was no error, else log.Fatal
func Assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CopyFile copies file from dst to scr
func CopyFile(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

// GetSHA256 calculates a file's sha256sum
func GetSHA256(name string) string {

	dat, err := ioutil.ReadFile(name)
	Assert(err)

	h256 := sha256.New()
	_, err = h256.Write(dat)
	Assert(err)

	return fmt.Sprintf("%x", h256.Sum(nil))
}

// RunCommand runs cmd on file
func RunCommand(cmd string, args ...string) string {

	cmdOut, err := exec.Command(cmd, args...).Output()
	if len(cmdOut) == 0 {
		Assert(err)
	}

	return string(cmdOut)
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

// RemoveDuplicates removes duplicate items from a list
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

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
