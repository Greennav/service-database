package sqlite

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/omniscale/imposm3/element"
)

const (
	JSONFILE   = "sqlite_test.json"
	TESTDB     = "test.db"
	SCHEMAFILE = "../../schemata/sqlite.sql"
)

type TestBundle struct {
	Node     element.Node
	Way      element.Way
	Relation element.Relation
}

func TestSqlite(t *testing.T) {
	data, err := ioutil.ReadFile(JSONFILE)
	if err != nil {
		t.Error(err)
	}
	elements := TestBundle{}
	err = json.Unmarshal(data, &elements)
	if err != nil {
		t.Error(err)
	}
	elements.Node.Id = 100
	elements.Way.Id = 200
	elements.Relation.Id = 100
	db := CreateTest()
	InsertTest(db, elements)
	readTest(db, elements)
	db.Close()
	os.Remove(TESTDB)

}

func CreateTest() *SQLiteDatabase {
	db, err := CreateEmpty(TESTDB, SCHEMAFILE)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InsertTest(db *SQLiteDatabase, elements TestBundle) {
	err := db.NewTransaction()
	if err != nil {
		log.Fatal(err)
	}
	nodeChan := make(chan element.Node)
	nodeTagChan := make(chan element.Node)

	wayChan := make(chan element.Way)
	wayNodeChan := make(chan element.Way)
	wayTagChan := make(chan element.Way)

	relationChan := make(chan element.Relation)
	relationTagChan := make(chan element.Relation)
	relationMemberChan := make(chan element.Relation)

	wg := sync.WaitGroup{}
	wg.Add(8)

	go func() {
		err := db.WriteNodes(nodeChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteNodeTags(nodeTagChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	go func() {
		err := db.WriteWays(wayChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteWayNodes(wayNodeChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteWayTags(wayTagChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	go func() {
		err := db.WriteRelation(relationChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteRelationTags(relationTagChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteRelationMembers(relationMemberChan)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	nodeChan <- elements.Node
	nodeTagChan <- elements.Node

	wayChan <- elements.Way
	wayNodeChan <- elements.Way
	wayTagChan <- elements.Way

	relationChan <- elements.Relation
	relationMemberChan <- elements.Relation
	relationTagChan <- elements.Relation

	close(nodeChan)
	close(nodeTagChan)

	close(wayTagChan)
	close(wayNodeChan)
	close(wayChan)

	close(relationTagChan)
	close(relationMemberChan)
	close(relationChan)

	wg.Wait()
	db.Commit()
}

func readTest(db *SQLiteDatabase, elements TestBundle) {
	checkNode(db, elements)
	checkWay(db, elements)
	checkRelation(db, elements)
}

func checkNode(db *SQLiteDatabase, elements TestBundle) {
	node, err := db.ReadNode(elements.Node.Id)
	if err != nil {
		log.Fatal(err)
	}
	nodeTags, err := db.ReadNodeTag(elements.Node.Id)
	if err != nil {
		log.Fatal(err)
	}
	if node.Lat != elements.Node.Lat && node.Long != elements.Node.Long {
		log.Fatal("Node:Mismatched Latitude and longitude")
	}
	for key, value := range nodeTags {
		if elements.Node.Tags[key] != value {
			log.Fatal("Node:Mismatched Tags")
		}
	}
}

func checkWay(db *SQLiteDatabase, elements TestBundle) {
	_, err := db.ReadWay(elements.Way.Id)
	if err != nil {
		log.Fatal(err)
	}
	wayTags, err := db.ReadWayTags(elements.Way.Id)
	if err != nil {
		log.Fatal(err)
	}
	wayNodes, err := db.ReadWayNodes(elements.Way.Id)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range wayTags {
		if elements.Way.Tags[key] != value {
			log.Fatal("Way:Mismatched Tags")
		}
	}
	if len(wayNodes) != len(elements.Way.Refs) {
		log.Fatal("Way:Mismatched Node length")
	}
	for i, num := range wayNodes {
		if elements.Way.Refs[i] != num {
			log.Fatal("Way:Mismatched Node")
		}
	}
}

func checkRelation(db *SQLiteDatabase, elements TestBundle) {
	_, err := db.ReadRelation(elements.Relation.Id)
	if err != nil {
		log.Fatal(err)
	}
	relationTags, err := db.ReadRelationTags(elements.Relation.Id)
	if err != nil {
		log.Fatal(err)
	}
	Members, err := db.ReadRelationMembers(elements.Relation.Id)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range relationTags {
		if elements.Relation.Tags[key] != value {
			log.Fatal("Relation:Mismatched Tags")
		}
	}
	if len(Members) != len(elements.Relation.Members) {
		log.Fatal("Relation:Mismatched Member length")
	}
	for i, member := range Members {
		orginalMember := elements.Relation.Members[i]
		if orginalMember.Id != member.Id && orginalMember.Role != member.Role && orginalMember.Type != member.Type {
			log.Fatal("Relation:Mismatched Member")
		}
	}

}
