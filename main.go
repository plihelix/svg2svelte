// This small program will convert a given svg into a svelte file.
package main

import (
	"fmt"
	"os"

	"github.com/plihelix/svg2svelte/svg"
)

type prepSvg struct {
	Filename string
	ViewBox  string
	Height   string
	Width    string
	Groups   []prepGroup
}

type prepGroup struct {
	Transform string
	Paths     []prepPath
	Rects     []prepRect
	Circles   []prepCircle
}

type prepPath struct {
	D    string
	Fill string
}

type prepRect struct {
	Width     string
	Height    string
	Transform string
	Style     string
	Rx        string
	Ry        string
	Fill      string
}

type prepCircle struct {
	Cx        string
	Cy        string
	R         string
	Transform string
	Style     string
	Fill      string
}

type clrs []colormap

type colormap struct {
	Hex    string
	ScrVar string
	SvgVar string
}

func main() {
	fmt.Println("Convert a .svg icon to a .svelte and extract the colors into external variables.")
	fmt.Print("Enter the name of the .svg file: ")
	var filename string
	fmt.Scanln(&filename)
	var iconname string
	// the name of the icon is the name of the file without the extension
	iconname = filename[:len(filename)-4]

	// Setup a reader for the file.
	r, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Opened file: ", filename)

	// use the svg package to unmarshal the svg
	icon, err := svg.ParseSvgFromReader(r, iconname, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("The icon has been parsed.")

	fmt.Printf("ViewBox: %s\n", icon.ViewBox)
	fmt.Printf("Height: %s\n", icon.Height)
	fmt.Printf("Width: %s\n", icon.Width)
	var paths []svg.Path
	var rects []svg.Rect
	var circles []svg.Circle
	for i, g := range icon.Groups {
		paths = icon.Groups[i].GetPaths()
		rects = icon.Groups[i].GetRects()
		circles = icon.Groups[i].GetCircles()
		fmt.Printf("Group: %d has %d paths.\n", i, len(paths))
		fmt.Printf("  TransformString: %s\n", g.TransformString)
	}
	for i, p := range paths {
		fmt.Printf("Path: %d\n", i)
		fmt.Printf("  d: %v\n", p.D)
		fmt.Printf("  fill: %s\n", p.Fill)
		// fmt.Printf("  style: %s\n", p.Style)
	}
	for i, r := range rects {
		fmt.Printf("Rect: %d\n", i)
		// fmt.Printf("  width: %s\n", r.Width)
		// fmt.Printf("  height: %s\n", r.Height)
		// fmt.Printf("  transform: %s\n", r.Transform)
		// fmt.Printf("  style: %s\n", r.Style)
		// fmt.Printf("  rx: %s\n", r.Rx)
		// fmt.Printf("  ry: %s\n", r.Ry)
		fmt.Printf("  fill: %s\n", r.Fill)
	}
	for i, c := range circles {
		fmt.Printf("Circle: %d\n", i)
		// fmt.Printf("  cx: %s\n", c.Cx)
		// fmt.Printf("  cy: %s\n", c.Cy)
		// fmt.Printf("  r: %s\n", c.R)
		// fmt.Printf("  transform: %s\n", c.Transform)
		// fmt.Printf("  style: %s\n", c.Style)
		fmt.Printf("  fill: %s\n", c.Fill)
	}

	sveltePrep := prepSvg{
		Filename: iconname + ".svelte",
		ViewBox:  icon.ViewBox,
		Height:   icon.Height,
		Width:    icon.Width,
	}

	for _, g := range icon.Groups {
		prepGroup := prepGroup{
			Transform: g.TransformString,
		}
		for _, p := range paths {
			prepPath := prepPath{
				D:    p.D,
				Fill: p.Fill,
			}
			prepGroup.Paths = append(prepGroup.Paths, prepPath)
		}
		sveltePrep.Groups = append(sveltePrep.Groups, prepGroup)
		for _, r := range rects {
			prepRect := prepRect{
				Width:     r.Width,
				Height:    r.Height,
				Transform: r.Transform,
				Style:     r.Style,
				Rx:        r.Rx,
				Ry:        r.Ry,
				Fill:      r.Fill,
			}
			prepGroup.Rects = append(prepGroup.Rects, prepRect)
		}
		for _, c := range circles {
			prepCircle := prepCircle{
				Cx:        fmt.Sprint(c.Cx),
				Cy:        fmt.Sprint(c.Cy),
				R:         fmt.Sprint(c.Radius),
				Transform: c.Transform,
				Style:     c.Style,
				Fill:      c.Fill,
			}
			prepGroup.Circles = append(prepGroup.Circles, prepCircle)
		}
	}
	// Show the built sveltePrep
	// fmt.Printf("Svelte Prep: %v\n", sveltePrep)
	fmt.Println("Svelte Prep: Completed")
	fmt.Println("Extracting colors from the paths...")

	pathcolors := ExtractPathColors(paths)
  rectcolors := ExtractRectColors(rects)
  circlecolors := ExtractCircleColors(circles)
  colors := UniqueColors(pathcolors, rectcolors, circlecolors)
	fmt.Printf("Path Colors: %v\n", colors)

	// Use the sveltePrep and colors to create the svelte file.
	result := GenerateSvelte(sveltePrep, colors)

	// save the svelte file
	saveSvelteFile(result, sveltePrep.Filename)
}

func ExtractPathColors(paths []svg.Path) []string {
	var colors []string

	for _, p := range paths {
		if !Contains(colors, p.Fill) {
			colors = append(colors, p.Fill)
		}
	}
	return colors
}

func ExtractRectColors(rects []svg.Rect) []string{
  	var colors []string

	for _, r := range rects {
		if !Contains(colors, r.Fill) {
			colors = append(colors, r.Fill)
		}
	}
	return colors
}

func ExtractCircleColors(circles []svg.Circle) []string{
	var colors []string

	for _, c := range circles {
		if !Contains(colors, c.Fill) {
			colors = append(colors, c.Fill)
		}
	}
	return colors
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GenerateSvelte(prep prepSvg, colors []string) string {
	fmt.Println("Generating the svelte file...")

	var svelteColors clrs
	var svelteOut string

	// For each colors, create a variable
	for i, c := range colors {
		svelteColors = append(svelteColors, colormap{
			Hex:    c,
			ScrVar: fmt.Sprintf("color%d", i),
			SvgVar: fmt.Sprintf("{color%d}", i),
		})
	}

	// replace the fill colors with the variables in prep
	for _, g := range prep.Groups {
		for _, p := range g.Paths {
			p.Fill = svelteColors.findSvgVarByColor(p.Fill)
		}
	}

	// Write the svelte file script section
	svelteOut = "<script>\n"
	for _, c := range colors {
		varName := svelteColors.findScrVarByColor(c)
		svelteOut += fmt.Sprintf("  export let %s = \"%s\";\n", varName, c)
	}
	svelteOut += "</script>\n"

	// Write the svelte output for the svg.
	svelteOut += fmt.Sprintf("<svg viewBox=\"%s\" height=\"%s\" width=\"%s\">\n", prep.ViewBox, prep.Height, prep.Width)
	for _, g := range prep.Groups {
		svelteOut += fmt.Sprintf("  <g transform=\"%s\">\n", g.Transform)
		for _, p := range g.Paths {
			svelteOut += fmt.Sprintf("    <path d=\"%s\" fill=\"%s\" />\n", p.D, svelteColors.findSvgVarByColor(p.Fill))
		}
    for _, r := range g.Rects {
      svelteOut += fmt.Sprintf("    <rect width=\"%s\" height=\"%s\" transform=\"%s\" style=\"%s\" rx=\"%s\" ry=\"%s\" fill=\"%s\" />\n", r.Width, r.Height, r.Transform, r.Style, r.Rx, r.Ry, r.Fill)
    }
    for _, c := range g.Circles {
      svelteOut += fmt.Sprintf("    <circle cx=\"%s\" cy=\"%s\" r=\"%s\" transform=\"%s\" style=\"%s\" fill=\"%s\" />\n", c.Cx, c.Cy, c.R, c.Transform, c.Style, c.Fill)
    }
		svelteOut += "  </g>\n"
	}
	svelteOut += "</svg>\n"

	fmt.Println()
	fmt.Println("Test output:")
	fmt.Println(svelteOut)

	return svelteOut
}

func (cm *clrs) findSvgVarByColor(s string) string {
	for _, c := range *cm {
		if c.Hex == s {
			return c.SvgVar
		}
	}
	return ""
}

func (cm *clrs) findScrVarByColor(s string) string {
	for _, c := range *cm {
		if c.Hex == s {
			return c.ScrVar
		}
	}
	return ""
}

func saveSvelteFile(svelteOut, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	f.WriteString(svelteOut)
}

func UniqueColors(pathcolors, rectcolors, circlecolors []string) []string {
  var colors []string
  colors = append(colors, pathcolors...)
  for _, c := range rectcolors {
    if !Contains(colors, c) {
      colors = append(colors, c)
    }
  }
  for _, c := range circlecolors {
    if !Contains(colors, c) {
      colors = append(colors, c)
    }
  }
  return colors
}
