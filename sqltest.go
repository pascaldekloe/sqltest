// Package sqltest provides utilities for integration tests.
//
// Tests can run isolated in a dedicated transaction. An automated rollback
// after each test keeps the data consistent and/or clean.
package sqltest

// BUG(pascaldekloe): DDL on MySQL causes an an implicit commit, which breaks
// the automated rollback.

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

func connect(t *testing.T) *sql.DB {
	if driverNameVar != "" {
		s, ok := os.LookupEnv(driverNameVar)
		if !ok {
			if driverName == "" {
				t.Fatalf("sqltest: need environment variable %q (with a driver name)", driverNameVar)
			}
		} else {
			if driverName != "" {
				t.Logf("sqltest: driver %q override with environment variable %q", driverName, driverNameVar)
			}
			driverName = s
		}
	}
	if dataSourceNameVar != "" {
		s, ok := os.LookupEnv(dataSourceNameVar)
		if !ok {
			if dataSourceName == "" {
				t.Fatalf("sqltest: need environment variable %q (with a data source name)", dataSourceNameVar)
			}
		} else {
			if dataSourceName != "" {
				t.Logf("sqltest: data source %q override with environment variable %q", dataSourceName, dataSourceNameVar)
			}
			dataSourceName = s
		}
	}

	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		t.Fatalf("sqltest: driver %q datasource %q unavailable: %s", driverName, dataSourceName, err)
	}
	return d
}

var dBMutex sync.Mutex
var dB *sql.DB

func getDB(t *testing.T) *sql.DB {
	dBMutex.Lock()
	defer dBMutex.Unlock()
	if dB == nil {
		dB = connect(t)
	}
	return dB
}

// NewTx returns a transaction with an automated rollback that fires after the
// test and its subtests complete. The test is skipped when in short mode.
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
