test:
	go test -race -cover $$(go list ./...)
