// vim: noet sw=8 sts=8
%{
package mysql

import (
    "fmt"
    "strconv"
    "errors"
)

type Token struct {
    tok int
    lit string
    pos Position
}

%}

%union{
    statements []Statement
    statement Statement
    table_names []TableNameIdentifier
    table_name TableNameIdentifier
    table_options []TableOption
    table_option TableOption
    database_name DatabaseNameIdentifier
    column_name ColumnNameIdentifier
    column_names []ColumnNameIdentifier
    index_name IndexNameIdentifier
    column_definition ColumnDefinition
    alter_specifications []AlterSpecification
    alter_specification AlterSpecification
    data_type DataTypeDefinition
    create_definitions []CreateDefinition
    create_definition CreateDefinition
    bool bool
    data_type_type DataType
    default_definition DefaultDefinition
    uint uint
    fraction_option [2]uint
    tok       Token
    str string
}

%type<statements> statements
%type<statement> statement
%type<table_names> table_names
%type<table_name> table_name
%type<database_name> database_name
%type<column_name> column_name
%type<column_names> index_column_names
%type<index_name> index_name skipable_index_name
%type<column_definition> column_definition
%type<alter_specifications> alter_specifications
%type<alter_specification> alter_specification
%type<create_definition> create_definition
%type<create_definitions> create_definitions
%type<data_type> data_type
%type<data_type_type> data_type_number data_type_fraction data_type_decimal
%type<bool> unsigned_option zerofill_option nullable autoincrement
%type<uint> length_option
%type<fraction_option> fraction_option decimal_option
%type<default_definition> default
%type<table_option> table_option
%type<table_options> skipable_table_options
%type<str> storage_engine_name string

%token<tok> IDENT NUMBER RAW COMMENT_START COMMENT_FINISH
%token<tok> DROP CREATE ALTER ADD
%token<tok> TABLE COLUMN DATABASE INDEX KEY NOT NULL AUTO_INCREMENT DEFAULT CURRENT_TIMESTAMP ON UPDATE PRIMARY UNIQUE
%token<tok> USING BTREE HASH CHARSET CHARACTER SET COLLATE
%token<tok> ENGINE AVG_ROW_LENGTH CHECKSUM COMMENT KEY_BLOCK_SIZE MAX_ROWS MIN_ROWS ROW_FORMAT DYNAMIC FIXED COMPRESSED REDUNDANT COMPACT
%token<tok> BIT TINYINT SMALLINT MEDIUMINT INT INTEGER BIGINT REAL DOUBLE FLOAT DECIMAL NUMERIC DATE TIME TIMESTAMP DATETIME YEAR CHAR VARCHAR BINARY VARBINARY TINYBLOB BLOB MEDIUMBLOB LONGBLOB TINYTEXT TEXT MEDIUMTEXT LONGTEXT UNSIGNED ZEROFILL

%%

statements
    :
    {
        $$ = nil
        if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
            l.statements = $$
        }
    }
    | statements statement
    {
        $$ = append($1, $2)
        if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
            l.statements = $$
        }
    }

statement
    : DROP TABLE table_names ';'
    {
        $$ = &DropTableStatement{TableNames: $3}
    }
    | DROP DATABASE database_name ';'
    {
        $$ = &DropDatabaseStatement{DatabaseName: $3}
    }
    | CREATE DATABASE database_name ';'
    {
        $$ = &CreateDatabaseStatement{DatabaseName: $3}
    }
    | CREATE TABLE table_name '(' create_definitions ')' skipable_table_options optional_statement_finish
    {
        $$ = &CreateTableStatement{TableName: $3, CreateDefinitions: $5, TableOptions: $7}
    }
    | ALTER TABLE table_name alter_specifications ';'
    {
        $$ = &AlterTableStatement{TableName: $3, AlterSpecifications: $4}
    }
    | COMMENT_START RAW COMMENT_FINISH ';'
    {
        $$ = &CommentStatement{$2.lit}
    }

optional_statement_finish
    :
    | ';'

create_definitions
    : create_definition
    {
        $$ = []CreateDefinition{$1}
    }
    | create_definitions ',' create_definition
    {
        $$ = append($1, $3)
    }

