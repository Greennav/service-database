package importer_test

import (
	"github.com/GreenNav/service-database/database"
	"github.com/GreenNav/service-database/database/sqlite"
	"github.com/GreenNav/service-database/importer"
	"os"
	"testing"
)

const (
	TESTDATABASE = "test_importer.db"
)

func TestWriteToDatabase(t *testing.T) {
	db, _ := sqlite.CreateEmpty(TESTDATABASE)
	importer.WriteToDatabase("./monaco-20150428.osm.pbf", database.OSMDatabase(db))
	os.Remove(TESTDATABASE)
}
