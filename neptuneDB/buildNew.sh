go mod tidy
./uprev/uprev
swag init
env GOOS=linux GOARCH=amd64 go build -o goapiTemplate
