package main

import (
	//"os"
	"runtime"

	"./ui"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	//"github.com/go-gl/gltext"
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

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	gl.Init()
	gl.ClearColor(0, 0, 0, 1)

	//file, err := os.Open("C:\\Windows\\Fonts\\times.ttf")

	//defer file.Close()
	//font, err = gltext.LoadTruetype(file, 24, 32, 127, gltext.LeftToRight)

	//defer font.Release()

	layout := ui.NewTableLayout()
	ui.CreateRenderableBox(&layout, 20, 10, 0, 0, 1, 1, [4]float32{1, 1, 1, 1})
	ui.CreateRenderableBox(&layout, 15, 10, 0, 1, 1, 1, [4]float32{1, 1, 0, 1})
	ui.CreateRenderableBox(&layout, 10, 10, 0, 2, 1, 1, [4]float32{1, 0, 1, 1})
	ui.CreateRenderableBox(&layout, 60, 10, 1, 0, 1, 2, [4]float32{1, 0, 0, 1})
	ui.CreateRenderableBox(&layout, 20, 10, 1, 2, 1, 1, [4]float32{0, 1, 1, 1})
	ui.CreateRenderableBox(&layout, 10, 10, 2, 0, 1, 1, [4]float32{0, 1, 0, 1})
	ui.CreateRenderableBox(&layout, 80, 10, 2, 1, 1, 2, [4]float32{0, 0, 1, 1})
	layout.Layout()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		width, height := window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//gl.MatrixMode(gl.MODELVIEW)
		//gl.LoadIdentity()
		//gl.MatrixMode(gl.PROJECTION)

		//x := float64(time.Now().UnixNano()) / 1000000000.0

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
		layout.Render()

		//drawString(0, 0, "test")

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
