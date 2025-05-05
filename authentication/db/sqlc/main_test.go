package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/demola234/authentication/pkg/utils"

	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	configs, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	var dbErr error
	testDB, dbErr = sql.Open(configs.DBDriver, configs.DBSource)
	if dbErr != nil {
		log.Fatal("cannot connect to db: ", dbErr)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
