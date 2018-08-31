package utils

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// AppHelpTemplate is a default malice plugin help template
var AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

// CamelCase converts strings to their camel case equivalent
func CamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}

// Getopt reads environment variables.
// If not found will return a supplied default value
func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

// Getopts reads from user input then environment variable and finally a sane default.
func Getopts(userInput, envVar, dfault string) string {

	if len(strings.TrimSpace(userInput)) > 0 {
		return userInput
	}
	value := os.Getenv(envVar)
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
func RunCommand(ctx context.Context, cmd string, args ...string) (string, error) {

	var c *exec.Cmd

	if ctx != nil {
		c = exec.CommandContext(ctx, cmd, args...)
	} else {
		c = exec.Command(cmd, args...)
	}

	output, err := c.Output()
	if err != nil {
		return string(output), err
	}

	// check for exec context timeout
	if ctx != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("command %s timed out", cmd)
		}
	}

	return string(output), nil
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
		return "", errors.New("this is not a valid hash")
	}
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

// Unzip unzips archive to target location
func Unzip(archive, target string) error {

	// fmt.Println("Unzip archive ", target)

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(target, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

// SliceContainsString returns if slice contains substring
func SliceContainsString(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(b, a) {
			return true
		}
	}
	return false
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
