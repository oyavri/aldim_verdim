run_fe:
	go run cmd/frontend/main.go

run_worker:
	go run cmd/worker/main.go

run:
	make run_fe
	make run_worker

tidy:
	go mod tidy -v
	go fmt ./...

test-coverage:
	rm -rf coverage
	mkdir coverage
	go test -v -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
