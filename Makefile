test:
	TOKEN_SECRET=wombat go test ./...

coverage:
	go test -cover ./... -coverprofile=coverage.out
