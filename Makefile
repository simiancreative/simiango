test:
	go test ./...

coverage:
	go test -cover ./... -coverprofile=coverage.out
