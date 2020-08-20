# SQL Test

… convenience library for the Go programming language.

```go
func init() {
	// fix driver to PostgreSQL
	sqltest.Setup("postgres", "")
	// read connect string from an environment variable
	sqltest.EnvSetup("", "TEST_CONNECT_STRING")
}

func TestFoo(t *testing.T) {
	// install connection in package
	DBExec = sqltest.NewTx(t).Exec

	…

	// automatic rollback
}
```

This is free and unencumbered software released into the
[public domain](https://creativecommons.org/publicdomain/zero/1.0).
