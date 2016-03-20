package importer_test

import (
	//"fmt"
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
	/*osm, err := db.ReadEverythingWithinCoordinates(43.731341, 7.421213, 43.729883, 7.422897)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Stats of ReadEverythingWithinCoordinates between 43.731341 N, 7.421213 E, 43.729883 N, 7.422897 E\n"+
		"Nodes:%d\nWay:%d\n", len(osm.Nodes), len(osm.Ways))*/
	db.Close()
	os.Remove(TESTDATABASE)
}
