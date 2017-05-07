package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/vault/api"
)

func main() {
	// tokenValue := os.Getenv("VAULT_TOKEN")
	config := api.DefaultConfig()
	// config.Address = "http://127.0.0.1:8200"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	sealStatus, err := client.Sys().SealStatus()
	if err != nil {
		fmt.Printf("Error checking seal status: %s", err)
	}

	fmt.Printf(
		"Sealed: %v\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Unseal Progress: %d\n"+
			"Unseal Nonce: %v\n"+
			"Version: %s",
		sealStatus.Sealed,
		sealStatus.N,
		sealStatus.T,
		sealStatus.Progress,
		sealStatus.Nonce,
		sealStatus.Version)
}

func httpRequest() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8200/v1/sys/seal-status", nil)
	req.Header.Add("X-Vault-Token", "f593a179-36a0-3c66-8b7d-b0dd718659f1")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	// Display Results
	fmt.Println(string(respBody))
}
