module api-gateway

go 1.24.2

require (
	github.com/BSanjik/LUNARY/services/common v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/time v0.12.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/BSanjik/LUNARY/services/common => ../common
