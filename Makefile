build:
	go build -o bin/loadbalancer cmd/balancer/main.go

run: build
	./bin/loadbalancer

docker:
	docker build -t bin/loadbalancer .
	docker-compose up

clean:
	rm bin/loadbalancer