create_definition
    : column_name column_definition
    {
        $$ = &CreateDefinitionColumn{ColumnName: $1, ColumnDefinition: $2}
    }
    | PRIMARY KEY skipable_index_type '(' index_column_names ')'
    {
        $$ = &CreateDefinitionPrimaryIndex{Columns: $5}
    }
    | index_or_key skipable_index_name skipable_index_type '(' index_column_names ')' skipable_index_type
    {
        $$ = &CreateDefinitionIndex{Name: $2, Columns: $5}
    }
    | UNIQUE index_or_key skipable_index_name skipable_index_type '(' index_column_names ')'
    {
        $$ = &CreateDefinitionUniqueIndex{Name: $3, Columns: $6}
    }

skipable_table_options
 :
 {
    $$ = []TableOption{}
 }
 | skipable_table_options table_option
 {
    $$ = append($1, $2)
 }

table_option
    : ENGINE skipable_equal storage_engine_name
    {
        var option TableOption
        option.Key = "ENGINE"
        option.Value = $3
        $$ = option
    }
    | AUTO_INCREMENT skipable_equal NUMBER
    {
        var option TableOption
        option.Key = "AUTO_INCREMENT"
        option.Value = $3.lit
        $$ = option
    }
    | AVG_ROW_LENGTH skipable_equal NUMBER
    {
        var option TableOption
        option.Key = "AVG_ROW_LENGTH"
        option.Value = $3.lit
        $$ = option
    }
    | CHECKSUM skipable_equal NUMBER
    {
        var option TableOption
        option.Key = "CHECKSUM"
        option.Value = $3.lit
        $$ = option
    }
    | COMMENT skipable_equal '\'' RAW '\''
    {
        var option TableOption
        option.Key = "COMMENT"
        option.Value = $4.lit
        $$ = option
    }
    | KEY_BLOCK_SIZE skipable_equal NUMBER
    {
        var option TableOption
        option.Key = "KEY_BLOCK_SIZE"
        option.Value = $3.lit
        $$ = option
    }
    | MAX_ROWS skipable_equal NUMBER
    {
        var option TableOption
        option.Key = "MAX_ROWS"
        option.Value = $3.lit
        $$ = option
    }
    | MIN_ROWS skipable_equal NUMBER
    {
        var option TableOption
        option.Key = "MIN_ROWS"
        option.Value = $3.lit
        $$ = option
    }
    | ROW_FORMAT skipable_equal storage_engine_name
    {
        var option TableOption
        option.Key = "ROW_FORMAT"
        option.Value = $3
        $$ = option
    }
    | DEFAULT CHARSET skipable_equal string
    {
        var option TableOption
        option.Key = "DEFAULT CHARACTER"
        option.Value = $4
        $$ = option
    }
    | COLLATE skipable_equal string
    {
        var option TableOption
        option.Key = "COLLATE"
        option.Value = $3
        $$ = option
    }

charset_or_character_set
    : CHARSET
    | CHARACTER SET

skipable_equal
    :
    | '='

index_column_names
    : column_name
    {
        result := []ColumnNameIdentifier{}
        result = append(result, $1)
        $$ = result
    }
    | index_column_names ',' column_name
    {
        result := append($1, $3)
        $$ = result
    }

skipable_index_type
    :
    | USING BTREE
    | USING HASH

table_names
    : table_name
    {
        $$ = []TableNameIdentifier{$1}
    }
    | table_names ',' table_name
    {
        $$ = append([]TableNameIdentifier{$3}, $1...)
    }

table_name
    : IDENT
    {
        $$ = TableNameIdentifier{Name: $1.lit}
    }
    | '`' RAW '`'
    {
        $$ = TableNameIdentifier{Name: $2.lit}
    }
    | IDENT '.' IDENT
    {
        $$ = TableNameIdentifier{Database: $1.lit, Name: $3.lit}
    }

database_name
    : IDENT
    {
        $$ = DatabaseNameIdentifier{Name: $1.lit}
    }
    | '`' RAW '`'
    {
        $$ = DatabaseNameIdentifier{Name: $2.lit}
    }

alter_specifications
    :
    {
        $$ = nil
    }
    | alter_specification
    {
        $$ = []AlterSpecification{$1}
    }
    | alter_specifications ',' alter_specification
    {
        $$ = append($1, $3)
    }

