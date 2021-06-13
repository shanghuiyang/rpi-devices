package car

import (
	"log"

	"github.com/jakefau/rpi-devices/util/geo"
	"github.com/shanghuiyang/a-star/astar"
	"github.com/shanghuiyang/a-star/tilemap"
)

const tilemapStr = `
############################################
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#      ###################                 #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
#                                          #
############################################
`

const gridsize float64 = 0.000010

var bbox = &geo.Bbox{
	Left:   116.444217,
	Right:  116.444652,
	Top:    39.956275,
	Bottom: 39.955711,
}

func findPath(org, des *geo.Point) (astar.PList, error) {
	m := tilemap.BuildFromStr(tilemapStr)

	orgXY := geo2xy(org)
	desXY := geo2xy(des)

	a := astar.New(m)
	path, err := a.FindPath(orgXY, desXY)
	if err != nil {
		log.Printf("[car]failed to find the path from A(%v) to B(%v)", org, des)
		return nil, err
	}
	log.Printf("[car]path: %v", path)
	a.Draw()
	return path, nil
}

func turnPoints(path astar.PList) astar.PList {
	if len(path) <= 2 {
		return path
	}

	var ks []float64
	for i := 0; i < len(path)-1; i++ {
		k := 99999.99
		if path[i].Y != path[i+1].Y {
			k = float64(path[i].X-path[i+1].X) / float64(path[i].Y-path[i+1].Y)
		}
		ks = append(ks, k)
	}
	log.Printf("ks: %v\n", ks)

	var turns astar.PList
	for i := 0; i < len(ks)-1; i++ {
		if ks[i] == ks[i+1] {
			continue
		}
		turns = append(turns, path[i+1])
	}
	turns = append(turns, path[len(path)-1])
	log.Printf("turn points(x,y): %v", turns)
	return turns
}

func geo2xy(p *geo.Point) *astar.Point {
	return &astar.Point{
		X: int((bbox.Top-p.Lat)/gridsize + 0.5),
		Y: int((p.Lon-bbox.Left)/gridsize + 0.5),
	}
}

func xy2geo(p *astar.Point) *geo.Point {
	return &geo.Point{
		Lat: bbox.Top - float64(p.X)*gridsize,
		Lon: bbox.Left + float64(p.Y)*gridsize,
	}
}
