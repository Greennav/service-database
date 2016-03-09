package main
import(
  "github.com/sitaramshelke/service-database/dbhelper"
  _ "fmt"
  "github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

func main(){

  api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/node", dbhelper.GetAllNodeData),
		rest.Get("/country",dbhelper.GetCountryName),
		// rest.Post("/device",AddDevice),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8888", api.MakeHandler()))
  
}
