LOCAL_BIN=$(shell pwd)/bin

installs:
	GOBIN=$(LOCAL_BIN) go install github.com/markbates/pkger/cmd/pkger
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik

packfiles: installs
	$(LOCAL_BIN)/statik -m -src=cert