alter_specification
    : ADD skipable_column column_name column_definition
    {
        $$ = &AlterSpecificationAddColumn{ColumnName: $3, ColumnDefinition: $4}
    }
    | DROP index_or_key index_name
    {
        $$ = &AlterSpecificationDropIndex{IndexName: $3}
    }
    | DROP skipable_column column_name
    {
        $$ = &AlterSpecificationDropColumn{ColumnName: $3}
    }

skipable_column
    :
    | COLUMN

column_definition
    : data_type nullable default autoincrement key_options column_comment
    {
        $$ = ColumnDefinition{$1, $2, $4, $3}
    }

nullable
    :
    {
        $$ = true
    }
    | NULL
    {
        $$ = true
    }
    | NOT NULL
    {
        $$ = false
    }

default
    :
    {
        $$ = &DefaultDefinitionEmpty{}
    }
    | DEFAULT NULL
    {
        $$ = &DefaultDefinitionNull{}
    }
    | DEFAULT NUMBER
    {
        value := DefaultDefinitionString{}
        value.Value = $2.lit
        $$ = &value
    }
    | DEFAULT '"' RAW '"'
    {
        value := DefaultDefinitionString{}
        value.Value = $3.lit
        $$ = &value
    }
    | DEFAULT '\'' RAW '\''
    {
        value := DefaultDefinitionString{}
        value.Value = $3.lit
        $$ = &value
    }
    | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    {
        $$ = &DefaultDefinitionCurrentTimestamp{true}
    }
    | DEFAULT CURRENT_TIMESTAMP
    {
        $$ = &DefaultDefinitionCurrentTimestamp{false}
    }

string
    : IDENT
    {
        $$ = $1.lit
    }
    | '\'' RAW '\''
    {
        $$ = $2.lit
    }
    | '"' RAW '"'
    {
        $$ = $2.lit
    }

autoincrement
    :
    {
        $$ = false
    }
    | AUTO_INCREMENT
    {
        $$ = true
    }

key_options
    :

column_comment
    :
    | COMMENT string

data_type
    : BIT
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_BIT }
    }
    | data_type_number length_option unsigned_option zerofill_option
    {
        $$ = &DataTypeDefinitionNumber{Type: $1, Length: $2, Unsigned: $3, Zerofill: $4 }
    }
    | data_type_fraction fraction_option unsigned_option zerofill_option
    {
        fraction := $2
        $$ = &DataTypeDefinitionFraction{Type: $1, Length: fraction[0], Decimals: fraction[1], Unsigned: $3, Zerofill: $4 }
    }
    | data_type_decimal decimal_option unsigned_option zerofill_option
    {
        fraction := $2
        $$ = &DataTypeDefinitionFraction{Type: $1, Length: fraction[0], Decimals: fraction[1], Unsigned: $3, Zerofill: $4 }
    }
    | DATE
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_DATE }
    }
    | TIME
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_TIME }
    }
    | TIMESTAMP
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_TIMESTAMP }
    }
    | DATETIME
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_DATETIME }
    }
    | YEAR
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_YEAR }
    }
    | CHAR length_option optional_character_set optional_collate
    {
        $$ = &DataTypeDefinitionString{Type: DATATYPE_CHAR, Length: $2 }
    }
    | VARCHAR length_option optional_character_set optional_collate
    {
        $$ = &DataTypeDefinitionString{Type: DATATYPE_VARCHAR, Length: $2 }
    }
    | BINARY
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_BINARY }
    }
    | VARBINARY
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_VARBINARY }
    }
    | TINYBLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_TINYBLOB }
    }
    | BLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_BLOB }
    }
    | MEDIUMBLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_MEDIUMBLOB }
    }
    | LONGBLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_LONGBLOB }
    }
    | TINYTEXT
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_TINYTEXT }
    }
    | TEXT optional_character_set optional_collate
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_TEXT }
    }
    | MEDIUMTEXT optional_character_set optional_collate
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_MEDIUMTEXT }
    }
    | LONGTEXT optional_character_set optional_collate
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_LONGTEXT }
    }

data_type_number
    : TINYINT
    {
        $$ = DATATYPE_TINYINT
    }
    | SMALLINT
    {
        $$ = DATATYPE_SMALLINT
    }
    | MEDIUMINT
    {
        $$ = DATATYPE_MEDIUMINT
    }
    | INT
    {
        $$ = DATATYPE_INT
    }
    | INTEGER
    {
        $$ = DATATYPE_INT
    }
    | BIGINT
    {
        $$ = DATATYPE_BIGINT
    }

