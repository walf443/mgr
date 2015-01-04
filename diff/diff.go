/* 	import "github.com/walf443/mig/diff"
	func main() {
		result := diff.Extract(before_schema, after_schema)
	}
*/
package diff

import (
	"github.com/walf443/sqlparser/mysql"
)

type SchemaDifference struct {
	Added []mysql.Statement
	Removed []mysql.Statement
	Modified []mysql.Statement
}

func Extract(before []mysql.Statement, after []mysql.Statement) *SchemaDifference {
	var result SchemaDifference

	tableNameOf := make(map[string]mysql.Statement)
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
