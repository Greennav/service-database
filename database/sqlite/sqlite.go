package sqlite

import (
  "database/sql"
  "io/ioutil"
  "log"
  
  "github.com/GreenNav/service-database/database"  
  "github.com/omniscale/imposm3/element"
  _ "github.com/mattn/go-sqlite3" // SQLite driver
)

type SQLiteDatabase struct {
  FileName string
  Database *sql.DB
}

func CreateEmpty(Name string) (*SQLiteDatabase, error) {
  schema, err := ioutil.ReadFile("./schema.sql")
  if err != nil {
    log.Fatalf("Could not find or access schema.sql to create an empty database from: %v", err.Error())
  }
  db, err := sql.Open("sqlite3", Name)
  if err != nil {
    log.Fatal(err)
  }
  _, err = db.Exec(string(schema))
  if err != nil {
    log.Fatal(err)
  }
  
  return &SQLiteDatabase{FileName: Name, Database: db}, err
}

func GetByName(Name string) (*SQLiteDatabase, error) {
  db, err := sql.Open("sqlite3", Name)
  return &SQLiteDatabase{FileName: Name, Database: db}, err
}

func (s SQLiteDatabase) WriteNodes(Nodes chan []element.Node) error {
  tx, err := s.Database.Begin()
  if err != nil {
    return err
  }
  
  sqlInsertNode, err := tx.Prepare(`insert into nodes(id, lon, lat) values(?, ?, ?)`)
  if err != nil {
    return err
  }
  defer sqlInsertNode.Close()
  
  //sqlInsertTags, err := tx.Prepare(`insert into node_tags(ref, key, value) values(?, ?, ?)`)
  if err != nil {
    return err
  }
  defer sqlInsertNode.Close()
  
  for nodes := range Nodes {
    for _, node := range nodes {
      _, err = sqlInsertNode.Exec(node.OSMElem.Id, node.Long, node.Lat)
      if err != nil {
        return err
      }
      
      /*
      for key, value := range node.OSMElem.Tags {
        _, err := sqlInsertTags.Exec(node.OSMElem.Id, key, value)
        if err != nil {
          return err
        }
      }
      */
    }
  }
  tx.Commit()
  
  return nil
}

func (s SQLiteDatabase) WriteWays(Ways chan []element.Way) error {

  return nil
}

func (s SQLiteDatabase) WriteRelations(Relations chan []element.Relation) error {

  return nil
}

func (s SQLiteDatabase) GetEverythingWithinCoordinates(FromLong, FromLat, ToLong, ToLat int) (*database.OSMData, error) {

  return nil, nil
}