
BINARY=engine

.PHONY: clean run test init domain fmt

fmt:
	gofmt -s -w .

# run golang service 
start:
	go run app/main.go

# run dev docker-compose
run:
	docker-compose up -d --build

# stop docker compose
stop:
	docker-compose down

# build docker image
docker:
	docker build -t $(basename $PWD):latest .

# run golang test
test:
	go test ./...

# clean binary build of this go service
clean: 
	if [ -f ${BINARY} ]; then rm ${BINARY} ; fi

# run goose local up
goose-local-up:
	goose mysql "root:root@/kafka-cluster-go?parseTime=true" up

# run goose local down
goose-local-down:
	goose mysql "root:root@/kafka-cluster-go?parseTime=true" down


vendor:
	go mod vendor

engine:
	go build -o engine app/main.go

.PHONY: start run stop docker test vendor engine goose-local-up goose-local-down
