package pbf_import

import (
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"net/http"
	"runtime"
  "testing"
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

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				nc++
			case *osmpbf.Way:
				// Process Way v.
				wc++
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

    if nc != 2060543 {
        t.Error("Node count not exact")
    }
    if wc != 345004 {
        t.Error("Way count not exact")
    }
    if rc != 6877 {
        t.Error("Relation count not exact")
    }
}