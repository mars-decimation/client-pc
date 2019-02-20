package main

import (
	"golang.org/x/image/math/fixed"
	//"os"
	"fmt"
	"runtime"

	"./ui"
	"./ui/font"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	//"github.com/go-gl/gltext"
	"github.com/golang/freetype/truetype"
)

func init() {
	runtime.LockOSThread()
}

//var font *gltext.Font

// test
func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	//glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	//window.

	window.MakeContextCurrent()

	gl.Init()
	gl.ClearColor(0, 0, 0, 1)

	//file, err := os.Open("C:\\Windows\\Fonts\\times.ttf")

	//defer file.Close()
	//font, err = gltext.LoadTruetype(file, 24, 32, 127, gltext.LeftToRight)

	//defer font.Release()

	layout := ui.NewTableLayout()
	ui.CreateRenderableBox(&layout, 200, 100, 0, 0, 1, 1, [4]float32{1, 1, 1, 1})
	ui.CreateRenderableBox(&layout, 150, 100, 0, 1, 1, 1, [4]float32{1, 1, 0, 1})
	ui.CreateRenderableBox(&layout, 100, 100, 0, 2, 1, 1, [4]float32{1, 0, 1, 1})
	ui.CreateRenderableBox(&layout, 600, 100, 1, 0, 1, 2, [4]float32{1, 0, 0, 1})
	ui.CreateRenderableBox(&layout, 200, 100, 1, 2, 1, 1, [4]float32{0, 1, 1, 1})
	ui.CreateRenderableBox(&layout, 100, 100, 2, 0, 1, 1, [4]float32{0, 1, 0, 1})
	ui.CreateRenderableBox(&layout, 800, 100, 2, 1, 1, 2, [4]float32{0, 0, 1, 1})
	layout.Layout()

	fontsize := float64(400)

	fnt, err := font.LoadFont(fontsize)
	if err != nil {
		fmt.Printf(err.Error())
	}

	gl.ClearColor(0, 0, 0, 0)

	indices := font.GetIndicesForString(fnt, "The 8 quick brown fox jumped over the lazy dog.")
	kerns := font.GetKernsForIndices(fnt, fixed.Int26_6(fontsize), indices)

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//gl.MatrixMode(gl.MODELVIEW)
		//gl.LoadIdentity()
		//gl.MatrixMode(gl.PROJECTION)

		//z := [16]float64{math.Cos(x) * 1.7, -math.Sin(x) * 1.7, 0, 0, math.Sin(x), math.Cos(x), 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
		z := [16]float64{2.0 / float64(width), 0, 0, 0, 0, -2.0 / float64(height), 0, 0, 1, 0, 0, 0, 0, 0, 0, 1}
		t := [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, -1, 1, 0, 1}
		gl.MatrixMode(gl.MODELVIEW)

		gl.LoadMatrixd(&t[0])
		gl.MultMatrixd(&z[0])

		//gl.Begin(gl.TRIANGLES)

		//x = x - float64(int(x))

		//gl.Color3d(x, 1, 0)
		//gl.Vertex2d(-0.5, -0.5)
		//gl.Color3d(0, 1, 1)
		//gl.Vertex2d(0.5, -0.5)
		//gl.Color3d(1, 0, 1)
		//gl.Vertex2d(0, 0.5)

		//gl.End()
		//layout.Render()

		//vao := fnt.Glyphs[50].GLVAO
		//gl.BindVertexArray(*vao)
		//gl.DrawArrays(gl.LINES, 0, (*fnt).Glyphs[0].Count)

		//drawString(0, 0, "test")
		//time := time.Now().Second()
		gl.Color3d(1, 1, 1)

		//var last *truetype.Point
		dx := 0.0
		for i, index := range indices {
			//endIndex := 0
			//offPointCount := 0
			firstOnPoint := true
			var curPoints []truetype.Point
			glyphBuf := fnt.Glyphs[index].Glyph
			//var beginPoint *truetype.Point
			glyphPoints := make([]truetype.Point, len(glyphBuf.Points)+len(glyphBuf.Ends))
			glyphEnds := make([]int, len(glyphBuf.Ends))
			//glyphPoints := glyphBuf.Points
			//copy(glyphPoints, glyphBuf.Points)
			// expand points to make complete contours
			for end := range glyphBuf.Ends {
				//end := len(glyphBuf.Ends) - 1
				if end == 0 {
					copy(glyphPoints[:glyphBuf.Ends[end]+end], glyphBuf.Points[:glyphBuf.Ends[end]])
					glyphPoints[glyphBuf.Ends[end]+end] = glyphBuf.Points[0]
					//glyphPoints = append(glyphPoints, glyphBuf.Points[:glyphBuf.Ends[end]]...)
				} else {
					copy(glyphPoints[glyphBuf.Ends[end-1]+end:glyphBuf.Ends[end]+end], glyphBuf.Points[glyphBuf.Ends[end-1]:glyphBuf.Ends[end]])
					glyphPoints[glyphBuf.Ends[end]+end] = glyphBuf.Points[glyphBuf.Ends[end-1]]
				}
				glyphEnds[end] = glyphBuf.Ends[end] + end

				//glyphPoints = append(glyphPoints[:end], glyphPoints[end:]...)
			}

			/*if len(glyphPoints) > time {

				glyphPoints = glyphPoints[:time]
			}*/

			// draw contours
			var endIndex int
			for b, point := range glyphPoints {

				//if beginPoint == nil {
				//beginPoint = &glyphBuf.Points[j]
				//}

				gl.Color3d(0, 1, 1)
				gl.PointSize(3)
				gl.Begin(gl.POINTS)
				gl.Vertex2d(float64(point.X)+10+dx, 500-float64(point.Y))
				gl.End()

				if b == glyphEnds[endIndex]+1 {
					endIndex++
					firstOnPoint = true
					curPoints = []truetype.Point{}
				}

				if point.Flags&1 == 0 {

					//offPointCount++

					// debug
					/*gl.Color3d(0, 1, 0)
					gl.PointSize(3)
					gl.Begin(gl.POINTS)
					gl.Vertex2d(float64(point.X)+10+dx, 500-float64(point.Y))
					gl.End()*/
					// -----

					curPoints = append(curPoints, point)
				} else if firstOnPoint == false {
					// debug
					/*gl.Color3d(1, 0, 0)
					gl.Begin(gl.POINTS)
					gl.Vertex2d(float64(point.X)+10+dx, 500-float64(point.Y))
					gl.End()*/
					// -----

					curPoints = append(curPoints, point)

					gl.Color3d(1, 1, 1)
					gl.Begin(gl.LINE_STRIP)
					for d := 0.0; d <= 1.0; d += 0.1 {

						if len(curPoints) == 2 {
							x, y := font.LinearBézier(d, float64(curPoints[0].X), float64(curPoints[0].Y), float64(curPoints[1].X), float64(curPoints[1].Y))
							gl.Vertex2d(x+10+dx, 500-y)
						} else if len(curPoints) == 3 {
							x, y := font.QuadraticBézier(d, float64(curPoints[0].X), float64(curPoints[0].Y), float64(curPoints[1].X), float64(curPoints[1].Y), float64(curPoints[2].X), float64(curPoints[2].Y))
							gl.Vertex2d(x+10+dx, 500-y)
						} else {
							//x, y := font.Bézier(d, curPoints)
							x, y := font.UnpackBézier(d, curPoints)
							gl.Vertex2d(x+10+dx, 500-y)
						}
					}
					gl.End()
					/*if j == glyphBuf.Ends[endIndex]-1 {
						curPoints = []truetype.Point{}
						firstOnPoint = true
						endIndex++
					} else {*/
					curPoints = []truetype.Point{point}
					firstOnPoint = false
					//}
					//offPointCount = 0

					//last = &point
				} else {
					// debug
					/*gl.Color3d(1, 1, 0)
					gl.Begin(gl.POINTS)
					gl.Vertex2d(float64(point.X)+10+dx, 500-float64(point.Y))
					gl.End()*/
					// -----

					firstOnPoint = false
					curPoints = append(curPoints, point)
				}
			}
			dx += fnt.Glyphs[index].AdvanceWidth + kerns[i]

		}

		window.SwapBuffers()
		glfw.PollEvents()
	}

}

/*func drawString(x, y float32, str string) error {
	//for i := range fonts {

	// We need to offset each string by the height of the
	// font. To ensure they don't overlap each other.
	//font.GlyphBounds()

	// Draw a rectangular backdrop using the string's metrics.
	sw, sh := font.Metrics(str)
	gl.Color4f(0.1, 0.1, 0.1, 0.7)
	gl.Rectf(x, y, x+float32(sw), y+float32(sh))

	// Render the string.
	gl.Color4f(1, 1, 1, 1)
	err := font.Printf(x, y, str)
	if err != nil {
		return err
	}
	//}

	return nil
}*/
