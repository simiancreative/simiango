test:
	./bin/test.sh

coverage:
	go test -cover ./... -coverprofile=coverage.out

coveragehtml:
	go tool cover -html=coverage.out
