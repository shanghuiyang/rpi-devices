package tile

import (
	"bytes"
	"image"

	"image/jpeg"

	"github.com/golang/geo/s2"
	sm "github.com/shanghuiyang/go-staticmaps"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

type Render struct {
	ctx *sm.Context
}

func NewRender() *Render {
	return &Render{
		ctx: sm.NewContext(),
	}
}

func (m *Render) SetSize(width, height int) {
	m.ctx.SetSize(width, height)
}

func (m *Render) SetZoom(zoom int) {
	m.ctx.SetZoom(zoom)
}

func (m *Render) SetCenter(pt *geo.Point) {
	m.ctx.SetCenter(s2.LatLngFromDegrees(pt.Lat, pt.Lon))
}

func (m *Render) SetCache(cache sm.TileCache) {
	m.ctx.SetCache(cache)
}

func (m *Render) SetTileProvider(tileProvider *sm.TileProvider) {
	m.ctx.SetTileProvider(tileProvider)
}

// func (m *Render) SetTileFetcher(tileFetecher *sm.TileFetcher) {
// 	m.ctx.SetTileFetcher(tileFetecher)
// }

func (m *Render) SetOnline(online bool) {
	m.ctx.SetOnline(online)
}

func (m *Render) AddMarker(marker *sm.Marker) {
	m.ctx.AddObject(marker)
}

func (m *Render) ClearMarker() {
	m.ctx.ClearObjects()
}

func (m *Render) Render() ([]byte, error) {
	img, err := m.ctx.Render()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Render) RenderImg() (*image.Image, error) {
	img, err := m.ctx.Render()
	if err != nil {
		return nil, err
	}
	return &img, nil
}
