test: mgr

mgr: cli/*.go cli_test diff_test sqlparser_test
	cd cli && go build -o ../mgr

sqlparser_test:
	cd sqlparser && make test

diff_test: sqlparser_test
	cd diff && make test

cli_test: diff_test
	cd cli && go test -v ./...

get-deps:
	go get github.com/k0kubun/pp
	go get gopkg.in/yaml.v2

.PHONY: test sqlparser_test diff_test get-deps
