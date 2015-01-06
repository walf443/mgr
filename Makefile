test: diff_test

diff_test:
	cd diff && make test

get-deps:
	go get github.com/walf443/sqlparser/mysql

.PHONY: test get-deps
