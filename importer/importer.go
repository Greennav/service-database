package importer

import (
	"log"
	"sync"

	"github.com/GreenNav/service-database/database"
	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/parser/pbf"
)

func WriteToDatabase(pbfFileName string, destination database.OSMDatabase) {
	nodes := make(chan []element.Node, 1000000)
	ways := make(chan []element.Way, 1000000)
	relations := make(chan []element.Relation, 1000000)

	pbfFile, err := pbf.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}

	parser := pbf.NewParser(pbfFile, nil, nodes, ways, relations)
	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		destination.WriteNodes(nodes)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		destination.WriteWays(ways)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		destination.WriteRelations(relations)
		wg.Done()
	}()
  
  parser.Parse()
	wg.Wait()
}