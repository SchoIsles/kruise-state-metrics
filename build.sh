GOOS=linux GOARCH=amd64 go build -o bin/kruise-state-metrics main.go
docker build -t kruise-state-metrics .
