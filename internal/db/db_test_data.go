package db

// Test cases for DSN string generation
// dsnTestCase define un caso de prueba para la generación de DSN.
type dsnTestCase struct {
	config  DBConfig
	wantDSN string // Expected DSN string after formatting
}

var dsnTestCases = map[string]dsnTestCase{
	"mysql_minimal": {
		config: DBConfig{
			Driver: MySQL,
			Host:   "localhost",
			Port:   3306,
			User:   "root",
			DBName: "test",
		},
		wantDSN: "root:@tcp(localhost:3306)/test?charset=utf8mb4&loc=Local&parseTime=True",
	},
	"mysql_with_password": {
		config: DBConfig{
			Driver:   MySQL,
			Host:     "localhost",
			Port:     3306,
			User:     "user",
			Password: "pass",
			DBName:   "test",
		},
		wantDSN: "user:pass@tcp(localhost:3306)/test?charset=utf8mb4&loc=Local&parseTime=True",
	},
	"postgres_minimal": {
		config: DBConfig{
			Driver: Postgres,
			Host:   "localhost",
			Port:   5432,
			User:   "postgres",
			DBName: "test",
		},
		wantDSN: "host=localhost port=5432 user=postgres dbname=test",
	},
	"sqlite_memory": {
		config: DBConfig{
			Driver: SQLite,
			DBName: ":memory:",
		},
		wantDSN: ":memory:",
	},
	"sqlserver_complete": {
		config: DBConfig{
			Driver:   SQLServer,
			Host:     "localhost",
			Port:     1433,
			User:     "sa",
			Password: "pass",
			DBName:   "test",
		},
		wantDSN: "sqlserver://sa:pass@localhost:1433?database=test",
	},
}

// connectTestCase define un caso de prueba para la función Connect.
type connectTestCase struct {
	config  DBConfig
	wantErr bool   // Whether we expect an error
	errMsg  string // Expected error message if wantErr is true
}

var connectTestCases = map[string]connectTestCase{
	"mysql_missing_host": {
		config: DBConfig{
			Driver: MySQL,
			Port:   3306,
			User:   "root",
			DBName: "test",
		},
		wantErr: true,
		errMsg:  "gormr: missing Host for MySQL connection",
	},
	"postgres_missing_dbname": {
		config: DBConfig{
			Driver: Postgres,
			Host:   "localhost",
			Port:   5432,
			User:   "postgres",
		},
		wantErr: true,
		errMsg:  "gormr: missing DBName for Postgres connection",
	},
	"sqlite_missing_dbname": {
		config: DBConfig{
			Driver: SQLite,
		},
		wantErr: true,
		errMsg:  "gormr: missing DBName (file path or :memory:) for SQLite connection",
	},
	"invalid_driver": {
		config: DBConfig{
			Driver: "invalid",
		},
		wantErr: true,
		errMsg:  "gormr: unsupported driver: invalid",
	},
}
