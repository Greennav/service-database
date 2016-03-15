package sqlite

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/GreenNav/service-database/database"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/omniscale/imposm3/element"
)

type SQLiteDatabase struct {
	FileName    string
	Database    *sql.DB
	Transaction *sql.Tx
}

func CreateEmpty(Name string) (*SQLiteDatabase, error) {
	schema, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", "file:"+Name+"?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return &SQLiteDatabase{FileName: Name, Database: db, Transaction: tx}, err
}

func GetByName(Name string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", Name)
	return &SQLiteDatabase{FileName: Name, Database: db}, err
}

func (s *SQLiteDatabase) WriteNodes(Nodes chan element.Node) error {
	stmt, err := s.Transaction.Prepare("insert into nodes(id, lon, lat) values(?, ?, ?)")
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

func (s *SQLiteDatabase) WriteNodeTags(Nodes chan element.Node) error {
	stmt, err := s.Transaction.Prepare("insert into node_tags(ref, key, value) values(?, ?, ?)")
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

func (s *SQLiteDatabase) WriteWays(Ways chan element.Way) error {
	stmt, err := s.Transaction.Prepare("insert into ways(id) values(?)")
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

func (s *SQLiteDatabase) WriteWayNodes(Ways chan element.Way) error {
	stmt, err := s.Transaction.Prepare("insert into way_nodes(way, num, node) values(?, ?, ?)")
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

func (s *SQLiteDatabase) WriteWayTags(Ways chan element.Way) error {
	stmt, err := s.Transaction.Prepare("insert into way_tags(ref, key, value) values(?, ?, ?)")
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

func (s *SQLiteDatabase) WriteRelation(Relations chan element.Relation) error {
	stmt, err := s.Transaction.Prepare("insert into relations(id) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for relation := range Relations {
		_, err := stmt.Exec(relation.Id)
		if err != nil {
			return err
		}
		log.Println("relation committed")
	}
	return nil
}

//To be Verified
func (s *SQLiteDatabase) WriteRelationTags(Relations chan element.Relation) error {
	stmt, err := s.Transaction.Prepare("insert into relation_tags(ref,key,value) values(?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for relation := range Relations {
		for key, value := range relation.OSMElem.Tags {
			_, err := stmt.Exec(relation.OSMElem.Id, key, value)
			if err != nil {
				return err
			}
			log.Println("relation tag committed")
		}
	}
	return nil
}

//To be Verified
func (s *SQLiteDatabase) WriteRelationMembers(Relations chan element.Relation) error {
	stmt, err := s.Transaction.Prepare("insert into members(relation,type,ref,role) values(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for relation := range Relations {
		for _, member := range relation.Members {
			_, err := stmt.Exec(relation.Id, member.Type, member.Id, member.Role)
			if err != nil {
				return err
			}
			log.Println("member committed")
		}
	}
	return nil
}

func (s *SQLiteDatabase) Close() error {
	err := s.Transaction.Commit()
	if err != nil {
		return err
	}
	err = s.Database.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteDatabase) GetEverythingWithinCoordinates(FromLong, FromLat, ToLong, ToLat int) (*database.OSMData, error) {
	return nil, nil
}
