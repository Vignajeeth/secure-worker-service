package helper

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Job struct {
	ID        int       `json:"id"`
	Status    uint      `json:"status"`
	Command   string    `json:"command"`
	UserId    string    `json:"user"` // Could be changed to the User struct.
	StartTime time.Time `json:"startTime"`
	StopTime  time.Time `json:"stopTime"`
	Result    string    `json:"result"`
}

type User struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Access    uint   `json:"access"`
}

type ErrorLog struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var Users []User // #GLOBAL variable
var JobDB []Job
var jobCounter int
var Ch [100]chan bool

var users = map[string]string{
	// UserId: Password
	"admin": "admin",
	"read":  "read",
	"write": "write",
	"no":    "no",
}

func authWrapper(w http.ResponseWriter, r *http.Request, access map[int]bool) (string, bool) {

	username, password, ok := r.BasicAuth()

	if !ok {
		log.Println("Error parsing basic auth")
		w.WriteHeader(http.StatusUnauthorized)
		return "", false
	}

	var userAccess int
	for i := range Users {
		if Users[i].UserId == username {
			userAccess = int(Users[i].Access)
			break
		}
	}

	_, present := access[userAccess]
	if val, ok := users[username]; ok {
		if val == password && present {
			log.Printf("Authentication Successful")
			w.WriteHeader(http.StatusOK)
			return username, true

		} else {
			log.Println("Credentials incorrect or access denied.")
			w.WriteHeader(http.StatusUnauthorized)
			return "", false
		}
	}
	return "", false
}

func GetTLSConfig(host, caCertFile string, clientAuth tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool           // Empty certificate pool.
	if clientAuth > tls.RequestClientCert { // RequestClientCert=1
		caCert, err = ioutil.ReadFile(caCertFile) // Reading public key of CA.
		if err != nil {
			log.Fatal("Error opening cert file", caCertFile, ", error ", err)
		}
		caCertPool = x509.NewCertPool()       // New empty pool.
		caCertPool.AppendCertsFromPEM(caCert) // Appends CA to the pool so that it is recognised.
	}
	return &tls.Config{ServerName: host, ClientAuth: clientAuth, ClientCAs: caCertPool, MinVersion: tls.VersionTLS12}
}

func StartJob(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Job Hit.")
	user, allowed := authWrapper(w, r, map[int]bool{1: true, 3: true})
	if !allowed {
		log.Println("Authentication Fail in StartJob")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var job Job
	_ = json.NewDecoder(r.Body).Decode(&job)

	j := Workerlib(job, jobCounter, user)
	jobCounter++
	// var job Job
	// _ = json.NewDecoder(r.Body).Decode(&job)
	// job.ID = strconv.Itoa(rand.Intn(5000))
	// Jobs = append(Jobs, job) // Appending to database.
	json.NewEncoder(w).Encode(j)
}

func GetJob(w http.ResponseWriter, r *http.Request) { // Resouce not found. Do
	_, allowed := authWrapper(w, r, map[int]bool{2: true, 3: true})
	if !allowed {
		log.Println("Authentication Fail in GetJob")
		return
	}

	// log.Println(JobDB)
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	jobid, _ := strconv.Atoi(params["id"])
	log.Println(jobid)

	var j Job
	for i := range JobDB {
		if JobDB[i].ID == jobid {
			j = JobDB[i]
			break
		}
	}
	json.NewEncoder(w).Encode(j)

}

func StopJob(w http.ResponseWriter, r *http.Request) {
	_, allowed := authWrapper(w, r, map[int]bool{1: true, 3: true})
	if !allowed {
		log.Println("Authentication Fail in StopJob")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	jobid, _ := strconv.Atoi(params["id"])
	log.Println(jobid)
	KillJob(jobid)

	var j Job
	for i := range JobDB {
		if JobDB[i].ID == jobid {
			j = JobDB[i]
			break
		}
	}
	json.NewEncoder(w).Encode(j)

}
