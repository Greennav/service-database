package importer

import (
	"log"
	"sync"

	"github.com/GreenNav/service-database/database"
	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/parser/pbf"
)

func WriteToDatabase(pbfFileName string, db database.OSMDatabase) {
	pbfNodes := make(chan []element.Node)
	pbfWays := make(chan []element.Way)
  
  nodesToWrite := make(chan element.Node)
  nodeTagsToWrite := make(chan element.Node)
  
	waysToWrite := make(chan element.Way)
	wayNodesToWrite := make(chan element.Way)
	wayTagsToWrite := make(chan element.Way)

	pbfFile, err := pbf.Open(pbfFileName)
	if err != nil {
		log.Fatal(err)
	}
  
  
  wg := sync.WaitGroup{}
	parser := pbf.NewParser(pbfFile, nil, pbfNodes, pbfWays, nil)
  
  wg.Add(5)
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
  
  parser.Parse()
  wg.Wait()
}