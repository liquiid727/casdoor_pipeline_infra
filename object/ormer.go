// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/beego/beego/v2/server/web"
	xormadapter "github.com/casdoor/xorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql" // db = mysql
	_ "github.com/lib/pq"              // db = postgres
	"github.com/liquiid727/pipeline-auth/conf"
	"github.com/liquiid727/pipeline-auth/util"
	_ "github.com/microsoft/go-mssqldb" // db = mssql
	"github.com/xorm-io/xorm"
	"github.com/xorm-io/xorm/core"
	"github.com/xorm-io/xorm/names"
	_ "modernc.org/sqlite" // db = sqlite
)

const (
	defaultConfigPath     = "conf/app.conf"
	defaultExportFilePath = "init_data_dump.json"
)

var (
	ormer          *Ormer = nil
	createDatabase        = true
	configPath            = defaultConfigPath
	exportData            = false
	exportFilePath        = defaultExportFilePath
)

func InitFlag() {
	createDatabasePtr := flag.Bool("createDatabase", false, "true if you need to create database")
	configPathPtr := flag.String("config", defaultConfigPath, "set it to \"/your/path/app.conf\" if your config file is not in: \"/conf/app.conf\"")
	exportDataPtr := flag.Bool("export", false, "export database to JSON file and exit (use -exportPath to specify custom location)")
	exportFilePathPtr := flag.String("exportPath", defaultExportFilePath, "path to the exported data file (used with -export)")
	flag.Parse()

	createDatabase = *createDatabasePtr
	configPath = *configPathPtr
	exportData = *exportDataPtr
	exportFilePath = *exportFilePathPtr

	// Load beego config from the specified config path
	err := web.LoadAppConfig("ini", configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config from %s: %v", configPath, err))
	}
}

func ShouldExportData() bool {
	return exportData
}

func GetExportFilePath() string {
	return exportFilePath
}

func InitConfig() {
	err := web.LoadAppConfig("ini", "../conf/app.conf")
	if err != nil {
		panic(err)
	}

	web.BConfig.WebConfig.Session.SessionOn = true

	InitAdapter()
	CreateTables()
}

func InitAdapter() {
	if conf.GetConfigString("driverName") == "" {
		if !util.FileExist(configPath) {
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			dir = strings.ReplaceAll(dir, "\\", "/")
			panic(fmt.Sprintf("The Casdoor config file: \"app.conf\" was not found, it should be placed at: \"%s/conf/app.conf\"", dir))
		}
	}

	if createDatabase {
		err := createDatabaseForPostgres(conf.GetConfigString("driverName"), conf.GetConfigDataSourceName(), conf.GetConfigString("dbName"))
		if err != nil {
			panic(err)
		}
	}

	var err error
	ormer, err = NewAdapter(conf.GetConfigString("driverName"), conf.GetConfigDataSourceName(), conf.GetConfigString("dbName"))
	if err != nil {
		panic(err)
	}

	tableNamePrefix := conf.GetConfigString("tableNamePrefix")
	tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, tableNamePrefix)
	ormer.Engine.SetTableMapper(tbMapper)
}

func CreateTables() {
	if createDatabase {
		err := ormer.CreateDatabase()
		if err != nil {
			panic(err)
		}
	}

	ormer.createTable()
}

// Ormer represents the MySQL adapter for policy storage.
type Ormer struct {
	driverName     string
	dataSourceName string
	dbName         string
	Db             *sql.DB
	Engine         *xorm.Engine
}

// finalizer is the destructor for Ormer.
func finalizer(a *Ormer) {
	err := a.Engine.Close()
	if err != nil {
		panic(err)
	}

	if a.Db != nil {
		err = a.Db.Close()
		if err != nil {
			panic(err)
		}
	}
}