data_type_fraction
    : REAL
    {
        $$ = DATATYPE_REAL
    }
    | DOUBLE
    {
        $$ = DATATYPE_DOUBLE
    }
    | FLOAT
    {
        $$ = DATATYPE_FLOAT
    }

data_type_decimal
    : DECIMAL
    {
        $$ = DATATYPE_DECIMAL
    }
    | NUMERIC
    {
        $$ = DATATYPE_NUMERIC
    }

length_option
    :
    {
        $$ = 0
    }
    | '(' NUMBER ')'
    {
        num, err := strconv.Atoi($2.lit)
        if err != nil {
            num = 0
        }
        $$ = uint(num)
    }

fraction_option
    :
    {
        $$ = [2]uint{0, 0}
    }
    | '(' NUMBER ',' NUMBER ')'
    {
        num1, err := strconv.Atoi($2.lit)
        if err != nil {
            num1 = 0
        }
        num2, err := strconv.Atoi($4.lit)
        if err != nil {
            num2 = 0
        }
        result := [2]uint{0, 0}
        result[0] = uint(num1)
        result[1] = uint(num2)
        $$ = result
    }

optional_character_set
    :
    | CHARACTER SET string

optional_collate
    :
    | COLLATE string

decimal_option
    :
    {
        $$ = [2]uint{0, 0}
    }
    | '(' NUMBER ')'
    {
        result := [2]uint{0, 0}
        num1, err := strconv.Atoi($2.lit)
        if err != nil {
            num1 = 0
        }
        result[0] = uint(num1)
        $$ = result
    }
    | '(' NUMBER ',' NUMBER ')'
    {
        num1, err := strconv.Atoi($2.lit)
        if err != nil {
            num1 = 0
        }
        num2, err := strconv.Atoi($4.lit)
        if err != nil {
            num2 = 0
        }
        result := [2]uint{0, 0}
        result[0] = uint(num1)
        result[1] = uint(num2)
        $$ = result
    }

unsigned_option
    :
    {
        $$ = false
    }
    | UNSIGNED
    {
        $$ = true
    }

zerofill_option
    :
    {
        $$ = false
    }
    | ZEROFILL
    {
        $$ = true
    }


index_or_key
    : INDEX
    | KEY

column_name
    : IDENT
    {
        $$ = ColumnNameIdentifier{Name: $1.lit}
    }
    | '`' RAW '`'
    {
        $$ = ColumnNameIdentifier{Name: $2.lit}
    }

skipable_index_name
    :
    {
        $$ = IndexNameIdentifier{Name: ""}
    }
    | index_name
    {
        $$ = $1
    }

index_name
    : IDENT
    {
        $$ = IndexNameIdentifier{Name: $1.lit}
    }
    | '`' RAW '`'
    {
        $$ = IndexNameIdentifier{Name: $2.lit}
    }

storage_engine_name
    : IDENT
    {
        $$ = $1.lit
    }
    | '\'' IDENT '\''
    {
        $$ = $2.lit
    }
    | '"' IDENT '\''
    {
        $$ = $2.lit
    }

skipable_default
    :
    | DEFAULT

%%

type LexerWrapper struct {
    scanner *Scanner
    recentLit   string
    recentPos   Position
    statements []Statement
}

func (l *LexerWrapper) Lex(lval *yySymType) int {
    tok, lit, pos := l.scanner.Scan()
    if tok == EOF {
        return 0
    }
    lval.tok = Token{tok: tok, lit: lit, pos: pos}
    l.recentLit = lit
    l.recentPos = pos
    return tok
}

func (l *LexerWrapper) Error(e string) {
}

func (l *LexerWrapper) GetError(e string) error {
    result := fmt.Sprintf("%s while processing near %q line %d, col: %d\n", e, l.recentLit, l.recentPos.Line, l.recentPos.Column)
    result += fmt.Sprintf("%s\n", l.scanner.CurrentLine())
    for i := 0; i < l.recentPos.Column-1; i++ {
        result += fmt.Sprintf(" ")
    }
    result += fmt.Sprintf("^\n")
    return errors.New(result)
}

func Parse(s *Scanner) ([]Statement, error) {
    l := LexerWrapper{scanner: s}
    if yyParse(&l) != 0 {
        return []Statement{}, l.GetError("syntax error")
    }
    return l.statements, nil
}
