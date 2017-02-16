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
	WriteNodes(Nodes chan element.Node) error
	WriteNodeTags(Nodes chan element.Node) error

	WriteWays(Ways chan element.Way) error
	WriteWayNodes(Ways chan element.Way) error
	WriteWayTags(Ways chan element.Way) error

	WriteRelation(Relations chan element.Relation) error
	WriteRelationTags(Relations chan element.Relation) error
	WriteRelationMembers(Relations chan element.Relation) error

	ReadNode(Id int64) (node element.Node, err error)
	ReadWay(Id int64) (way element.Way, err error)
	ReadRelation(id int64) (relation element.Relation, err error)

	ReadNodesByCoordinates(FromLat, FromLon, ToLat, ToLon float64) (nodes []element.Node, err error)
	ReadWaysByCoordinates(FromLat, FromLon, ToLat, ToLon float64) (ways []element.Way, err error)
	ReadRelationsByCoordinates() (relations []element.Member, err error)

	ReadEverythingWithinCoordinates(FromLat, FromLon, ToLat, ToLon float64) (*OSMData, error)

	NewTransaction() error
	Commit() error
	Close() error
}
