package font

import (
	"errors"
	"io"
	"os"
	"unsafe"

	"golang.org/x/image/font"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

// LoadedFont struct for storing info about each loaded font
type LoadedFont struct {
	Font   *truetype.Font
	Size   fixed.Int26_6
	Glyphs []GlyphVBO
}

// GlyphVBO stores pointers to OpenGL VBOs and things
type GlyphVBO struct {
	GLVAO        *uint32
	Count        int32
	Glyph        truetype.GlyphBuf
	AdvanceWidth float64
}

// GetIndicesForString gets the glyph indices for each character in the string
func GetIndicesForString(font *LoadedFont, str string) []truetype.Index {
	var indices []truetype.Index
	for _, char := range str {
		indices = append(indices, font.Font.Index(char))
	}
	return indices
}

// GetKernsForIndices gets the kernings for the indices
func GetKernsForIndices(font *LoadedFont, indices []truetype.Index) []float64 {
	var kerns []float64
	var last truetype.Index
	for i, index := range indices {
		if i == 0 {
			last = index
			continue
		}
		fixedkern := font.Font.Kern(font.Size, last, index)
		kerns = append(kerns, float64(fixedkern))
		last = index
	}
	kerns = append(kerns, 0.0)
	return kerns
}

// DrawSlowString draws a string slowly by calculating and unpacking vertices on the fly
func DrawSlowString(font *LoadedFont, str string, x, y float64) {
	// get glyph indices for the runes in the string
	indices := GetIndicesForString(font, str)

	// get kernings for the indices
	kerns := GetKernsForIndices(font, indices)
	dx := 0.0

	// draw each glyph
	for i, index := range indices {
		firstOnPoint := true
		var curPoints []truetype.Point
		glyphBuf := font.Glyphs[index].Glyph
		glyphPoints := make([]truetype.Point, len(glyphBuf.Points)+len(glyphBuf.Ends))
		glyphEnds := make([]int, len(glyphBuf.Ends))

		// expand points to make complete contours
		for end := range glyphBuf.Ends {
			if end == 0 {
				copy(glyphPoints[:glyphBuf.Ends[end]+end], glyphBuf.Points[:glyphBuf.Ends[end]])
				glyphPoints[glyphBuf.Ends[end]+end] = glyphBuf.Points[0]
			} else {
				copy(glyphPoints[glyphBuf.Ends[end-1]+end:glyphBuf.Ends[end]+end], glyphBuf.Points[glyphBuf.Ends[end-1]:glyphBuf.Ends[end]])
				glyphPoints[glyphBuf.Ends[end]+end] = glyphBuf.Points[glyphBuf.Ends[end-1]]
			}
			glyphEnds[end] = glyphBuf.Ends[end] + end
		}

		// draw contours
		var endIndex int
		for b, point := range glyphPoints {
			if b == glyphEnds[endIndex]+1 {
				endIndex++
				firstOnPoint = true
				curPoints = []truetype.Point{}
			}

			if point.Flags&1 == 0 {
				curPoints = append(curPoints, point)
			} else if firstOnPoint == false {
				curPoints = append(curPoints, point)

				gl.Color3d(1, 1, 1)
				gl.Begin(gl.LINE_STRIP)
				for d := 0.0; d <= 1.0; d += 0.1 {

					if len(curPoints) == 2 {
						_x, _y := LinearBézier(d, float64(curPoints[0].X), float64(curPoints[0].Y), float64(curPoints[1].X), float64(curPoints[1].Y))
						gl.Vertex2d(x+_x+dx, y-_y)
					} else if len(curPoints) == 3 {
						_x, _y := QuadraticBézier(d, float64(curPoints[0].X), float64(curPoints[0].Y), float64(curPoints[1].X), float64(curPoints[1].Y), float64(curPoints[2].X), float64(curPoints[2].Y))
						gl.Vertex2d(x+_x+dx, y-_y)
					} else {
						_x, _y := UnpackBézier(d, curPoints)
						gl.Vertex2d(x+_x+dx, y-_y)
					}
				}
				gl.End()
				curPoints = []truetype.Point{point}
				firstOnPoint = false
			} else {
				firstOnPoint = false
				curPoints = append(curPoints, point)
			}
		}
		dx += font.Glyphs[index].AdvanceWidth + kerns[i]

	}
}

// LoadFont loads a font at the specified scale
func LoadFont(path string, scale float64) (*LoadedFont, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	strct, err := file.Stat()

	if err != nil {
		return nil, err
	}

	filedata := make([]byte, strct.Size())

	_, err = io.ReadFull(file, filedata)

	if err != nil {
		return nil, err
	}

	fnt, err := truetype.Parse(filedata)

	if err != nil {
		return nil, err
	}

	hint := font.HintingNone

	lfont := LoadedFont{
		Font:   fnt,
		Size:   fixed.Int26_6(scale),
		Glyphs: []GlyphVBO{},
	}
	for i := 0; i < 4096; i++ {
		vao, cnt, err := buildGlyph(fnt, i, fixed.Int26_6(scale), hint)
		if err != nil {
			continue
		}
		glyphbuf := truetype.GlyphBuf{}
		err = glyphbuf.Load(fnt, fixed.Int26_6(scale), truetype.Index(i), hint)
		if err != nil {
			continue
		}
		glyphvbo := GlyphVBO{
			GLVAO:        vao,
			Count:        cnt,
			Glyph:        glyphbuf,
			AdvanceWidth: float64(fnt.HMetric(fixed.Int26_6(scale), truetype.Index(i)).AdvanceWidth),
		}
		lfont.Glyphs = append(lfont.Glyphs, glyphvbo)
	}

	return &lfont, nil

}

// buildGlyph builds a glyph VBO for the specified font
func buildGlyph(font *truetype.Font, index int, scale fixed.Int26_6, hint font.Hinting) (*uint32, int32, error) {

	vao := [1]uint32{}
	gl.GenVertexArrays(1, &vao[0])
	gl.BindVertexArray(vao[0])

	var quit bool

	defer func() {
		if recover() != nil {
			quit = true
		}
	}()

	glyphbuf := truetype.GlyphBuf{}
	err := glyphbuf.Load(font, scale, truetype.Index(index), hint)

	if quit {
		return nil, 0, errors.New("Could not load glyph")
	}

	if err != nil {
		return nil, 0, err
	}

	vbo := [1]uint32{}
	gl.GenBuffers(1, &vbo[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])

	size := int32(len(glyphbuf.Points))

	for _, pt := range glyphbuf.Points {
		gpt := [2]float64{float64(pt.X), float64(pt.Y)}
		gl.BufferData(vbo[0], 2, unsafe.Pointer(&gpt[0]), gl.STATIC_DRAW)
	}

	indices := []int32{}
	lastindex := int32(0)
	for _, index := range glyphbuf.Ends {
		for i := lastindex; i < int32(index); i++ {
			indices = append(indices, i, i+1)
		}
	}

	gl.VertexAttribPointer(0, size, gl.DOUBLE, false, 0, unsafe.Pointer(nil))
	gl.EnableVertexAttribArray(0)

	// generate VBO indices
	//vboIndices := [1]uint32{}
	//gl.GenBuffers(1, &vboIndices[0])
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboIndices[0])
	//gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices), unsafe.Pointer(&indices[0]), gl.STATIC_DRAW)

	gl.BindVertexArray(0)

	return &vao[0], size, nil

}

