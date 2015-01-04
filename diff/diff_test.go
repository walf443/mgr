package diff

import (
	// "github.com/k0kubun/pp"
	"testing"

	"github.com/walf443/sqlparser/mysql"
)

func TestDiff(t *testing.T) {
	before := "CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT); CREATE TABLE foo (id int unsigned not null AUTO_INCREMENT);"
	after := "CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT); CREATE TABLE bar (id int unsigned not null AUTO_INCREMENT);"
	beforeStmt := parseSQL(t, before)
	afterStmt := parseSQL(t, after)
	result := Extract(beforeStmt, afterStmt)
	// pp.Print(result)
	if !checkTable(result.Added[0], "`bar`") {
		t.Errorf("bar should be added")
	}
	if !checkTable(result.Removed[0], "`foo`") {
		t.Errorf("foo should be added")
	}
}

func parseSQL(t *testing.T, sql string) []mysql.Statement {
	s := new(mysql.Scanner)
	s.Init(sql)
	stmt, err := mysql.Parse(s)
	if err != nil {
		t.Errorf("Faied to parse SQL: %s, error: %q", sql, err)
	}
	return stmt
}

func checkTable(target mysql.Statement, tableName string) bool {
	if v, ok := target.(*mysql.CreateTableStatement); ok {
		return v.TableName.ToQuery() == tableName
	}
	return false
}
