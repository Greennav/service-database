package importer

import (
	"github.com/GreenNav/service-database/database"
	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/parser/pbf"
	"log"
	"sync"
)

const (
	CHANNELSIZE = 10
)

func WriteToDatabase(pbfFileName string, db database.OSMDatabase) {
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

	pbfFile, err := pbf.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	parser := pbf.NewParser(pbfFile, nil, pbfNodes, pbfWays, pbfRelations)
	wg.Add(8)
	go func() {
		err := db.WriteNodes(nodesToWrite)
		if err != nil {
			log.Fatal(err.Error() + "(table nodes)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteNodeTags(nodeTagsToWrite)
		if err != nil {
			log.Fatal(err.Error() + "(table node_tags)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteWays(waysToWrite)
		if err != nil {
			log.Fatal(err.Error() + "(table ways)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteWayTags(wayTagsToWrite)
		if err != nil {
			log.Fatal(err.Error() + "(table way_tags)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteWayNodes(wayNodesToWrite)
		if err != nil {
			log.Fatal(err.Error() + "(table way_nodes)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteRelation(relationChannel)
		if err != nil {
			log.Println(err.Error(), "(table relation_nodes)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteRelationTags(relationTagsChannel)
		if err != nil {
			log.Println(err.Error(), "(table relation_tags)")
		}
		wg.Done()
	}()
	go func() {
		err := db.WriteRelationMembers(relationMembersChannel)
		if err != nil {
			log.Println(err.Error(), "(table relation_members)")
		}
		wg.Done()
	}()
	go func() {
		for nodes := range pbfNodes {
			for _, node := range nodes {
				nodesToWrite <- node
				nodeTagsToWrite <- node
			}
		}
		close(nodesToWrite)
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
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func WriteNodesToDatabase(pbfFileName string, db database.OSMDatabase) {
	pbfNodes := make(chan []element.Node)

	nodesToWrite := make(chan element.Node)
	pbfFile, err := pbf.Open(pbfFileName)
	wg := sync.WaitGroup{}
	if err != nil {
		log.Fatal(err)
	}
	pbfParser := pbf.NewParser(pbfFile, nil, pbfNodes, nil, nil)
	wg.Add(1)
	go func() {
		err := db.WriteNodes(nodesToWrite)
		if err != nil {
			log.Fatal(err.Error() + "(table nodes)")
		}
		wg.Done()
	}()
	go func() {
		for nodes := range pbfNodes {
			for _, node := range nodes {
				nodesToWrite <- node
			}
		}
		close(nodesToWrite)
	}()
	pbfParser.Parse()
	wg.Wait()
}
