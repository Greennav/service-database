package sqlite

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
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
	db, err := CreateTest()
	if err != nil {
		t.Error(err)
	}
	InsertTest(db, elements)
	err = readTest(db, elements)
	if err != nil {
		t.Error(err)
	}
	db.Close()
	os.Remove(TESTDB)

}

func CreateTest() (*SQLiteDatabase, error) {
	db, err := CreateEmpty(TESTDB, SCHEMAFILE)
	if err != nil {
		return nil, err
	}
	return db, err
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

func readTest(db *SQLiteDatabase, elements TestBundle) (err error) {
	if err = checkNode(db, elements); err != nil {
		return
	}
	if err = checkWay(db, elements); err != nil {
		return
	}
	if err = checkRelation(db, elements); err != nil {
		return
	}
	return
}

func checkNode(db *SQLiteDatabase, elements TestBundle) (err error) {
	node, err := db.ReadNode(elements.Node.Id)
	if err != nil {
		return
	}
	if !reflect.DeepEqual(node, elements.Node) {
		err = errors.New("Node:Corrupt")
	}
	return
}

func checkWay(db *SQLiteDatabase, elements TestBundle) (err error) {
	way, err := db.ReadWay(elements.Way.Id)
	if err != nil {
		return
	}
	if !reflect.DeepEqual(way, elements.Way) {
		err = errors.New("Way:Corrupt")
		return
	}
	return
}

func checkRelation(db *SQLiteDatabase, elements TestBundle) (err error) {
	relation, err := db.ReadRelation(elements.Relation.Id)
	if err != nil {
		return
	}
	if !reflect.DeepEqual(relation, elements.Relation) {
		err = errors.New("Relation:Corrupt")
		return
	}
	return
}
