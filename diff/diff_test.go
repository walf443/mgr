package diff

import (
	"github.com/k0kubun/pp"
	"io/ioutil"
	"os"
	"testing"

	"github.com/walf443/mig/sqlparser/mysql"
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
	before := "CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT, foo int(10) unsigned not null, key foo (foo))"
	after := "CREATE TABLE hoge (id int unsigned not null AUTO_INCREMENT, bar int(10) unsigned not null, key bar (bar))"
	beforeStmt := parseCreateTableStatement(t, before)
	afterStmt := parseCreateTableStatement(t, after)
	result := ExtractTableSchemaDifference(beforeStmt, afterStmt)
	if len(result.Added) != 2 {
		t.Errorf("Expect len(Added) to 2, But Got %d", len(result.Added))
		return
	}
	if len(result.Removed) != 2 {
		t.Errorf("Expect len(Removed) to 2, But Got %d", len(result.Removed))
		return
	}
	if !checkColumn(result.Added[0], "`bar`") {
		t.Errorf("column bar should be added")
	}
	if !checkColumn(result.Removed[0], "`foo`") {
		t.Errorf("column foo should be added")
	}
	if !checkIndex(result.Added[1], "`bar`") {
		t.Errorf("index bar should be added")
	}
	if !checkIndex(result.Removed[1], "`foo`") {
		t.Errorf("index bar should be added")
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
