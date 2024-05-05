module github.com/GetterSethya/userService

go 1.21.1

require (
	github.com/GetterSethya/library v1.0.0
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.22
)

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
)

replace github.com/GetterSethya/library => ../library
