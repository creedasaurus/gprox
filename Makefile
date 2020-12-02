LOCAL_BIN=$(shell pwd)/bin

installs:
	GOBIN=$(LOCAL_BIN) go install github.com/markbates/pkger/cmd/pkger

packfiles: installs
	$(LOCAL_BIN)/pkger -include /localhost.cert -include /localhost.key

