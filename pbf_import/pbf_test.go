package pbf_import

import (
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"net/http"
	"runtime"
  "testing"
)
const (
	EXACT_NODE_COUNT = 2060543
	EXACT_WAY_COUNT = 345004
	EXACT_RELATION_COUNT = 6877
)
func TestImportFromHttp(t *testing.T) {
	// f, err := os.Open("greater-london-140324.osm.pbf")
	response, err := http.Get("http://download.bbbike.org/osm/bbbike/Luebeck/Luebeck.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}
	f := response.Body
	// defer f.Close()

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1)) // use several goroutines for faster decoding
	if err != nil {
		log.Fatal(err)
	}

	var nodecount, waycount, relationcount uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				nodecount++
			case *osmpbf.Way:
				// Process Way v.
				waycount++
			case *osmpbf.Relation:
				// Process Relation v.
				relationcount++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

    if nodecount != EXACT_NODE_COUNT {
        t.Error("Node count not exact")
    }
    if waycount != EXACT_WAY_COUNT {
        t.Error("Way count not exact")
    }
    if relationcount != EXACT_RELATION_COUNT {
        t.Error("Relation count not exact")
    }
}