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
	FileName    string
	Database    *sql.DB
	Transaction *sql.Tx
}

func CreateEmpty(Name string, SchemaFile string) (*SQLiteDatabase, error) {
	schema, err := ioutil.ReadFile(SchemaFile)
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
	return &SQLiteDatabase{FileName: Name, Database: db}, nil
}

func GetByName(Name string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", "file:"+Name+"?cache=shared&mode=rwc")
	return &SQLiteDatabase{FileName: Name, Database: db}, err
}

func (s *SQLiteDatabase) NewTransaction() error {
	if s.Transaction != nil {
		return errors.New("Uncommited transaction still open")
	}
	tx, err := s.Database.Begin()
	if err != nil {
		return err
	}
	s.Transaction = tx
	return nil
}

func (s *SQLiteDatabase) WriteNodes(Nodes chan element.Node) error {
	stmt, err := s.Transaction.Prepare("insert into nodes(id, lon, lat) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for node := range Nodes {
		_, err = stmt.Exec(node.OSMElem.Id, node.Long, node.Lat)
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
		for index, nodeId := range way.Refs {
			_, err = stmt.Exec(way.Id, index, nodeId)
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
	}
	return nil
}

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
		}
	}
	return nil
}

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
		}
	}
	return nil
}

func (s *SQLiteDatabase) ReadNode(Id int64) (node element.Node, err error) {
	node.Id = Id
	row := s.Database.QueryRow("select lon,lat from nodes where id=?", Id)
	err = row.Scan(&node.Long, &node.Lat)
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadNodeTag(Id int64) (tagMap element.Tags, err error) {
	rows, err := s.Database.Query("select key,value from node_tags where ref=?", Id)
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		key   string
		value string
	)
	tagMap = make(element.Tags)
	for rows.Next() {
		rows.Scan(&key, &value)
		tagMap[key] = value
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadWay(Id int64) (way element.Way, err error) {
	row := s.Database.QueryRow("select id from ways where id=?", Id)
	if err != nil {
		return
	}
	err = row.Scan(&way.Id)
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadWayTags(Id int64) (tagMap element.Tags, err error) {
	rows, err := s.Database.Query("select key,value from way_tags where ref=?", Id)
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		key   string
		value string
	)
	tagMap = make(element.Tags)
	for rows.Next() {
		rows.Scan(&key, &value)
		tagMap[key] = value
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadWayNodes(Id int64) (nodes []int64, err error) {
	rows, err := s.Database.Query("select node from way_nodes where way=? order by num", Id)
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		node int64
	)
	for rows.Next() {
		rows.Scan(&node)
		nodes = append(nodes, node)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadRelation(Id int64) (relation element.Relation, err error) {
	row := s.Database.QueryRow("select id from relations where id=?", Id)
	if err != nil {
		return
	}
	err = row.Scan(&relation.Id)
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadRelationTags(Id int64) (tagMap element.Tags, err error) {
	rows, err := s.Database.Query("select key,value from relation_tags where ref=?", Id)
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		key   string
		value string
	)
	tagMap = make(element.Tags)
	for rows.Next() {
		rows.Scan(&key, &value)
		tagMap[key] = value
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) ReadRelationMembers(Id int64) (members []element.Member, err error) {
	rows, err := s.Database.Query("select type,ref,role from members where relation=?", Id)
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		typeStr int64
		ref     int64
		role    string
	)
	var member element.Member
	for rows.Next() {
		rows.Scan(&typeStr, &ref, &role)
		member.Id = ref
		member.Role = role
		member.Type = element.MemberType(typeStr)
		members = append(members, member)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) Commit() error {
	if s.Transaction == nil {
		return errors.New("Empty transaction")
	}
	err := s.Transaction.Commit()
	if err != nil {
		return err
	}
	s.Transaction = nil
	return nil
}

func (s *SQLiteDatabase) Close() error {
	err := s.Database.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteDatabase) GetEverythingWithinCoordinates(FromLong, FromLat, ToLong, ToLat int) (*database.OSMData, error) {
	return nil, nil
}
