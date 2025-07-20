package database_connections

import (
	"os"
	"path/filepath"
	"semita/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseConfiguration struct {
	Driver   string
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// DatabaseConnectSQL retorna una conexión SQL a través del adapter
func DatabaseConnectSQL() SQLAdapter {
	var dbConfig = config.DatabaseConfig()

	switch dbConfig.Driver {
	case "mysql":
		return MysqlDatabaseConnectSQL(dbConfig.MySQL)
	case "postgres":
		return PostgresDatabaseConnectSQL(dbConfig.PgSQL)
	case "sqlite":
		return SqliteDatabaseConnectSQL(dbConfig.SQLite)
	default:
		panic("Unsupported database driver: " + dbConfig.Driver)
	}
}

// SQL Direct Connections
func MysqlDatabaseConnectSQL(config config.Mysql) SQLAdapter {
	var driver = config.Driver
	var host = config.Host
	var port = config.Port
	var dbname = config.Database
	var user = config.Username
	var password = config.Password

	db, err := SqlOpen(driver, user+":"+password+"@tcp("+host+":"+port+")/"+dbname)
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return NewSQLAdapter(db)
}

func PostgresDatabaseConnectSQL(config config.Pgsql) SQLAdapter {
	var driver = config.Driver
	var host = config.Host
	var port = config.Port
	var dbname = config.Database
	var user = config.Username
	var password = config.Password

	db, err := SqlOpen(driver, "host="+host+" port="+port+" user="+user+" password="+password+" dbname="+dbname+" sslmode=disable")
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return NewSQLAdapter(db)
}

func SqliteDatabaseConnectSQL(config config.Sqlite) SQLAdapter {
	var driver = config.Driver
	var dbname = config.Database

	var dbPath = filepath.Join("storage", dbname+".db")

	if errorMkdir := os.MkdirAll("storage", os.ModePerm); errorMkdir != nil {
		panic("Error creating storage directory: " + errorMkdir.Error())
	}

	db, err := SqlOpen(driver, dbPath)
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return NewSQLAdapter(db)
}
