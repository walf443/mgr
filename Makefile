test: diff_test

sqlparser_test:
	cd sqlparser && make test

diff_test: sqlparser_test
	cd diff && make test

.PHONY: test
