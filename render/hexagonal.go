/*
Copyright (c) 2022 Andre Renaud <andre@ignavus.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package render

import (
	"image"

	"github.com/disintegration/imaging"
	tiled "github.com/lafriks/go-tiled"
)

// HexagonalRendererEngine represents hexangonal rendering engine.
type HexagonalRendererEngine struct {
	m *tiled.Map
}

// Init initializes rendering engine with provided map options.
func (e *HexagonalRendererEngine) Init(m *tiled.Map) {
	e.m = m
}

// GetFinalImageSize returns final image size based on map data.
func (e *HexagonalRendererEngine) GetFinalImageSize() image.Rectangle {
	switch e.m.StaggerAxis {
	case tiled.AxisX:
		return image.Rect(0, 0, e.m.Width*e.m.TileWidth, e.m.Height*e.m.TileHeight+e.m.TileHeight/2)
	case tiled.AxisY:
		return image.Rect(0, 0, e.m.Width*e.m.TileWidth+e.m.TileWidth/2, (e.m.Height+1)*e.m.TileHeight*3/4)
	}
	return image.Rectangle{}
}

// RotateTileImage rotates provided tile layer.
func (e *HexagonalRendererEngine) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
	timg := img
	if tile.HorizontalFlip {
		timg = imaging.FlipH(timg)
	}
	if tile.VerticalFlip {
		timg = imaging.FlipV(timg)
	}
	if tile.DiagonalFlip {
		timg = imaging.FlipH(imaging.Rotate90(timg))
	}

	return timg
}

// GetTilePosition returns tile position in image.
func (e *HexagonalRendererEngine) GetTilePosition(x, y int) image.Rectangle {
	switch e.m.StaggerAxis {
	case tiled.AxisX:
		oddColumn := (x % 2) == 1
		offsetWidth := e.m.TileWidth * 3 / 4
		yBump := 0
		if oddColumn {
			yBump = e.m.TileHeight / 2
		}
		return image.Rect(x*offsetWidth,
			y*e.m.TileHeight+yBump,
			x*offsetWidth+e.m.TileWidth,
			(y+2)*e.m.TileHeight+yBump)
	case tiled.AxisY:
		oddRow := (y % 2) == 1
		offsetHeight := e.m.TileHeight * 3 / 4
		xBump := 0
		if oddRow {
			xBump = e.m.TileWidth / 2
		}
		return image.Rect(x*e.m.TileHeight+xBump,
			y*offsetHeight,
			(x+2)*e.m.TileWidth+xBump,
			(y+1)*offsetHeight+e.m.TileHeight)
	}
	return image.Rectangle{}
}
