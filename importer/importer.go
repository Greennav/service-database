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
	go func() {
		err := Db.WriteNodes(nodesToWrite)
		if err != nil {
			log.Fatal(err.Error() + " (table nodes)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteNodeTags(nodeTagsToWrite)
		if err != nil {
			log.Fatal(err.Error() + " (table node_tags)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteWays(waysToWrite)
		if err != nil {
			log.Fatal(err.Error() + " (table ways)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteWayTags(wayTagsToWrite)
		if err != nil {
			log.Fatal(err.Error() + " (table way_tags)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteWayNodes(wayNodesToWrite)
		if err != nil {
			log.Fatal(err.Error() + " (table way_nodes)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteRelation(relationChannel)
		if err != nil {
			log.Fatal(err.Error(), " (table relation_nodes)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteRelationTags(relationTagsChannel)
		if err != nil {
			log.Fatal(err.Error(), " (table relation_tags)")
		}
		wg.Done()
	}()
	go func() {
		err := Db.WriteRelationMembers(relationMembersChannel)
		if err != nil {
			log.Fatal(err.Error(), " (table relation_members)")
		}
		wg.Done()
	}()
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
	if err := Db.Close(); err != nil {
		return err
	}

	return nil
}
