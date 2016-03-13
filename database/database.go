package database

import (
	"github.com/omniscale/imposm3/element"
)

type OSMData struct {
	Nodes     []element.Node
	Ways      []element.Way
	Relations []element.Relation
}

type OSMDatabase interface {
	WriteNodes(Nodes chan []element.Node) error
	WriteWays(Ways chan []element.Way) error
	WriteRelations(Relations chan []element.Relation) error

	GetEverythingWithinCoordinates(FromLong, FromLat, ToLong, ToLat int) (*OSMData, error)
}