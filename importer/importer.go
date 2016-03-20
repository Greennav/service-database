package importer

import (
	"log"
	"sync"

	"github.com/GreenNav/service-database/database"
	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/parser/pbf"
)

const (
	CHANNELSIZE = 4 // CHANNELSIZE is used as size of the cache for nodes/ways/relations
	// during conversion from pbf file to the database

	TRANSACTIONLIMIT = 300000 //Number of elements allowed in cache

	NODES = iota //Element types
	COORDS
	WAYS
	RELATIONS
)

type elementWriter struct {
	element     interface{}
	elementType int
}

// WriteToDatabase takes a pbf file and writes it to any database
// that is conform to the database/OSMDatabase type
func WriteToDatabase(PbfFileName string, Db database.OSMDatabase) error {
	pbfCoords := make(chan []element.Node, CHANNELSIZE)
	pbfNodes := make(chan []element.Node, CHANNELSIZE)
	pbfWays := make(chan []element.Way, CHANNELSIZE)
	pbfRelations := make(chan []element.Relation, CHANNELSIZE)

	writeManagerChannel := make(chan elementWriter, CHANNELSIZE)

	pbfFile, err := pbf.Open(PbfFileName)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	writerWg := sync.WaitGroup{}

	parser := pbf.NewParser(pbfFile, pbfCoords, pbfNodes, pbfWays, pbfRelations)
	wg.Add(4)

	go func() {
		var writer elementWriter
		writer.elementType = COORDS
		for nodes := range pbfCoords {
			writer.element = nodes
			writeManagerChannel <- writer
		}
		wg.Done()
	}()
	go func() {
		var writer elementWriter
		writer.elementType = NODES
		for nodes := range pbfNodes {
			writer.element = nodes
			writeManagerChannel <- writer
		}
		wg.Done()
	}()
	go func() {
		var writer elementWriter
		writer.elementType = WAYS
		for ways := range pbfWays {
			writer.element = ways
			writeManagerChannel <- writer
		}
		wg.Done()
	}()
	go func() {
		var writer elementWriter
		writer.elementType = RELATIONS
		for relations := range pbfRelations {
			writer.element = relations
			writeManagerChannel <- writer
		}
		wg.Done()
	}()
	go func() {
		wg.Wait()
		close(writeManagerChannel)
	}()

	writerWg.Add(1)
	go writeManager(Db, writeManagerChannel, &writerWg)
	parser.Parse()
	writerWg.Wait()
	return nil
}

func writeManager(Db database.OSMDatabase, channel chan elementWriter, writerWg *sync.WaitGroup) {
	for elementCount := 0; ; elementCount = 0 {
		breakToExit := true
		wg := sync.WaitGroup{}
		err := Db.NewTransaction()
		if err != nil {
			log.Fatal(err)
		}
		nodesToWrite := make(chan element.Node, CHANNELSIZE)
		nodeTagsToWrite := make(chan element.Node, CHANNELSIZE)

		waysToWrite := make(chan element.Way, CHANNELSIZE)
		wayNodesToWrite := make(chan element.Way, CHANNELSIZE)
		wayTagsToWrite := make(chan element.Way, CHANNELSIZE)

		relationChannel := make(chan element.Relation, CHANNELSIZE)
		relationTagsChannel := make(chan element.Relation, CHANNELSIZE)
		relationMembersChannel := make(chan element.Relation, CHANNELSIZE)

		wg.Add(8)
		go writeNodeHelper(Db.WriteNodes, nodesToWrite, "Nodes", &wg)
		go writeNodeHelper(Db.WriteNodeTags, nodeTagsToWrite, "Node Tags", &wg)

		go writeWayHelper(Db.WriteWays, waysToWrite, "Ways", &wg)
		go writeWayHelper(Db.WriteWayTags, wayTagsToWrite, "Way Tags", &wg)
		go writeWayHelper(Db.WriteWayNodes, wayNodesToWrite, "Way Nodes ", &wg)

		go writeRelationHelper(Db.WriteRelation, relationChannel, "Relations", &wg)
		go writeRelationHelper(Db.WriteRelationMembers, relationMembersChannel, "Relation Member", &wg)
		go writeRelationHelper(Db.WriteRelationTags, relationTagsChannel, "Relation Tags", &wg)

		for object := range channel {
			switch object.elementType {
			case COORDS:
				nodes := object.element.([]element.Node)
				for _, node := range nodes {
					nodesToWrite <- node
				}
				elementCount += len(nodes)
			case NODES:
				nodes := object.element.([]element.Node)
				for _, node := range nodes {
					nodeTagsToWrite <- node
				}
				elementCount += len(nodes)
			case WAYS:
				ways := object.element.([]element.Way)
				for _, way := range ways {
					waysToWrite <- way
					wayTagsToWrite <- way
					wayNodesToWrite <- way
				}
				elementCount += len(ways)
			case RELATIONS:
				relations := object.element.([]element.Relation)
				for _, relation := range relations {
					relationChannel <- relation
					relationTagsChannel <- relation
					relationMembersChannel <- relation
				}
				elementCount += len(relations)
			}
			if elementCount >= TRANSACTIONLIMIT {
				log.Println("writing at ", elementCount)
				breakToExit = false
				break
			}
		}
		close(nodesToWrite)
		close(nodeTagsToWrite)
		close(wayTagsToWrite)
		close(waysToWrite)
		close(wayNodesToWrite)
		close(relationTagsChannel)
		close(relationChannel)
		close(relationMembersChannel)
		wg.Wait()
		err = Db.Commit()
		if err != nil {
			log.Println(err)
		}
		if breakToExit {
			break
		}
	}
	writerWg.Done()
}

func writeNodeHelper(function func(chan element.Node) error, channel chan element.Node, tag string, wg *sync.WaitGroup) {
	if err := function(channel); err != nil {
		log.Println(tag, ":", err)
	}
	wg.Done()
}

func writeWayHelper(function func(chan element.Way) error, channel chan element.Way, tag string, wg *sync.WaitGroup) {
	if err := function(channel); err != nil {
		log.Println(tag, ":", err)
	}
	wg.Done()
}

func writeRelationHelper(function func(chan element.Relation) error, channel chan element.Relation, tag string, wg *sync.WaitGroup) {
	if err := function(channel); err != nil {
		log.Println(tag, ":", err)
	}
	wg.Done()
}
