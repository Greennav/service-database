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

func (s *SQLiteDatabase) readNodeById(Id int64) (node element.Node, err error) {
	node.Id = Id
	row := s.Database.QueryRow("select lon,lat from nodes where id=?", Id)
	err = row.Scan(&node.Long, &node.Lat)
	if err != nil {
		return
	}
	return
}

func (s *SQLiteDatabase) readNodeTagById(Id int64) (tagMap element.Tags, err error) {
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

func (s *SQLiteDatabase) readWayById(Id int64) (way element.Way, err error) {
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

func (s *SQLiteDatabase) readWayTagsById(Id int64) (tagMap element.Tags, err error) {
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

func (s *SQLiteDatabase) readWayNodesById(Id int64) (nodes []int64, err error) {
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

func (s *SQLiteDatabase) readRelationById(Id int64) (relation element.Relation, err error) {
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

func (s *SQLiteDatabase) readRelationTagsById(Id int64) (tagMap element.Tags, err error) {
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

func (s *SQLiteDatabase) readRelationMembersById(Id int64) (members []element.Member, err error) {
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

func (s *SQLiteDatabase) ReadNode(Id int64) (node element.Node, err error) {
	node, err = s.readNodeById(Id)
	if err != nil {
		return
	}
	tags, err := s.readNodeTagById(Id)
	if err != nil {
		return
	}
	node.Tags = tags
	return
}

func (s *SQLiteDatabase) ReadWay(Id int64) (way element.Way, err error) {
	way, err = s.readWayById(Id)
	if err != nil {
		return
	}
	tags, err := s.readWayTagsById(Id)
	if err != nil {
		return
	}
	way.Tags = tags
	nodes, err := s.readWayNodesById(Id)
	if err != nil {
		return
	}
	way.Refs = nodes
	return
}

func (s *SQLiteDatabase) ReadRelation(Id int64) (relation element.Relation, err error) {
	relation, err = s.readRelationById(Id)
	if err != nil {
		return
	}
	tags, err := s.readRelationTagsById(Id)
	if err != nil {
		return
	}
	relation.Tags = tags
	members, err := s.readRelationMembersById(Id)
	if err != nil {
		return
	}
	relation.Members = members
	return
}

//TO be verified
func (s *SQLiteDatabase) ReadNodesByCoordinates(FromLat, FromLon, ToLat, ToLon float64) (nodes []element.Node, err error) {
	var (
		id          int64
		key, value  string
		node        element.Node
		nodesTagMap = make(map[int64]element.Tags)
	)
	nodeTagRows, err := s.Database.Query("select t.ref,t.key,t.value from node_tags t join nodes n on n.id=t.ref where n.lon between ? and ? and n.lat between ? and ?",
		FromLon, ToLon, FromLat, ToLat)
	if err != nil {
		return
	}
	if err = nodeTagRows.Err(); err != nil {
		return
	}
	defer nodeTagRows.Close()
	for nodeTagRows.Next() {
		err = nodeTagRows.Scan(&id, &key, &value)
		if err != nil {
			return
		}
		if _, available := nodesTagMap[id]; !available {
			nodesTagMap[id] = make(element.Tags)
		}
		nodesTagMap[id][key] = value
	}

	nodeRows, err := s.Database.Query("select id, lon,lat from nodes where lon between ? and ? and lat between ? and ?", FromLon, ToLon, FromLat, ToLat)
	if err != nil {
		return
	}
	if err = nodeRows.Err(); err != nil {
		return
	}
	defer nodeRows.Close()
	for nodeRows.Next() {
		err = nodeRows.Scan(&node.Id, &node.Long, &node.Lat)
		if err != nil {
			return
		}
		if _, available := nodesTagMap[node.Id]; available {
			node.Tags = nodesTagMap[node.Id]
		} else {
			node.Tags = nil
		}
		nodes = append(nodes, node)
	}
	return
}

//To be verified
func (s *SQLiteDatabase) ReadWaysByCoordinates(FromLat, FromLon, ToLat, ToLon float64) (ways []element.Way, err error) {
	var (
		id, nodeId int64
		key, value string
		waysMap    = make(map[int64]*element.Way)
	)
	wayTagRows, err := s.Database.Query("select t.ref,t.key,t.value from way_tags t join way_nodes wn join nodes n on(n.id=wn.node and wn.way=t.ref) where n.lat between ? and ? and n.lon between ? and ?",
		FromLat, ToLat, FromLon, ToLon)
	if err != nil {
		return
	}
	if err = wayTagRows.Err(); err != nil {
		return
	}
	defer wayTagRows.Close()
	for wayTagRows.Next() {
		err = wayTagRows.Scan(&id, &key, &value)
		if err != nil {
			return
		}
		if _, available := waysMap[id]; !available {
			waysMap[id] = new(element.Way)
			waysMap[id].Tags = make(element.Tags)
			waysMap[id].Id = id
		}
		waysMap[id].Tags[key] = value
	}

	wayNodeRows, err := s.Database.Query("select w.way,w.node from way_nodes w join nodes n on (w.node=n.id) where n.lat between ? and ? and n.lon between ? and ? order by w.way,w.num",
		FromLat, ToLat, FromLon, ToLon)
	if err != nil {
		return
	}
	if err = wayNodeRows.Err(); err != nil {
		return
	}
	defer wayNodeRows.Close()
	for wayNodeRows.Next() {
		err = wayNodeRows.Scan(&id, &nodeId)
		if err != nil {
			return
		}
		if _, available := waysMap[id]; !available {
			waysMap[id] = new(element.Way)
			waysMap[id].Id = id
		}
		waysMap[id].Refs = append(waysMap[id].Refs, nodeId)
	}

	for _, value := range waysMap {
		ways = append(ways, *value)
	}
	return
}

//Help needed wrt the implementation
func (s *SQLiteDatabase) ReadRelationsByCoordinates() (relations []element.Member, err error) {
	return
}

func (s *SQLiteDatabase) ReadEverythingWithinCoordinates(FromLat, FromLon, ToLat, ToLon float64) (osmData *database.OSMData, err error) {
	if FromLat > ToLat {
		FromLat, ToLat = ToLat, FromLat
	}
	if FromLon > ToLon {
		FromLon, ToLon = ToLon, FromLon
	}
	nodes, err := s.ReadNodesByCoordinates(FromLat, FromLon, ToLat, ToLon)
	if err != nil {
		return
	}
	ways, err := s.ReadWaysByCoordinates(FromLat, FromLon, ToLat, ToLon)
	if err != nil {
		return
	}
	osmData = new(database.OSMData)
	osmData.Nodes = nodes
	osmData.Ways = ways
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
