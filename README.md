# SQL Test

… convenience library for the Go programming language. Tests can run isolated in
a dedicated transaction. An automated rollback after each test keeps the data
consistent and/or clean.

[![API](https://pkg.go.dev/badge/github.com/pascaldekloe/sqltest.svg)](https://pkg.go.dev/github.com/pascaldekloe/sqltest)

```go
func init() {
	// database configuration
	sqltest.Setup("pgx", "host=localhost user=test database=postgres")
	// optional connect string override with an environment variable
	sqltest.EnvSetup("", "TEST_CONNECT_STRING")
}

func TestWithSQLInteraction(t *testing.T) {
	db := sqltest.NewTx(t)
	// install package variables
	DBExec = db.ExecContext
	DBQuery = db.QueryContext

	// test …

	// automatic rollback
}
```

This is free and unencumbered software released into the
[public domain](https://creativecommons.org/publicdomain/zero/1.0).