// LinearBézier computes the points on a bézier curve given 2 points and a t value.
func LinearBézier(t, p0x, p0y, p1x, p1y float64) (float64, float64) {
	x := (1-t)*p0x + t*p1x
	y := (1-t)*p0y + t*p1y
	return x, y
}

// QuadraticBézier computes the points on a bézier curve given 3 points and a t value.
func QuadraticBézier(t, p0x, p0y, p1x, p1y, p2x, p2y float64) (float64, float64) {
	x := (1-t)*((1-t)*p0x+t*p1x) + t*((1-t)*p1x+t*p2x)
	y := (1-t)*((1-t)*p0y+t*p1y) + t*((1-t)*p1y+t*p2y)
	return x, y
}

// CubicBézier computes the points on a bézier curve given 4 points and a t value.
func CubicBézier(t, p0x, p0y, p1x, p1y, p2x, p2y, p3x, p3y float64) (float64, float64) {
	x := (1-t)*(1-t)*(1-t)*p0x + 3*(1-t)*(1-t)*t*p1x + 3*(1-t)*t*t*p2x + t*t*t*p3x
	y := (1-t)*(1-t)*(1-t)*p0y + 3*(1-t)*(1-t)*t*p1y + 3*(1-t)*t*t*p2y + t*t*t*p3y
	return x, y
}

// UnpackBézier turns the incomplete bézier curves into their respective complete curves
func UnpackBézier(t float64, points []truetype.Point) (float64, float64) {
	subcurves := len(points) - 2
	newpoints := []truetype.Point{points[0]}
	for i := 1; i < len(points)-2; i++ {
		newpoints = append(newpoints, points[i])
		midx, midy := LinearBézier(0.5, float64(points[i].X), float64(points[i].Y), float64(points[i+1].X), float64(points[i+1].Y))
		newpoints = append(newpoints, truetype.Point{
			X:     fixed.Int26_6(midx),
			Y:     fixed.Int26_6(midy),
			Flags: 1,
		})
	}
	newpoints = append(newpoints, points[len(points)-2])
	newpoints = append(newpoints, points[len(points)-1])
	curveindex := int32(t * float64(subcurves))
	x, y := QuadraticBézier(t*float64(subcurves)-float64(curveindex), float64(newpoints[curveindex*2].X), float64(newpoints[curveindex*2].Y), float64(newpoints[curveindex*2+1].X), float64(newpoints[curveindex*2+1].Y), float64(newpoints[curveindex*2+2].X), float64(newpoints[curveindex*2+2].Y))

	return x, y
}

/*func UnpackBézier2(points []truetype.Point) []truetype.Point {
	//subcurves := len(points) - 2
	newpoints := []truetype.Point{points[0]}
	for i := 1; i < len(points)-2; i++ {
		newpoints = append(newpoints, points[i])
		midx, midy := LinearBézier(0.5, float64(points[i].X), float64(points[i].Y), float64(points[i+1].X), float64(points[i+1].Y))
		newpoints = append(newpoints, truetype.Point{
			X:     fixed.Int26_6(midx),
			Y:     fixed.Int26_6(midy),
			Flags: 1,
		})
	}
	newpoints = append(newpoints, points[len(points)-2])
	newpoints = append(newpoints, points[len(points)-1])
	return newpoints
}

func Bézier(t float64, points []truetype.Point) (float64, float64) {
	var x, y float64
	order := len(points) - 1
	for i := 0; i <= order; i++ {
		coeff := float64(Binomial(int32(order), int32(i))) * math.Pow(1-t, float64(order-i)) * math.Pow(t, float64(i))
		x += coeff * float64(points[i].X)
		y += coeff * float64(points[i].Y)
	}
	return x, y
}

func Binomial(n, k int32) int32 {
	return Factorial(n) / (Factorial(k) * Factorial(n-k))
}

func Factorial(x int32) int32 {
	if x <= 1 {
		return 1
	}
	return x * Factorial(x-1)
}*/
