package pbf_import

import (
	"github.com/qedus/osmpbf"
	"io"
	"os"
	"log"
	_ "net/http"
	"runtime"
	"fmt"
	"github.com/sitaramshelke/service-database/dbhelper"
)

func ImportFromHttp() {
	// to save time read from file but works for both 
	// response, err := http.Get("http://download.bbbike.org/osm/bbbike/Luebeck/Luebeck.osm.pbf")
	f, err := os.Open("pbf_import/Luebeck.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}
	// f := response.Body

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(-1)) // use several goroutines for faster decoding
	if err != nil {
		log.Fatal(err)
	}

	var nc, wc, rc uint64
	count := 100 //for shorter test time
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				dbhelper.InsertNode(v)
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
		count -= 1
		if count <= 0 {
			break
		}
	}
	fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)


}