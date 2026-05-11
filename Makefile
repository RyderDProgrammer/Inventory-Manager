run:
	go run ./cmd/server/main.go

test:
	go test ./...

build:
	go build -o bin/server ./cmd/server/main.go

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

k8s-deploy:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/
