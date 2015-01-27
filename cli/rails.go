package main

import (
	"io/ioutil"
	"os"
	"gopkg.in/yaml.v2"
)

type RubyOnRailsFramework struct {
}

func (f *RubyOnRailsFramework) DetectFramework() bool {
	file, err := os.Open(f.GetSchemaFile())
	defer file.Close()
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *RubyOnRailsFramework) GetSchemaFile() string {
	return "db/schema.sql"
}

func (f *RubyOnRailsFramework) GetDatabaseConnectionSettingFile() string {
	return "config/database.yml"
}

func (f *RubyOnRailsFramework) loadDBConfig() (string, error) {
	buf, err := ioutil.ReadFile(f.GetDatabaseConnectionSettingFile())
	if err != nil {
		return "", err
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		return "", err
	}

	return "TODO", nil
}

func (f *RubyOnRailsFramework) GetSchema() (string, error) {
	bytes, err := ioutil.ReadFile(f.GetSchemaFile())
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}

func (f *RubyOnRailsFramework) GetCurrentSchema() (string, error) {
	return "", nil
}
