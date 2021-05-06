/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package appbackend

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func addText(mw *imagick.MagickWand, text, color string, size, stroke, x, y float64) {
	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	dw.SetFont("/usr/local/share/appbackend/plume.otf")
	dw.SetFontSize(size)

	pw.SetColor("white")
	dw.SetStrokeColor(pw)
	dw.SetStrokeWidth(stroke)

	pw.SetColor(color)
	dw.SetFillColor(pw)
	dw.Annotation(x, y, text)

	mw.DrawImage(dw)
}

func addPic(mw *imagick.MagickWand, file string, x, y, scale float64) {
	pic := imagick.NewMagickWand()
	defer pic.Destroy()

	pic.ReadImage(file)

	dw := imagick.NewDrawingWand()
	dw.Composite(imagick.COMPOSITE_OP_ATOP, x, y, float64(pic.GetImageWidth())*scale, float64(pic.GetImageHeight())*scale, pic)

	mw.DrawImage(dw)
}

func drawGraphLine(mw *imagick.MagickWand, pts []imagick.PointInfo, color string) {
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()

	dw.SetStrokeAntialias(true)
	dw.SetStrokeWidth(2)
	dw.SetStrokeLineCap(imagick.LINE_CAP_ROUND)
	dw.SetStrokeLineJoin(imagick.LINE_JOIN_ROUND)

	cw.SetColor(color)
	dw.SetStrokeColor(cw)

	cw.SetColor("none")
	dw.SetFillColor(cw)

	dw.Polyline(pts)

	mw.DrawImage(dw)
}

func drawGraphBackground(mw *imagick.MagickWand, pts []imagick.PointInfo, color string) {
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()

	dw.SetStrokeAntialias(true)
	dw.SetStrokeWidth(2)
	dw.SetStrokeLineCap(imagick.LINE_CAP_ROUND)
	dw.SetStrokeLineJoin(imagick.LINE_JOIN_ROUND)

	cw.SetColor("none")
	dw.SetStrokeColor(cw)

	cw.SetColor(color)
	cw.SetOpacity(0.4)
	dw.SetFillColor(cw)

	dw.Polygon(pts)

	mw.DrawImage(dw)
}

func addGraph(mw *imagick.MagickWand, x, y, width, height, min, max float64, mv TimeSeries, color string) {
	var (
		spanX = width / float64(len(mv)-1)
	)

	pts := make([]imagick.PointInfo, 0, len(mv)+2)
	pts = append(pts, imagick.PointInfo{
		X: x, Y: y,
	})
	for i, v := range mv {
		pts = append(pts, imagick.PointInfo{
			X: x + float64(i)*spanX,
			Y: y - ((v[1] - min) * (height - 60) / (max - min)),
		})
	}
	pts = append(pts, imagick.PointInfo{
		X: x + width, Y: y,
	})

	drawGraphBackground(mw, []imagick.PointInfo{
		{x, y}, {x + width, y}, {x + width, y - height}, {x, y - height},
	}, "white")
	drawGraphLine(mw, pts[1:len(pts)-1], color)
	drawGraphBackground(mw, pts, color)

	cw := imagick.NewPixelWand()
	defer cw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw.SetColor("white")
	dw.SetStrokeColor(cw)
	dw.SetStrokeWidth(3)
	dw.Line(x, y, x, y-height)
	dw.Line(x, y, x+width, y)

	mw.DrawImage(dw)
}

func init() {
	imagick.Initialize()
}

func AddSGLOverlays(box Box, plant Plant, meta MetricsMeta, img *bytes.Buffer) (*bytes.Buffer, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	mw.ReadImageBlob(img.Bytes())

	addText(mw, box.Name, "#3BB30B", 55, 2, 10, 50)
	addText(mw, plant.Name, "#FF4B4B", 45, 2, 10, 100)

	if meta.Temperature != nil || meta.Humidity != nil {
		var (
			x = float64(5)
			y = float64(mw.GetImageHeight() - 5)
		)
		if meta.Temperature != nil {
			addGraph(mw, x, y, 210, 140, 10, 40, *meta.Temperature, "#3BB30B")
			addText(mw, fmt.Sprintf("%d°C", int(meta.Temperature.current())), "#3BB30B", 60, 2, x+20, y-120)
			addText(mw, fmt.Sprintf("(%d°F)", int(meta.Temperature.current()*9/5+32)), "#3BB30B", 40, 2, x+20, y-80)
		}

		if meta.Humidity != nil {
			addGraph(mw, x+225, y, 210, 140, 10, 90, *meta.Humidity, "#0B81B3")
			addText(mw, fmt.Sprintf("%d%%", int(meta.Humidity.current())), "#0B81B3", 60, 2, x+245, y-120)
		}
	}

	t := meta.Date
	d := t.Format("2006/01/02")
	addText(mw, d, "#3BB30B", 25, 1, float64(mw.GetImageWidth()-170), float64(mw.GetImageHeight()-15))
	d = t.Format("15:04")
	addText(mw, d, "#3BB30B", 35, 1, float64(mw.GetImageWidth()-170), float64(mw.GetImageHeight()-40))

	addPic(mw, "/usr/local/share/appbackend/watermark-logo.png", float64(mw.GetImageWidth()-100), 10, 0.3)

	return bytes.NewBuffer(mw.GetImageBlob()), nil
}

func AddSGLOverlaysForFile(box Box, plant Plant, meta MetricsMeta, file string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer(f)

	buff, err = AddSGLOverlays(box, plant, meta, buff)
	if err != nil {
		return err
	}

	if ioutil.WriteFile(file, buff.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
