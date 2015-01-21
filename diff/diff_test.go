package diff

import (
	"github.com/k0kubun/pp"
	"io/ioutil"
	"os"
	"testing"

	"github.com/walf443/mgr/sqlparser/mysql"
)

func init() {
	if os.Getenv("DEBUG") == "" {
		pp.SetDefaultOutput(ioutil.Discard)
	}
}


func TestDiffDatabase(t *testing.T) {
	before := "CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT); CREATE TABLE foo (id int unsigned not null AUTO_INCREMENT);"
	after := "CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT); CREATE TABLE bar (id int unsigned not null AUTO_INCREMENT);"
	beforeStmt := parseSQL(t, before)
	afterStmt := parseSQL(t, after)
	result := Extract(beforeStmt, afterStmt)
	pp.Print(result)
	if !checkTable(result.Added[0], "`bar`") {
		t.Errorf("bar should be added")
	}
	if !checkTable(result.Removed[0], "`foo`") {
		t.Errorf("foo should be added")
	}
}

func TestDiffTable(t *testing.T) {
	testDiffTable(
		t,
		"general case",
		"CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT, foo int(10) unsigned not null, key foo (foo))",
		"CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT, bar int(10) unsigned not null, key bar (bar))",
		"ALTER TABLE `hoge` DROP `foo`, DROP INDEX `foo`, ADD `bar` INT(10) UNSIGNED NOT NULL , ADD INDEX `bar` (`bar`);",
	)
	testDiffTable(
		t,
		"unique key",
		"CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT, foo int(10) unsigned not null, unique key foo (foo))",
		"CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT, bar int(10) unsigned not null, unique key bar (bar))",
		"ALTER TABLE `hoge` DROP `foo`, DROP INDEX `foo`, ADD `bar` INT(10) UNSIGNED NOT NULL , ADD UNIQUE INDEX `bar` (`bar`);",
	)
}

func testDiffTable(t *testing.T, name, before string, after string, expected string) {
	beforeStmt := parseCreateTableStatement(t, before)
	afterStmt := parseCreateTableStatement(t, after)
	result := ExtractTableSchemaDifference(beforeStmt, afterStmt)
	sql := result.ToQuery()
	if sql != expected {
		t.Errorf("failed to testDiffTable \"%s\":\nBefore schema:\n%s\nAfter schema:\n%s\nExpected diff: \t%s\nBut got: \t%s", name, before, after, expected, sql);
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

func parseCreateTableStatement(t *testing.T, sql string) *mysql.CreateTableStatement {
	stmt := parseSQL(t, sql)
	v, ok := stmt[0].(*mysql.CreateTableStatement)
	if !ok {
		t.Errorf("Faied to extract CreateTableStatement")
	}
	return v
}

func checkTable(target mysql.Statement, tableName string) bool {
	if v, ok := target.(*mysql.CreateTableStatement); ok {
		return v.TableName.ToQuery() == tableName
	}
	return false
}

func checkColumn(target mysql.CreateDefinition, columnName string) bool {
	if v, ok := target.(*mysql.CreateDefinitionColumn); ok {
		return v.ColumnName.ToQuery() == columnName
	}
	return false
}

func checkIndex(target mysql.CreateDefinition, indexName string) bool {
	if v, ok := target.(*mysql.CreateDefinitionIndex); ok {
		return v.Name.ToQuery() == indexName
	} else if v, ok := target.(*mysql.CreateDefinitionUniqueIndex); ok {
		return v.Name.ToQuery() == indexName
	}
	return false
}
