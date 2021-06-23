package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/Vignajeeth/secure-worker-service/server/helper"
	"github.com/gorilla/mux"
)

// Error response and testing
// Check User in jobdb
func main() {

	for i := range helper.Ch {
		helper.Ch[i] = make(chan bool)
	}

	// Sample data
	helper.Users = append(helper.Users, helper.User{UserId: "admin", FirstName: "admin", LastName: "admin", Access: 3})
	helper.Users = append(helper.Users, helper.User{UserId: "read", FirstName: "read", LastName: "read", Access: 2})
	helper.Users = append(helper.Users, helper.User{UserId: "write", FirstName: "write", LastName: "write", Access: 1})
	helper.Users = append(helper.Users, helper.User{UserId: "no", FirstName: "no", LastName: "no", Access: 0})

	// Initialising a router.
	r := mux.NewRouter()
	port := "8080"
	host := "localhost"
	clientAuth := 0
	caCert := "certs/out/Cert_Auth.crt"     // Public key of certificate authority.
	serverCert := "certs/out/localhost.crt" // Public key of server authorised by CA.
	srcKey := "certs/out/localhost.key"     // Private key of server.

	server := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    helper.GetTLSConfig(host, caCert, tls.ClientAuthType(clientAuth)),
	}

	// Route Handlers

	r.HandleFunc("/api/jobs", helper.StartJob).Methods("POST")          // curl -X POST --cacert certs/out/Cert_Auth.crt https://localhost:8080/api/jobs -d @json/request_body4.json -v
	r.HandleFunc("/api/jobs/{id}", helper.GetJob).Methods("GET")        // curl -X GET --cacert certs/out/Cert_Auth.crt https://localhost:8080/api/jobs/0  -v
	r.HandleFunc("/api/jobs/{id}/stop", helper.StopJob).Methods("POST") // curl -X POST --cacert certs/out/Cert_Auth.crt https://localhost:8080/api/jobs/0/stop  -v

	// Server listens at 8000 and serves requests.
	log.Println("Starting a server...")
	log.Fatal(server.ListenAndServeTLS(serverCert, srcKey))

}
