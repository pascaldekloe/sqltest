// Package sqltest provides convenience with integration testing.
// Tests using this package are skipped go test -short
package sqltest

import (
	"database/sql"
	"os"
	"sync"
	"testing"
)

var driverName, dataSourceName string

// Setup sets the data source configuration
// with a driver name and a connect string.
func Setup(driver, dataSource string) {
	driverName = driver
	dataSourceName = dataSource
}

var driverNameVar, dataSourceNameVar string

// EnvSetup causes the data source configuration to be read from
// environment variables (if the respective name is not empty).
func EnvSetup(driverVar, dataSourceVar string) {
	driverNameVar = driverVar
	dataSourceNameVar = dataSourceVar
}

func connect(t *testing.T) *sql.DB {
	if driverNameVar != "" {
		driverName = os.Getenv(driverNameVar)
		if driverName == "" {
			t.Fatalf("sqltest: need environment variable %q (with a driver name)", driverNameVar)
		}
	}
	if dataSourceNameVar != "" {
		dataSourceName = os.Getenv(dataSourceNameVar)
		if dataSourceName == "" {
			t.Fatalf("sqltest: need environment variable %q (with a connect string)", dataSourceNameVar)
		}
	}

	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		t.Fatalf("sqltest: %q on %q unavailable: %s", dataSourceName, driverName, err)
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

// NewTx returns a transaction with automatic rollback after the test
// and its subtests complete. The test is skipped when in short mode.
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
