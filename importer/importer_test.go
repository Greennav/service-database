package importer_test

import (
	"os"
	"testing"

	"github.com/GreenNav/service-database/database"
	"github.com/GreenNav/service-database/database/sqlite"
	"github.com/GreenNav/service-database/importer"
)

const (
	TESTDATABASE = "test_importer.db"
)

func TestWriteToSQLiteDatabase(t *testing.T) {
	db, err := sqlite.CreateEmpty(TESTDATABASE, "../schemata/sqlite.sql")
	if err != nil {
		t.Error(err)
	}
	err = importer.WriteToDatabase("./monaco.osm.pbf", database.OSMDatabase(db))
	if err != nil {
		t.Error(err)
	}
	os.Remove(TESTDATABASE)
}
