package helper

import "time"

type Job struct {
	ID        int        `json:"id"`
	Status    StatusCode `json:"status"`
	Command   string     `json:"command"`
	UserId    string     `json:"user"` // Could be changed to the User struct.
	StartTime time.Time  `json:"startTime"`
	StopTime  time.Time  `json:"stopTime"`
	Result    string     `json:"result"`
}

type User struct {
	UserId    string      `json:"userId"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Access    AccessLevel `json:"access"`
}
type AccessLevel uint

const (
	NoAccess AccessLevel = iota
	WriteOnly
	ReadOnly
	FullAccess
)

func (d AccessLevel) String() string {
	return [...]string{"NoAccess", "WriteOnly", "ReadOnly", "FullAccess"}[d]
}

type StatusCode uint

const (
	Created StatusCode = iota
	Running
	Completed
	Failed
)

func (d StatusCode) String() string {
	return [...]string{"Created", "Running", "Completed", "Failed"}[d]
}

type ErrorLog struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// These variables could be made into a DB later.
var Users []User
var JobDB []Job
var jobCounter int
var Ch [100]chan bool

// These should be salted and hashed later.
var users = map[string]string{
	"admin": "admin",
	"read":  "read",
	"write": "write",
	"no":    "no",
}
