package db

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBDriver is an enum-like type for supported database drivers.
type DBDriver string

const (
	MySQL      DBDriver = "mysql"
	Postgres   DBDriver = "postgres"
	PostgreSQL DBDriver = "postgresql"
	SQLite     DBDriver = "sqlite"
	SQLServer  DBDriver = "sqlserver"
)

// DBConfig is a user friendly config that will be used to build the DSN.
// It supports MySQL, Postgres, SQLite and SQL Server.
type DBConfig struct {
	// Database driver: MySQL, Postgres, SQLite, SQLServer
	Driver DBDriver
	// Host for networked DBs (ignored for SQLite)
	Host string
	// Port for networked DBs (ignored for SQLite)
	Port int
	// Username for authentication
	User string
	// Password for authentication
	Password string
	// Database name (or SQLite file path, or ":memory:")
	DBName string
	// Driver-specific DSN parameters (optional)
	Params map[string]string

	// Connection pool settings:
	// Maximum number of open connections to the database (default: 10)
	MaxOpenConns int
	// Maximum number of idle connections in the pool (default: 5)
	MaxIdleConns int
	// Maximum lifetime of a connection in seconds (default: 3600, i.e. 1 hour)
	ConnMaxLifeSec int
}

// dsn builds a driver-specific DSN string from DBConfig.
func (c DBConfig) dsn() string {
	driver := strings.ToLower(string(c.Driver))
	switch driver {
	case string(MySQL):
		return c.mysqlDSN()
	case string(Postgres), string(PostgreSQL):
		return c.postgresDSN()
	case string(SQLite):
		return c.sqliteDSN()
	case string(SQLServer):
		return c.sqlserverDSN()
	default:
		panic("unsupported driver: " + string(c.Driver))
	}
}

func (c DBConfig) mysqlDSN() string {
	if c.Params == nil {
		c.Params = map[string]string{}
	}
	if _, ok := c.Params["charset"]; !ok {
		c.Params["charset"] = "utf8mb4"
	}
	if _, ok := c.Params["parseTime"]; !ok {
		c.Params["parseTime"] = "True"
	}
	if _, ok := c.Params["loc"]; !ok {
		c.Params["loc"] = "Local"
	}
	pairs := []string{}
	for k, v := range c.Params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	params := ""
	if len(pairs) > 0 {
		params = "?" + strings.Join(pairs, "&")
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", c.User, c.Password, c.Host, c.Port, c.DBName, params)
}

func (c DBConfig) postgresDSN() string {
	parts := []string{fmt.Sprintf("host=%s", c.Host)}
	if c.Port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", c.Port))
	}
	if c.User != "" {
		parts = append(parts, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", c.Password))
	}
	if c.DBName != "" {
		parts = append(parts, fmt.Sprintf("dbname=%s", c.DBName))
	}
	for k, v := range c.Params {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, " ")
}

func (c DBConfig) sqliteDSN() string {
	return c.DBName
}

func (c DBConfig) sqlserverDSN() string {
	pairs := []string{}
	for k, v := range c.Params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	params := ""
	if len(pairs) > 0 {
		params = "&" + strings.Join(pairs, "&")
	}
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s%s", c.User, c.Password, c.Host, c.Port, c.DBName, params)
}

// Connect opens a gorm.DB connection using DBConfig and configures the pool.
func Connect(cfg DBConfig) (*gorm.DB, error) {
	dialector, err := getDialector(cfg)
	if err != nil {
		return nil, err
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	gdb, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("gormr: failed to open DB: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("gormr: failed to get sql.DB: %w", err)
	}

	// Set default values if not provided
	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 10
	}
	sqlDB.SetMaxOpenConns(maxOpen)

	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 5
	}
	sqlDB.SetMaxIdleConns(maxIdle)

	maxLife := cfg.ConnMaxLifeSec
	if maxLife <= 0 {
		maxLife = 3600 // 1 hour
	}
	sqlDB.SetConnMaxLifetime(time.Duration(maxLife) * time.Second)

	return gdb, nil
}

func getDialector(cfg DBConfig) (gorm.Dialector, error) {
	driver := strings.ToLower(string(cfg.Driver))
	switch driver {
	case string(MySQL):
		return getMySQLDialector(cfg)
	case string(Postgres), string(PostgreSQL):
		return getPostgresDialector(cfg)
	case string(SQLite):
		return getSQLiteDialector(cfg)
	case string(SQLServer):
		return getSQLServerDialector(cfg)
	default:
		return nil, fmt.Errorf("gormr: unsupported driver: %s", cfg.Driver)
	}
}

func getMySQLDialector(cfg DBConfig) (gorm.Dialector, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("gormr: missing Host for MySQL connection")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("gormr: missing Port for MySQL connection")
	}
	if cfg.User == "" {
		return nil, fmt.Errorf("gormr: missing User for MySQL connection")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("gormr: missing DBName for MySQL connection")
	}
	return mysql.Open(cfg.dsn()), nil
}

func getPostgresDialector(cfg DBConfig) (gorm.Dialector, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("gormr: missing Host for Postgres connection")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("gormr: missing Port for Postgres connection")
	}
	if cfg.User == "" {
		return nil, fmt.Errorf("gormr: missing User for Postgres connection")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("gormr: missing DBName for Postgres connection")
	}
	return postgres.Open(cfg.dsn()), nil
}

func getSQLiteDialector(cfg DBConfig) (gorm.Dialector, error) {
	if cfg.DBName == "" {
		return nil, fmt.Errorf("gormr: missing DBName (file path or :memory:) for SQLite connection")
	}
	return sqlite.Open(cfg.dsn()), nil
}

func getSQLServerDialector(cfg DBConfig) (gorm.Dialector, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("gormr: missing Host for SQLServer connection")
	}
	if cfg.Port == 0 {
		return nil, fmt.Errorf("gormr: missing Port for SQLServer connection")
	}
	if cfg.User == "" {
		return nil, fmt.Errorf("gormr: missing User for SQLServer connection")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("gormr: missing DBName for SQLServer connection")
	}
	return sqlserver.Open(cfg.dsn()), nil
}
