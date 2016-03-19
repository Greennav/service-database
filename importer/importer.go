package importer

import (
	"log"
	"sync"

	"github.com/GreenNav/service-database/database"
	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/parser/pbf"
)

const (
	// CHANNELSIZE is used as size of the cache for nodes/ways/relations
	// during conversion from pbf file to the database
	CHANNELSIZE = 4
)

// WriteToDatabase takes a pbf file and writes it to any database
// that is conform to the database/OSMDatabase type
func WriteToDatabase(PbfFileName string, Db database.OSMDatabase) error {
	err := Db.NewTransaction()
	if err != nil {
		log.Fatal(err)
	}
	pbfCoords := make(chan []element.Node, CHANNELSIZE)
	pbfNodes := make(chan []element.Node, CHANNELSIZE)
	pbfWays := make(chan []element.Way, CHANNELSIZE)
	pbfRelations := make(chan []element.Relation, CHANNELSIZE)

	nodesToWrite := make(chan element.Node, CHANNELSIZE)
	nodeTagsToWrite := make(chan element.Node, CHANNELSIZE)

	waysToWrite := make(chan element.Way, CHANNELSIZE)
	wayNodesToWrite := make(chan element.Way, CHANNELSIZE)
	wayTagsToWrite := make(chan element.Way, CHANNELSIZE)

	relationChannel := make(chan element.Relation, CHANNELSIZE)
	relationTagsChannel := make(chan element.Relation, CHANNELSIZE)
	relationMembersChannel := make(chan element.Relation, CHANNELSIZE)

	pbfFile, err := pbf.Open(PbfFileName)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	parser := pbf.NewParser(pbfFile, pbfCoords, pbfNodes, pbfWays, pbfRelations)

	wg.Add(8)

	go writeNodeHelper(Db.WriteNodes, nodesToWrite, "Nodes", &wg)
	go writeNodeHelper(Db.WriteNodeTags, nodeTagsToWrite, "Node Tags", &wg)

	go writeWayHelper(Db.WriteWays, waysToWrite, "Ways", &wg)
	go writeWayHelper(Db.WriteWayTags, wayTagsToWrite, "Way Tags", &wg)
	go writeWayHelper(Db.WriteWayNodes, wayNodesToWrite, "Way Nodes ", &wg)

	go writeRelationHelper(Db.WriteRelation, relationChannel, "Relations", &wg)
	go writeRelationHelper(Db.WriteRelationMembers, relationMembersChannel, "Relation Member", &wg)
	go writeRelationHelper(Db.WriteRelationTags, relationTagsChannel, "Relation Tags", &wg)

	go func() {
		for nodes := range pbfCoords {
			for _, node := range nodes {
				nodesToWrite <- node
			}
		}
		close(nodesToWrite)
	}()
	go func() {
		for nodes := range pbfNodes {
			for _, node := range nodes {
				nodeTagsToWrite <- node
			}
		}
		close(nodeTagsToWrite)
	}()
	go func() {
		for ways := range pbfWays {
			for _, way := range ways {
				waysToWrite <- way
				wayNodesToWrite <- way
				wayTagsToWrite <- way
			}
		}
		close(waysToWrite)
		close(wayNodesToWrite)
		close(wayTagsToWrite)
	}()
	go func() {
		for relations := range pbfRelations {
			for _, relation := range relations {
				relationChannel <- relation
				relationMembersChannel <- relation
				relationTagsChannel <- relation
			}
		}
		close(relationChannel)
		close(relationMembersChannel)
		close(relationTagsChannel)
	}()

	parser.Parse()
	wg.Wait()
	if err := Db.Commit(); err != nil {
		return err
	}

	return nil
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
