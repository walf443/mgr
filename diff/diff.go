package diff

import (
	"github.com/walf443/mgr/sqlparser/mysql"
)

type DatabaseSchemaDifference struct {
	Added    []mysql.Statement
	Removed  []mysql.Statement
	Modified []TableSchemaDifference
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
				if v.ToQuery() != tableNameOf[key].ToQuery() {
					result.Modified = append(result.Modified, ExtractTableSchemaDifference(tableNameOf[key], v))
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
func ExtractTableSchemaDifference(x *mysql.CreateTableStatement, y *mysql.CreateTableStatement) TableSchemaDifference {
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

	return result
}

func (x *DatabaseSchemaDifference) Changes() []string {
	var sqls []string
	for _, stmt := range x.Added {
		if v, ok := stmt.(*mysql.CreateTableStatement); ok {
			newStmt := convertCreateTableStatement(v)
			sqls = append(sqls, newStmt.ToQuery())
		}
	}
	for _, stmt := range x.Removed {
		if v, ok := stmt.(*mysql.CreateTableStatement); ok {
			newStmt := convertDropTableStatement(v)
			sqls = append(sqls, newStmt.ToQuery())
		}
	}

	for _, stmt := range x.Modified {
		sqls = append(sqls, stmt.ToQuery())
	}

	return sqls
}

func (x *TableSchemaDifference) ToQuery() string {

	if ( len(x.Added) == 0 && len(x.Removed) == 0 && len(x.Modified) == 0 ) {
		return "";
	}

	var specs []mysql.AlterSpecification
	for _, def := range x.Removed {
		switch def := def.(type) {
		case *mysql.CreateDefinitionColumn:
			spec := mysql.AlterSpecificationDropColumn{}
			spec.ColumnName = def.ColumnName
			specs = append(specs, &spec)
		case *mysql.CreateDefinitionPrimaryIndex:
			// TODO
		case *mysql.CreateDefinitionUniqueIndex:
			spec := mysql.AlterSpecificationDropIndex{}
			spec.Name = def.Name
			specs = append(specs, &spec)
		case *mysql.CreateDefinitionIndex:
			spec := mysql.AlterSpecificationDropIndex{}
			spec.Name = def.Name
			specs = append(specs, &spec)
		default:
		}
	}
	for _, def := range x.Added {
		switch def := def.(type) {
		case *mysql.CreateDefinitionColumn:
			spec := mysql.AlterSpecificationAddColumn{}
			spec.ColumnName = def.ColumnName
			spec.ColumnDefinition = def.ColumnDefinition
			specs = append(specs, &spec)
		case *mysql.CreateDefinitionPrimaryIndex:
		case *mysql.CreateDefinitionUniqueIndex:
			spec := mysql.AlterSpecificationAddIndex{}
			spec.Name = def.Name
			spec.Columns = def.Columns
			spec.Unique = true
			specs = append(specs, &spec)
		case *mysql.CreateDefinitionIndex:
			spec := mysql.AlterSpecificationAddIndex{}
			spec.Name = def.Name
			spec.Columns = def.Columns
			spec.Unique = false
			specs = append(specs, &spec)
		default:
		}
	}
	stmt := new(mysql.AlterTableStatement)
	stmt.TableName = x.Before.TableName
	stmt.AlterSpecifications = specs

	return stmt.ToQuery()
}

func convertCreateTableStatement(stmt *mysql.CreateTableStatement) mysql.Statement {
	return stmt
}

func convertDropTableStatement(stmt *mysql.CreateTableStatement) mysql.Statement {
	newStmt := mysql.DropTableStatement{[]mysql.TableNameIdentifier{}}
	newStmt.TableNames = append(newStmt.TableNames, stmt.TableName)
	return &newStmt
}
