package render

import (
	"github.com/disintegration/imaging"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/utils"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// RenderVisibleGroups renders all visible groups
func (r *Renderer) RenderVisibleGroups() error {
	for _, group := range r.m.Groups {
		if !group.Visible {
			continue
		}
		r._renderGroup(group)
	}
	return nil
}

// RenderGroup renders single group.
func (r *Renderer) RenderGroup(groupIdx int) error {
	group := r.m.Groups[groupIdx]
	return r._renderGroup(group)
}

func (r *Renderer) _renderGroup(group *tiled.Group) error {
	for _, layer := range group.Layers {
		if !layer.Visible {
			continue
		}
		if err := r._renderLayer(layer); err != nil {
			return err
		}
	}

	for _, objectGroup := range group.ObjectGroups {
		if !objectGroup.Visible {
			continue
		}
		err := r._renderObjectGroup(objectGroup)
		if err != nil {
			return err
		}
	}

	return nil
}

// RenderVisibleLayersAndObjectGroups render all layers and object groups, layer first, objectGroup second
// so the order may be incorrect,
// you may put them into different groups, then call RenderVisibleGroups
func (r *Renderer) RenderVisibleLayersAndObjectGroups() error {
	// TODO: The order maybe incorrect

	if err := r.RenderVisibleLayers(); err != nil {
		return err
	}
	return r.RenderVisibleObjectGroups()
}

// RenderVisibleObjectGroups renders all visible object groups
func (r *Renderer) RenderVisibleObjectGroups() error {
	for i, layer := range r.m.ObjectGroups {
		if !layer.Visible {
			continue
		}
		if err := r.RenderObjectGroup(i); err != nil {
			return err
		}
	}
	return nil
}

// RenderObjectGroup renders a single object group
func (r *Renderer) RenderObjectGroup(i int) error {
	layer := r.m.ObjectGroups[i]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) _renderObjectGroup(objectGroup *tiled.ObjectGroup) error {
	objs := objectGroup.Objects
	objs = utils.SortAny(objs, utils.SortObjectsLess)
	for _, obj := range objs {
		if err := r.renderOneObject(objectGroup, obj); err != nil {
			return err
		}
	}
	return nil
}

// RenderGroupObjectGroup renders single object group in a certain group.
func (r *Renderer) RenderGroupObjectGroup(groupIdx, objectGroupId int) error {
	group := r.m.Groups[groupIdx]
	layer := group.ObjectGroups[objectGroupId]
	return r._renderObjectGroup(layer)
}

func (r *Renderer) renderOneObject(layer *tiled.ObjectGroup, o *tiled.Object) error {
	if !o.Visible {
		return nil
	}

	if o.GID == 0 {
		// TODO: o.GID == 0
		return nil
	}

	tile, err := r.m.TileGIDToTile(o.GID)
	if err != nil {
		return err
	}

	img, err := r.getTileImage(tile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	srcSize := bounds.Size()
	dstSize := image.Pt(int(o.Width), int(o.Height))

	if !srcSize.Eq(dstSize) {
		img = imaging.Resize(img, dstSize.X, dstSize.Y, imaging.NearestNeighbor)
	}

	var originPoint image.Point

	img, originPoint = r._rotateObjectImage(img, -o.Rotation)

	bounds = img.Bounds()
	pos := bounds.Add(image.Pt(int(o.X), int(o.Y)).Sub(originPoint))

	if layer.Opacity < 1 {
		mask := image.NewUniform(color.Alpha{uint8(layer.Opacity * 255)})

		draw.DrawMask(r.Result, pos, img, img.Bounds().Min, mask, mask.Bounds().Min, draw.Over)
	} else {
		draw.Draw(r.Result, pos, img, img.Bounds().Min, draw.Over)
	}

	return nil
}

func (r *Renderer) _rotateObjectImage(img image.Image, rotation float64) (newImage image.Image, originPoint image.Point) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	points := []image.Point{
		image.Pt(0, 0),
		image.Pt(w-1, 0),
		image.Pt(w-1, h-1),
		image.Pt(0, h-1),
	}

	sin, cos := math.Sincos(-math.Pi * rotation / 180)

	rotatedPointsX := []float64{}
	rotatedPointsY := []float64{}

	for _, p := range points {
		x := float64(p.X)
		y := float64(p.Y)

		rotatedPointsX = append(rotatedPointsX, x*cos-y*sin)
		rotatedPointsY = append(rotatedPointsY, x*sin+y*cos)
	}

	rotatedMinX := rotatedPointsX[0]
	rotatedMinY := rotatedPointsY[0]

	for i := 1; i < 4; i++ {
		rotatedMinX = math.Min(rotatedMinX, rotatedPointsX[i])
		rotatedMinY = math.Min(rotatedMinY, rotatedPointsY[i])
	}

	originPoint = image.Pt(int(rotatedPointsX[3]-rotatedMinX), int(rotatedPointsY[3]-rotatedMinY))

	return imaging.Rotate(img, rotation, color.RGBA{}), originPoint
}
