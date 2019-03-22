package main

import (
	"fmt"
	"io/ioutil"

	"github.com/fogleman/gg"

	geojson "github.com/paulmach/go.geojson"
)

func main() {
	dc := gg.NewContext(1366, 1024)
	dc.SetHexColor("fff")

	dc.InvertY()
	dc.Scale(6, 10)

	rawFeatureJSON, err := ioutil.ReadFile("borders.geojson")
	if err != nil {
		panic(err)
	}

	fc1, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	if err != nil {
		panic(err)
	}

	dc.SetLineWidth(1.0)
	dc.SetRGB(1, 0, 0)
	drawBackground("888", dc)
	for _, multiPolygon := range fc1.Features[0].Geometry.MultiPolygon {
		for _, polygon := range multiPolygon {
			drawPolygon(polygon, dc)
		}

	}

	dc.SetHexColor("f00")
	dc.Fill()
	dc.SavePNG("out.png")

}

//
func drawPolygon(polygon [][]float64, c *gg.Context) {
	c.MoveTo(polygon[0][0], polygon[0][1])
	for i := 0; i < len(polygon); i++ {
		fmt.Println(polygon[i][0])
		fmt.Println(polygon[i][1])
		c.LineTo(polygon[i][0], polygon[i][1])
	}
}

func drawBackground(color string, c *gg.Context) {
	c.MoveTo(0, 0)
	c.LineTo(1366, 0)
	c.LineTo(1366, 1024)
	c.LineTo(0, 1024)
	c.SetHexColor(color)
	c.Fill()
}
