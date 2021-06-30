package helper

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func checkPassword(inputPassword, storedPassword string) (string, bool) {
	if storedPassword == inputPassword {
		message := "Correct password."
		log.Println(message)
		return message, true
	} else {
		message := "Password is not correct."
		log.Println(message)
		return message, false
	}

}

func checkAccess(userAccess int, access map[int]bool) (string, bool) {
	_, accessAllowed := access[userAccess]
	if accessAllowed {
		message := "Access granted."
		log.Println(message)
		return message, true
	} else {
		message := "Access denied."
		log.Println(message)
		return message, false
	}

}

// authWrapper checks if the user is authorized to perform a certain job. If the bool is false, then the
// string returns the reason for error. If true, it returns the username.
func authWrapper(w http.ResponseWriter, r *http.Request, access map[int]bool) (string, string, bool) {
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Println("Error parsing username and password")
		w.WriteHeader(http.StatusUnauthorized)
		return "", "Error parsing username and password", false
	}

	userAccess := -1
	for i := range Users {
		if Users[i].UserId == username {
			userAccess = int(Users[i].Access)
			break
		}
	}
	if userAccess == -1 {
		return "", "User is not present in the DB.", false
	}

	if savedPassword, ok := users[username]; ok {

		passwordMessage, isPasswordCorrect := checkPassword(password, savedPassword)
		accessMessage, isAccessAllowed := checkAccess(userAccess, access)

		// Authentication and authorisation successful.
		if isPasswordCorrect && isAccessAllowed {
			log.Printf("Authentication Successful")
			w.WriteHeader(http.StatusOK)
			return username, "", true
		} else {

			return "", passwordMessage + " \n " + accessMessage, false

		}
	}
	return "", "User is not present in the credentials data.", false
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
	user, erMessage, allowed := authWrapper(w, r, map[int]bool{1: true, 3: true})
	if !allowed {
		log.Println("Authentication fail in StartJob")
		e := ErrorLog{
			Message: erMessage,
		}
		json.NewEncoder(w).Encode(e)
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
	_, erMessage, allowed := authWrapper(w, r, map[int]bool{2: true, 3: true})
	if !allowed {
		log.Println("Authentication Fail in GetJob")
		e := ErrorLog{
			Message: erMessage,
		}
		json.NewEncoder(w).Encode(e)
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
	// User is not required because stopping jobs solely depends on access level of the user. Therefore
	// any user with full access can stop any job by any other person.
	_, erMessage, allowed := authWrapper(w, r, map[int]bool{1: true, 3: true})
	if !allowed {
		log.Println("Authentication Fail in StopJob")
		e := ErrorLog{
			Message: erMessage,
		}
		json.NewEncoder(w).Encode(e)
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
