package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/fogleman/gg"

	geojson "github.com/paulmach/go.geojson"
)

func main() {
	var fileDest = flag.String("d", "borders.geojson", "Destination of file")
	flag.Parse()
	fmt.Println(*fileDest)
	dc := gg.NewContext(1366, 1024)

	dc.InvertY()
	dc.Scale(7, 12)

	rawFeatureJSON, err := ioutil.ReadFile(*fileDest)
	if err != nil {
		panic(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	if err != nil {
		panic(err)
	}

	err = drawMap(fc, dc)
	if err != nil {
		panic(err)
	}
	dc.SavePNG("out.png")

}

func drawMap(fc *geojson.FeatureCollection, c *gg.Context) error {
	drawBackground("000", c)
	for _, feature := range fc.Features {
		name, err := feature.PropertyString("name")
		if err != nil {
			return err
		}
		fmt.Println(name)
		drawGeometry(feature.Geometry, c)
	}
	return nil
}

func drawGeometry(g *geojson.Geometry, c *gg.Context) {
	if g.IsCollection() {
		for _, geometry := range g.Geometries {
			drawGeometry(geometry, c)
		}
		return
	}
	if g.IsLineString() {
		drawLine(g.LineString, c)
		return
	}
	if g.IsMultiLineString() {
		for _, line := range g.MultiLineString {
			drawLine(line, c)
		}
		return
	}
	if g.IsMultiPoint() {
		for _, point := range g.MultiPoint {
			drawPoint(point, c)
		}
		return
	}
	if g.IsMultiPolygon() {
		for _, multiPolygon := range g.MultiPolygon {
			for _, polygon := range multiPolygon {
				drawPolygon(polygon, c)
			}
		}
		return
	}
	if g.IsPoint() {
		drawPoint(g.Point, c)
		return
	}
	if g.IsPolygon() {
		for _, polygon := range g.Polygon {
			drawPolygon(polygon, c)
		}
		return
	}
}

func drawLine(l [][]float64, c *gg.Context) {
	c.SetHexColor("0f0")
	c.MoveTo(l[0][0], l[0][1])
	for _, p := range l {
		c.LineTo(p[0], p[1])
	}
}

func drawPoint(p []float64, c *gg.Context) {
	c.SetHexColor("00f")
	c.DrawPoint(p[0], p[1], 5)
}

func drawPolygon(polygon [][]float64, c *gg.Context) {
	c.SetRGB(0, 0, 0)
	c.SetLineWidth(1.0)
	c.MoveTo(polygon[0][0], polygon[0][1])
	var x float64
	var y float64
	for i := 0; i < len(polygon); i++ {
		if polygon[i][0] < 0 {
			x = polygon[i][0] + 180*2
		} else {
			x = polygon[i][0]
		}
		y = polygon[i][1]
		c.LineTo(x, y)
	}
	c.SetHexColor("fff")
	c.Fill()
}

func drawBackground(color string, c *gg.Context) {
	c.MoveTo(0, 0)
	c.LineTo(1366, 0)
	c.LineTo(1366, 1024)
	c.LineTo(0, 1024)
	c.SetHexColor(color)
	c.Fill()
}
