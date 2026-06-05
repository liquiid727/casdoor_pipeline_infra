package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/liquiid727/pipeline-auth/conf"
	"github.com/liquiid727/pipeline-auth/util"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

type finding struct {
	level   string
	title   string
	details string
}

func main() {
	configPath := flag.String("config", "conf/app.conf", "path to app.conf")
	flag.Parse()

	if err := web.LoadAppConfig("ini", *configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(2)
	}

	driverName := conf.GetConfigString("driverName")
	dataSourceName := conf.GetConfigDataSourceName()
	dbName := conf.GetConfigString("dbName")
	if driverName == "" || dataSourceName == "" {
		fmt.Fprintln(os.Stderr, "driverName or dataSourceName is empty")
		os.Exit(2)
	}

	db, err := openDB(driverName, dataSourceName, dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(2)
	}
	defer db.Close()

	tablePrefix := conf.GetConfigString("tableNamePrefix")

	findings := []finding{}

	findings = append(findings, checkLegacyColumn(db, driverName, tablePrefix+"site", "casdoor_application", "auth_application")...)
	findings = append(findings, checkLegacyColumn(db, driverName, tablePrefix+"user", "casdoor", "oidc")...)
	findings = append(findings, checkProviderCount(db, tablePrefix+"provider", "warn", "Legacy OAuth provider type", "category = 'OAuth' AND type = 'Casdoor'")...)
	findings = append(findings, checkProviderCount(db, tablePrefix+"provider", "error", "Unsupported Storage provider type", "category = 'Storage' AND type = 'Casdoor'")...)
	findings = append(findings, checkProviderCount(db, tablePrefix+"provider", "error", "Unsupported Log provider type", "type = 'Casdoor Permission Log'")...)

	hasError := false
	fmt.Println("Pipeline Auth de-Casdoor preflight")
	fmt.Println()
	for _, item := range findings {
		fmt.Printf("[%s] %s\n", strings.ToUpper(item.level), item.title)
		fmt.Printf("  %s\n", item.details)
		if item.level == "error" {
			hasError = true
		}
	}

	if len(findings) == 0 {
		fmt.Println("[OK] No legacy Casdoor blockers detected.")
		return
	}

	if hasError {
		os.Exit(1)
	}
}

func openDB(driverName string, dataSourceName string, dbName string) (*sql.DB, error) {
	driver := driverName
	dsn := dataSourceName
	if driverName == "mysql" {
		dsn = dataSourceName + dbName
	}
	if driver == "sqlite3" {
		driver = "sqlite"
	}
	return sql.Open(driver, dsn)
}

func checkLegacyColumn(db *sql.DB, driverName string, tableName string, oldColumn string, newColumn string) []finding {
	findings := []finding{}

	oldExists, err := columnExists(db, driverName, tableName, oldColumn)
	if err != nil {
		return append(findings, finding{
			level:   "error",
			title:   fmt.Sprintf("Column inspection failed for %s.%s", tableName, oldColumn),
			details: err.Error(),
		})
	}

	newExists, err := columnExists(db, driverName, tableName, newColumn)
	if err != nil {
		return append(findings, finding{
			level:   "error",
			title:   fmt.Sprintf("Column inspection failed for %s.%s", tableName, newColumn),
			details: err.Error(),
		})
	}

	switch {
	case oldExists && newExists:
		findings = append(findings, finding{
			level:   "warn",
			title:   fmt.Sprintf("Both legacy and new columns exist on %s", tableName),
			details: fmt.Sprintf("Found both %s and %s. Startup migration should copy forward data, but verify old column can be retired safely.", oldColumn, newColumn),
		})
	case oldExists && !newExists:
		findings = append(findings, finding{
			level:   "warn",
			title:   fmt.Sprintf("Legacy column still present on %s", tableName),
			details: fmt.Sprintf("Found %s without %s. Startup migration is expected to rename this column.", oldColumn, newColumn),
		})
	case !oldExists && !newExists:
		findings = append(findings, finding{
			level:   "error",
			title:   fmt.Sprintf("Expected column missing on %s", tableName),
			details: fmt.Sprintf("Neither %s nor %s exists. Inspect schema before rollout.", oldColumn, newColumn),
		})
	}

	return findings
}

func checkProviderCount(db *sql.DB, tableName string, level string, title string, where string) []finding {
	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tableName, where))
	count := 0
	if err := row.Scan(&count); err != nil {
		return []finding{{
			level:   "error",
			title:   title,
			details: fmt.Sprintf("failed to count providers: %v", err),
		}}
	}
	if count == 0 {
		return nil
	}

	return []finding{{
		level:   level,
		title:   title,
		details: fmt.Sprintf("Found %d matching provider record(s): %s", count, where),
	}}
}

func columnExists(db *sql.DB, driverName string, tableName string, columnName string) (bool, error) {
	switch driverName {
	case "postgres":
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name = $1 AND column_name = $2
			)
		`, tableName, columnName).Scan(&exists)
		return exists, err
	case "mysql":
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?
			)
		`, tableName, columnName).Scan(&exists)
		return exists, err
	case "sqlite", "sqlite3":
		rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
		if err != nil {
			return false, err
		}
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name string
			var typ string
			var notNull int
			var dflt sql.NullString
			var pk int
			if err := rows.Scan(&cid, &name, &typ, &notNull, &dflt, &pk); err != nil {
				return false, err
			}
			if strings.EqualFold(name, columnName) {
				return true, nil
			}
		}
		return false, rows.Err()
	case "mssql":
		var exists bool
		err := db.QueryRow(`
			SELECT CASE WHEN EXISTS (
				SELECT 1
				FROM INFORMATION_SCHEMA.COLUMNS
				WHERE TABLE_NAME = @p1 AND COLUMN_NAME = @p2
			) THEN 1 ELSE 0 END
		`, tableName, columnName).Scan(&exists)
		return exists, err
	default:
		return false, fmt.Errorf("unsupported driver %q", driverName)
	}
}

func init() {
	if conf.GetConfigString("driverName") == "" && util.FileExist("conf/app.conf") {
		_ = web.LoadAppConfig("ini", "conf/app.conf")
	}
}
