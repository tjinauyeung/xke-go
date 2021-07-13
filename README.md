# XKE Go demo

## Run
```go
go run .
```

## Build 
```go
go build
```

## Watch
```go
// install fresh pkg
go get github.com/pilu/fresh

// start watcher
fresh
```

## DB
Update the db config to connect to your own postgres db

```go
// fill in your db config here
const (
	host     = "localhost"
	user     = "postgres"
	password = "password"
	dbname   = "golangdb"
	port     = 5432
)
```
