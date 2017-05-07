package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	httpRequest()
	// // tokenValue := os.Getenv("VAULT_TOKEN")
	// config := api.DefaultConfig()
	// // config.Address = "http://127.0.0.1:8200"
	// client, err := api.NewClient(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // client.SetToken(tokenValue)
	// // if v := client.Token(); v != tokenValue {
	// // 	log.Fatalf("bad: %s", v)
	// // }

	// resp, err := client.RawRequest(client.NewRequest("GET", "/v1/sys/init"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Copy the response
	// var buf bytes.Buffer
	// io.Copy(&buf, resp.Body)
	// fmt.Println(buf.String())
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
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))
	fmt.Println(string(respBody))
}
