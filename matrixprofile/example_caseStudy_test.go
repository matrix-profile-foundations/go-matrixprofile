package matrixprofile

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aouyang1/go-matrixprofile/siggen"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func Points(a []float64, n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := 0; i < n; i++ {
		pts[i].X = float64(i)
		if i < len(a) {
			pts[i].Y = a[i]
		}
	}
	return pts
}

func CreatePlot(pts []plotter.XYs, labels []string, title string) (*plot.Plot, error) {
	if labels != nil && len(pts) != len(labels) {
		return nil, fmt.Errorf("number of XYs, %d, does not match number of labels, %d", len(pts), len(labels))
	}

	p, err := plot.New()
	if err != nil {
		return p, err
	}

	p.Title.Text = title
	for i := 0; i < len(pts); i++ {
		line, points, err := plotter.NewLinePoints(pts[i])
		if err != nil {
			return p, err
		}
		line.Color = plotutil.Color(i)
		points.Color = plotutil.Color(i)
		points.Shape = nil
		p.Add(line, points)
		if labels != nil {
			p.Legend.Add(labels[i], line)
		}
	}
	return p, err
}

func PlotMP(sigPts, mpPts, cacPts plotter.XYs, motifPts [][]plotter.XYs, discordPts []plotter.XYs, discordLabels []string, filename string) error {
	var err error
	rows, cols := len(motifPts), 2
	if rows < 4 {
		rows = 4
	}
	plots := make([][]*plot.Plot, rows)

	for i := 0; i < len(motifPts); i++ {
		plots[i] = make([]*plot.Plot, cols)
	}

	plots[0][0], err = CreatePlot([]plotter.XYs{sigPts}, nil, "signal")
	if err != nil {
		return err
	}

	plots[1][0], err = CreatePlot([]plotter.XYs{mpPts}, nil, "matrix profile")
	if err != nil {
		return err
	}

	plots[2][0], err = CreatePlot([]plotter.XYs{cacPts}, nil, "corrected arc curve")
	if err != nil {
		return err
	}

	plots[3][0], err = CreatePlot(discordPts, discordLabels, "discords")
	if err != nil {
		return err
	}

	for i := 0; i < len(motifPts); i++ {
		plots[i][1], err = CreatePlot(motifPts[i], nil, fmt.Sprintf("motif %d", i))
		if err != nil {
			return err
		}
	}

	img := vgimg.New(vg.Points(1200), vg.Points(600))
	dc := draw.New(img)

	t := draw.Tiles{
		Rows: rows,
		Cols: cols,
	}

	canvases := plot.Align(plots, t, dc)
	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			if plots[j][i] != nil {
				plots[j][i].Draw(canvases[j][i])
			}
		}
	}

	w, err := os.Create(filename)
	if err != nil {
		return err
	}

	png := vgimg.PngCanvas{Canvas: img}
	_, err = png.WriteTo(w)
	return err
}

func Example_caseStudy() {
	sin := siggen.Sin(1, 5, 0, 0, 100, 2)
	sin2 := siggen.Sin(0.25, 10, 0, 0.75, 100, 0.25)
	saw := siggen.Sawtooth(0.5, 7, 0, 0, 100, 1)
	noise := siggen.Noise(0.3, len(sin2)*4)
	sig := siggen.Append(sin, sin2, sin, noise, sin2, noise, sin2, noise, saw)

	noise = siggen.Noise(0.1, len(sig))
	sig = siggen.Add(sig, noise)

	var m, k int
	var r float64
	m = 32
	k = 6
	r = 3
	mp, err := New(sig, nil, m)
	if err != nil {
		panic(err)
	}

	if err = mp.Stomp(2); err != nil {
		panic(err)
	}

	_, _, cac := mp.Segment()

	motifs, err := mp.TopKMotifs(k, r)
	if err != nil {
		panic(err)
	}

	discords := mp.Discords(3, mp.m/2)
	if err != nil {
		panic(err)
	}

	sigPts := Points(sig, len(sig))
	mpPts := Points(mp.MP, len(sig))
	cacPts := Points(cac, len(sig))
	motifPts := make([][]plotter.XYs, k)
	discordPts := make([]plotter.XYs, k)
	discordLabels := make([]string, k)

	for i := 0; i < k; i++ {
		motifPts[i] = make([]plotter.XYs, len(motifs[i].Idx))
	}

	for i := 0; i < k; i++ {
		for j, idx := range motifs[i].Idx {
			motifPts[i][j] = Points(sig[idx:idx+m], m)
		}
	}

	for i, idx := range discords {
		discordPts[i] = Points(sig[idx:idx+m], m)
		discordLabels[i] = strconv.Itoa(idx)
	}

	if err = PlotMP(sigPts, mpPts, cacPts, motifPts, discordPts, discordLabels, "../mp_sine.png"); err != nil {
		panic(err)
	}

	fmt.Println("Saved png file result to mp_sine.png")
	// Output: Saved png file result to mp_sine.png
}
