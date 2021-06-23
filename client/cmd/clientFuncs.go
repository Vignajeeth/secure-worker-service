package cmd

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func withTLS() http.Transport {
	caCertFile := "out/Cert_Auth.crt"  // Public key
	clientCertFile := "out/client.crt" // Public key
	clientKeyFile := "out/client.key"  // Private key

	var cert tls.Certificate
	var err error

	if clientCertFile != "" && clientKeyFile != "" {
		cert, err = tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
		}
	}

	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}
	return *t
}

func startRequest(command, username, password string) {

	t := withTLS()
	client := &http.Client{Transport: &t}

	postBody, _ := json.Marshal(map[string]string{
		"command": command,
	})
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", "https://localhost:8080/api/jobs", responseBody)
	if err != nil {
		log.Fatalln(err)
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}
func stopRequest(jobId int, username, password string) {

	t := withTLS()
	client := &http.Client{Transport: &t}

	req, err := http.NewRequest("POST", "https://localhost:8080/api/jobs/"+strconv.Itoa(jobId)+"/stop", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}
func queryRequest(jobId int, username, password string) {

	t := withTLS()

	// Initialising a http client for more control. Can be used without initialising but will be default.
	// client := &http.Client{}
	client := &http.Client{Transport: &t}

	req, err := http.NewRequest("GET", "https://localhost:8080/api/jobs/"+strconv.Itoa(jobId), nil)

	// If get request fails, control goes to the below block
	if err != nil {
		log.Fatalln(err)
	}
	// Setting up basic auth
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	// If get request is successful, we should close the body after the function exits to prevent resource leaks.
	defer resp.Body.Close()

	// Reading the message from the response.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}

func Client() {
	withTLS()

}
