test:
	go test -race -cover $$(go list ./...)

benchmark:
	go test -bench=.  -benchmem=true -cpu=1,2,4,8 -run Bench
