/* 	import "github.com/walf443/mgr/diff"
func main() {
	result := diff.Extract(before_schema, after_schema)
}
*/
package diff

import (
	// "github.com/k0kubun/pp"
	"github.com/walf443/mgr/sqlparser/mysql"
)

type DatabaseSchemaDifference struct {
	Added    []mysql.Statement
	Removed  []mysql.Statement
	Modified []mysql.Statement
}

type TableSchemaDifference struct {
	Before   *mysql.CreateTableStatement
	After    *mysql.CreateTableStatement
	Added    []mysql.CreateDefinition
	Removed  []mysql.CreateDefinition
	Modified []mysql.CreateDefinition
}

func Extract(before []mysql.Statement, after []mysql.Statement) *DatabaseSchemaDifference {
	var result DatabaseSchemaDifference

	tableNameOf := make(map[string]*mysql.CreateTableStatement)
	for _, stmt := range before {
		if v, ok := stmt.(*mysql.CreateTableStatement); ok {
			key := v.TableName.ToQuery()
			tableNameOf[key] = v
		}
	}

	for _, stmt := range after {
		if v, ok := stmt.(*mysql.CreateTableStatement); ok {
			key := v.TableName.ToQuery()
			if _, ok := tableNameOf[key]; ok {
				// TODO: detect Modified
				if v.ToQuery() != tableNameOf[key].ToQuery() {
					result.Modified = append(result.Modified, v) // 変更前のデータが構造的に必要
				}
				delete(tableNameOf, key)
			} else {
				result.Added = append(result.Added, v)
			}
		}
	}

	for _, statement := range tableNameOf {
		result.Removed = append(result.Removed, statement)
	}

	return &result
}

// TODO: How to check primary key difference?
func ExtractTableSchemaDifference(x *mysql.CreateTableStatement, y *mysql.CreateTableStatement) *TableSchemaDifference {
	var result TableSchemaDifference
	result.Before = x
	result.After = y

	columnNameOf := make(map[string]mysql.CreateDefinition)
	indexNameOf := make(map[string]mysql.CreateDefinition)
	for _, definition := range x.CreateDefinitions {
		switch v := definition.(type) {
		case *mysql.CreateDefinitionColumn:
			key := v.ColumnName.ToQuery()
			columnNameOf[key] = definition
		case *mysql.CreateDefinitionIndex:
			key := v.Name.ToQuery()
			indexNameOf[key] = definition
		case *mysql.CreateDefinitionUniqueIndex:
			key := v.Name.ToQuery()
			indexNameOf[key] = definition
		}
	}

	for _, definition := range y.CreateDefinitions {
		switch v := definition.(type) {
		case *mysql.CreateDefinitionColumn:
			key := v.ColumnName.ToQuery()
			if _, ok := columnNameOf[key]; ok {
				delete(columnNameOf, key)
				// TODO: check modified
			} else {
				result.Added = append(result.Added, definition)
			}
		case *mysql.CreateDefinitionIndex:
			key := v.Name.ToQuery()
			if _, ok := indexNameOf[key]; ok {
				delete(indexNameOf, key)
				// TODO: check modified
			} else {
				result.Added = append(result.Added, definition)
			}
		case *mysql.CreateDefinitionUniqueIndex:
			key := v.Name.ToQuery()
			if _, ok := indexNameOf[key]; ok {
				delete(indexNameOf, key)
				// TODO: check modified
			} else {
				result.Added = append(result.Added, definition)
			}
		}
	}

	for _, definition := range columnNameOf {
		result.Removed = append(result.Removed, definition)
	}
	for _, definition := range indexNameOf {
		result.Removed = append(result.Removed, definition)
	}

	return &result
}
