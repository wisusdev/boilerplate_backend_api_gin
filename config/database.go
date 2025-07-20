package config

type Connections struct {
	Driver string `json:"driver"`
	SQLite Sqlite `json:"sqlite"`
	MySQL  Mysql  `json:"mysql"`
	PgSQL  Pgsql  `json:"pgsql"`
	Redis  Redis  `json:"redis"`
}

type Sqlite struct {
	Driver   string
	Database string
}

type Mysql struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Charset  string
}

type Pgsql struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

func (c *Connections) GetSQLiteConfig() Sqlite {
	return c.SQLite
}

func (c *Connections) GetMySQLConfig() Mysql {
	return c.MySQL
}

func (c *Connections) GetPgSQLConfig() Pgsql {
	return c.PgSQL
}

func (c *Connections) GetRedisConfig() Redis {
	return c.Redis
}

func DatabaseConfig() *Connections {
	return &Connections{
		Driver: GetEnv("DB_DRIVER", "sqlite"),

		SQLite: Sqlite{
			Driver:   "sqlite3",
			Database: "database.db",
		},
		MySQL: Mysql{
			Driver:   "mysql",
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "3306"),
			Database: GetEnv("DB_NAME", "semita"),
			Username: GetEnv("DB_USER", "root"),
			Password: GetEnv("DB_PASSWORD", ""),
			Charset:  GetEnv("DB_CHARSET", "utf8mb4"),
		},
		PgSQL: Pgsql{
			Driver:   "postgres",
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			Database: GetEnv("DB_NAME", "semita"),
			Username: GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", ""),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		Redis: Redis{
			Host: GetEnv("REDIS_HOST", "localhost"),
			Port: GetEnv("REDIS_PORT", "6379"),
		},
	}
}
