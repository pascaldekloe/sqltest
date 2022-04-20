// Package sqltest provides utilities for integration tests.
//
// Tests can run isolated in a dedicated transaction. An automated rollback
// after each test keeps the data consistent and/or clean.
package sqltest

import (
	"database/sql"
	"os"
	"sync"
	"testing"
)

// argument values for sql.Open
var driverName, dataSourceName string

// Setup configures a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a database
// name and connection information.
func Setup(driver, dataSource string) {
	driverName = driver
	dataSourceName = dataSource
}

var driverNameVar, dataSourceNameVar string

// EnvSetup configures the datasource with environment variables. When the
// respective variable (name) is present then its value overrides Setup.
func EnvSetup(driverVar, dataSourceVar string) {
	driverNameVar = driverVar
	dataSourceNameVar = dataSourceVar
}

// Open returns a new connection to the database. Use NewTx instead if possible.
func Open(tb testing.TB) *sql.DB {
	tb.Helper()

	driver := driverName
	dataSource := dataSourceName

	if driverNameVar != "" {
		s, ok := os.LookupEnv(driverNameVar)
		if !ok {
			if driver == "" {
				tb.Fatalf("sqltest: need environment variable %q (with a driver name)", driverNameVar)
			}
		} else {
			if driver != "" {
				tb.Logf("sqltest: driver %q override with environment variable %q", driver, driverNameVar)
			}
			driver = s
		}
	}

	if dataSourceNameVar != "" {
		s, ok := os.LookupEnv(dataSourceNameVar)
		if !ok {
			if dataSource == "" {
				tb.Fatalf("sqltest: need environment variable %q (with a data source name)", dataSourceNameVar)
			}
		} else {
			if dataSource != "" {
				tb.Logf("sqltest: data source %q override with environment variable %q", dataSource, dataSourceNameVar)
			}
			dataSource = s
		}
	}

	d, err := sql.Open(driver, dataSource)
	if err != nil {
		tb.Fatalf("sqltest: driver %q datasource %q unavailable: %s", driver, dataSource, err)
	}
	return d
}

var dBMutex sync.Mutex
var dB *sql.DB

func getDB(t *testing.T) *sql.DB {
	dBMutex.Lock()
	defer dBMutex.Unlock()
	if dB == nil {
		dB = Open(t)
	}
	return dB
}

// NewTx returns a transaction with an automated rollback that fires after the
// test and its subtests complete. The test is skipped when in short mode.
//
// BUG(pascaldekloe): DDL on MySQL causes an an implicit commit, which breaks
// the automated rollback.
func NewTx(t *testing.T) *sql.Tx {
	t.Helper()

	if testing.Short() {
		t.Skip("sqltest: no DB in short mode")
	}

	tx, err := getDB(t).Begin()
	if err != nil {
		t.Fatal("sqltest: transaction launch:", err)
	}
	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil {
			t.Error("sqltest: automatic rollback:", err)
		}
	})
	return tx
}
