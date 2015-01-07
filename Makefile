test: diff_test

sqlparser_test:
	cd sqlparser && make test

diff_test: sqlparser_test
	cd diff && make test

get-deps:
	go get github.com/walf443/sqlparser/mysql

.PHONY: test get-deps
