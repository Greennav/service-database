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

	GetEverythingWithinCoordinates(FromLong, FromLat, ToLong, ToLat int) (*OSMData, error)

	NewTransaction() error
	Commit() error
	Close() error
}
