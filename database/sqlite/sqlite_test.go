package sqlite_test

import (
  "testing"
  "os"
  
  _ "github.com/qedus/osmpbf"
  "github.com/GreenNav/service-database/database/sqlite"
)

func TestCreateEmpty(t *testing.T) {
  _, err := sqlite.CreateEmpty("./test_deleteme.db")
  
  if err != nil {
    t.Error("An error occurred while creating a new SQLite database:\n" + err.Error())
  }
  
  os.Remove("./test_deleteme.db")
}

func TestGetByName(t *testing.T) {
  sqlite.CreateEmpty("test_canyougetme.db")
  _, err := sqlite.GetByName("test_canyougetme.db")
  if err != nil {
    t.Error("Could not get a SQLite database by name:\n" + err.Error())
  }
  
  os.Remove("test_canyougetme.db")
}

func TestWriteNode(t *testing.T) {
  t.FailNow()
}

func TestWriteWay(t *testing.T) {
  t.FailNow()
}

func TestWriteRelation(t *testing.T) {
  t.FailNow()
}

func TestGetNodeFromID(t *testing.T) {
  t.FailNow()
}

func TestGetWayFromID(t *testing.T) {
  t.FailNow()
}

func TestGetRelationFromID(t *testing.T) {
  t.FailNow()
}

func TestGetEverythingWithinCoordinates(t *testing.T) {
  t.FailNow()
}