// NewAdapter is the constructor for Ormer.
func NewAdapter(driverName string, dataSourceName string, dbName string) (*Ormer, error) {
	a := &Ormer{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName
	a.dbName = dbName

	// Open the DB, create it if not existed.
	err := a.open()
	if err != nil {
		return nil, err
	}

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

// NewAdapterFromDb is the constructor for Ormer.
func NewAdapterFromDb(driverName string, dataSourceName string, dbName string, db *sql.DB) (*Ormer, error) {
	a := &Ormer{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName
	a.dbName = dbName
	a.Db = db

	// Open the DB, create it if not existed.
	err := a.openFromDb(a.Db)
	if err != nil {
		return nil, err
	}

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

func refineDataSourceNameForPostgres(dataSourceName string) string {
	reg := regexp.MustCompile(`dbname=[^ ]+\s*`)
	return reg.ReplaceAllString(dataSourceName, "dbname=postgres")
}

func createDatabaseForPostgres(driverName string, dataSourceName string, dbName string) error {
	if driverName == "postgres" {
		db, err := sql.Open(driverName, refineDataSourceNameForPostgres(dataSourceName))
		if err != nil {
			return err
		}
		defer db.Close()

		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\";", dbName))
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
		schema := util.GetValueFromDataSourceName("search_path", dataSourceName)
		if schema != "" {
			db, err = sql.Open(driverName, dataSourceName)
			if err != nil {
				return err
			}
			defer db.Close()

			_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s;", schema))
			if err != nil {
				if !strings.Contains(err.Error(), "already exists") {
					return err
				}
			}
		}

		return nil
	} else {
		return nil
	}
}

func (a *Ormer) CreateDatabase() error {
	if a.driverName == "postgres" {
		return nil
	}

	engine, err := xorm.NewEngine(a.driverName, a.dataSourceName)
	if err != nil {
		return err
	}
	defer engine.Close()

	_, err = engine.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8mb4 COLLATE utf8mb4_general_ci", a.dbName))
	return err
}

func (a *Ormer) open() error {
	dataSourceName := a.dataSourceName + a.dbName
	if a.driverName != "mysql" {
		dataSourceName = a.dataSourceName
	}

	driverName := a.driverName
	if driverName == "sqlite3" {
		driverName = "sqlite"
	}
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return err
	}

	if a.driverName == "postgres" {
		schema := util.GetValueFromDataSourceName("search_path", dataSourceName)
		if schema != "" {
			engine.SetSchema(schema)
		}
	}

	a.Engine = engine
	return nil
}

func (a *Ormer) openFromDb(db *sql.DB) error {
	dataSourceName := a.dataSourceName + a.dbName
	if a.driverName != "mysql" {
		dataSourceName = a.dataSourceName
	}

	xormDb := core.FromDB(db)

	driverName := a.driverName
	if driverName == "sqlite3" {
		driverName = "sqlite"
	}
	engine, err := xorm.NewEngineWithDB(driverName, dataSourceName, xormDb)
	if err != nil {
		return err
	}

	if a.driverName == "postgres" {
		schema := util.GetValueFromDataSourceName("search_path", dataSourceName)
		if schema != "" {
			engine.SetSchema(schema)
		}
	}

	a.Engine = engine
	return nil
}

func (a *Ormer) close() {
	_ = a.Engine.Close()
	a.Engine = nil
}

func (a *Ormer) createTable() {
	showSql := conf.GetConfigBool("showSql")
	a.Engine.ShowSQL(showSql)

	err := a.Engine.Sync2(new(Organization))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Group))
	if err != nil {
		panic(err)
	}

	err = a.runPipelineAuthFieldMigrations()
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(User))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Invitation))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Application))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Provider))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Resource))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Cert))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Key))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Role))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Permission))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Model))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Adapter))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Enforcer))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Session))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Token))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Syncer))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Record))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Webhook))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(WebhookEvent))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(VerificationRecord))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Ldap))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(RadiusAccounting))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(xormadapter.CasbinRule))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Form))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Ticket))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Agent))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Server))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Entry))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Site))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Rule))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(ThirdPartyLink))
	if err != nil {
		panic(err)
	}

	err = a.runPipelineAuthProviderTypeMigration()
	if err != nil {
		panic(err)
	}
}

func (a *Ormer) runPipelineAuthFieldMigrations() error {
	migrations := []struct {
		table     string
		oldColumn string
		newColumn string
		columnDef string
	}{
		{table: "site", oldColumn: "casdoor_application", newColumn: "auth_application", columnDef: "varchar(100)"},
		{table: "user", oldColumn: "casdoor", newColumn: "oidc", columnDef: "varchar(100)"},
	}

	for _, migration := range migrations {
		if err := a.renameColumnIfNeeded(migration.table, migration.oldColumn, migration.newColumn, migration.columnDef); err != nil {
			return err
		}
	}

	return nil
}

func (a *Ormer) runPipelineAuthProviderTypeMigration() error {
	if _, err := a.Engine.Where("category = ? AND type = ?", "OAuth", "Casdoor").Cols("type").Update(&Provider{Type: "OIDC"}); err != nil {
		return err
	}

	return nil
}

func (a *Ormer) renameColumnIfNeeded(table string, oldColumn string, newColumn string, columnDef string) error {
	tableName := names.NewPrefixMapper(names.SnakeMapper{}, conf.GetConfigString("tableNamePrefix")).Obj2Table(table)
	oldExists, err := a.columnExists(tableName, oldColumn)
	if err != nil {
		return err
	}
	if !oldExists {
		return nil
	}

	newExists, err := a.columnExists(tableName, newColumn)
	if err != nil {
		return err
	}
	if newExists {
		return a.copyColumnData(tableName, oldColumn, newColumn)
	}

	sqlText, err := a.buildRenameColumnSQL(tableName, oldColumn, newColumn, columnDef)
	if err != nil {
		return err
	}

	_, err = a.Engine.Exec(sqlText)
	return err
}

func (a *Ormer) columnExists(tableName string, columnName string) (bool, error) {
	columns, err := a.Engine.DBMetas()
	if err != nil {
		return false, err
	}

	for _, table := range columns {
		if table.Name != tableName {
			continue
		}

		for _, col := range table.Columns() {
			if strings.EqualFold(col.Name, columnName) {
				return true, nil
			}
		}

		return false, nil
	}

	return false, nil
}

func (a *Ormer) buildRenameColumnSQL(tableName string, oldColumn string, newColumn string, columnDef string) (string, error) {
	switch a.driverName {
	case "postgres", "sqlite3", "sqlite":
		return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", tableName, oldColumn, newColumn), nil
	case "mysql":
		return fmt.Sprintf("ALTER TABLE %s CHANGE %s %s %s", tableName, oldColumn, newColumn, columnDef), nil
	case "mssql":
		return fmt.Sprintf("EXEC sp_rename '%s.%s', '%s', 'COLUMN'", tableName, oldColumn, newColumn), nil
	default:
		return "", fmt.Errorf("unsupported database driver for column rename: %s", a.driverName)
	}
}

func (a *Ormer) copyColumnData(tableName string, oldColumn string, newColumn string) error {
	sqlText := fmt.Sprintf(
		"UPDATE %s SET %s = %s WHERE (%s IS NULL OR %s = '') AND (%s IS NOT NULL AND %s <> '')",
		tableName,
		newColumn,
		oldColumn,
		newColumn,
		newColumn,
		oldColumn,
		oldColumn,
	)
	_, err := a.Engine.Exec(sqlText)
	return err
}
