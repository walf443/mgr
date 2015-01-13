test: mgr

mgr: cli/*.go diff_test sqlparser_test
	cd cli && go build -o ../mgr

sqlparser_test:
	cd sqlparser && make test

diff_test: sqlparser_test
	cd diff && make test

.PHONY: test
