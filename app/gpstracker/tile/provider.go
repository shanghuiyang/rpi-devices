package tile

import (
	sm "github.com/shanghuiyang/go-staticmaps"
)

const (
	OsmTile             = "osm"
	BingSatelliteTile   = "bing-satellite"
	GoogleSatelliteTile = "google-satellite"
)

func NewLocalTileProvider(name string) *sm.TileProvider {
	return &sm.TileProvider{
		Name:           name,
		Attribution:    "",
		IgnoreNotFound: true,
		TileSize:       256,
		URLPattern:     "",
		Shards:         []string{},
	}
}
