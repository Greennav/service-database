package sqlite

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	"github.com/GreenNav/service-database/database"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/omniscale/imposm3/element"
)

type SQLiteDatabase struct {
	FileName string
	Database *sql.DB
	tx       *sql.Tx
}

func (s SQLiteDatabase) getTransaction() *sql.Tx {
	if s.tx == nil {
		s.tx, _ = s.Database.Begin()
	}
	return s.tx
}

func CreateEmpty(Name string) (*SQLiteDatabase, error) {
	schema, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
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

func (s SQLiteDatabase) WriteNodes(Nodes chan element.Node) error {
	stmt, err := s.Database.Prepare("insert into nodes(id, lon, lat) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for node := range Nodes {
		_, err = stmt.Exec(node.OSMElem.Id, node.Long, node.Lat)
		log.Print("node commited")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s SQLiteDatabase) WriteNodeTags(Nodes chan element.Node) error {
	stmt, err := s.Database.Prepare("insert into node_tags(ref, key, value) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for node := range Nodes {
		for key, value := range node.OSMElem.Tags {
			_, err := stmt.Exec(node.OSMElem.Id, key, value)
			log.Print("node tag commited")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s SQLiteDatabase) WriteWays(Ways chan element.Way) error {
	stmt, err := s.Database.Prepare("insert into ways(id) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for way := range Ways {
		_, err = stmt.Exec(way.Id)
		log.Print("way commited")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s SQLiteDatabase) WriteWayNodes(Ways chan element.Way) error {
	stmt, err := s.Database.Prepare("insert into way_nodes(way, num, node) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for way := range Ways {
		for num, node := range way.Nodes {
			_, err = stmt.Exec(way.Id, num, node.Id)
			log.Print("way node commited")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s SQLiteDatabase) WriteWayTags(Ways chan element.Way) error {
	tx, _ := s.Database.Begin()
	defer tx.Commit()
	stmt, err := tx.Prepare("insert into way_tags(ref, key, value) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for way := range Ways {
		for key, value := range way.OSMElem.Tags {
			_, err := stmt.Exec(way.OSMElem.Id, key, value)
			log.Print("way tag commited")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s SQLiteDatabase) WriteRelations(Relations chan element.Relation) error {
	return errors.New("Not implemented")
}

func (s SQLiteDatabase) GetEverythingWithinCoordinates(FromLong, FromLat, ToLong, ToLat int) (*database.OSMData, error) {
	return nil, nil
}
