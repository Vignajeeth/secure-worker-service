# Secure Worker Service

This is a worker library which can orchestrate Linux processes using a CLI to start, stop and get the status of a running process. It has an HTTPS API using a strong set of cipher suites for TLS; Basic Authentication and a simple authorization scheme. The CLI connects to the API, which in turn acts as a wrapper to the library. User management with various levels of access is also implemented.

## Usage

Run the server in the terminal by:

```bash
cd server
go run server.go
```

Use the client to send requests by:

```bash
cd client
go build
./client start echo foo
```
