module github.com/ksrinimba/demo-jwt-service

go 1.21.3
replace       github.com/ksrinimba/ssd-jwt-auth => ./../ssd-jwt-auth

require (
	github.com/gorilla/mux v1.8.1
	github.com/ksrinimba/ssd-jwt-auth v0.0.3
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
)
