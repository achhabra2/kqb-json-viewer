package stats

import (
	"image"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

const dpi = 96

func PlotStats(data StatsJSON) image.Image {
	stats := data.AdvancedStats()

	labels := make([]string, 8)
	QueenKills := make(plotter.Values, 8)
	WarriorKills := make(plotter.Values, 8)
	WorkerKills := make(plotter.Values, 8)
	WarriorDeaths := make(plotter.Values, 8)
	WorkerDeaths := make(plotter.Values, 8)

	counter := 0
	for name, playerStats := range stats {
		labels[counter] = name
		// QueenKills = append(QueenKills, float64(playerStats["QueenKills"]))
		QueenKills[counter] = float64(playerStats["QueenKills"])
		WarriorKills[counter] = float64(playerStats["WarriorKills"])
		WorkerKills[counter] = float64(playerStats["WorkerKills"])
		WarriorDeaths[counter] = float64(playerStats["WarriorDeaths"])
		WorkerDeaths[counter] = float64(playerStats["WorkerDeaths"])
		counter++
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Advanced Stats"
	p.Y.Label.Text = "Count"

	w := vg.Points(20)

	barsA, err := plotter.NewBarChart(QueenKills, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -2 * w

	barsB, err := plotter.NewBarChart(WarriorKills, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)
	barsB.Offset = -w

	barsC, err := plotter.NewBarChart(WorkerKills, w)
	if err != nil {
		panic(err)
	}
	barsC.LineStyle.Width = vg.Length(0)
	barsC.Color = plotutil.Color(2)

	barsD, err := plotter.NewBarChart(WarriorDeaths, w)
	if err != nil {
		panic(err)
	}
	barsD.LineStyle.Width = vg.Length(0)
	barsD.Color = plotutil.Color(3)
	barsD.Offset = w

	barsE, err := plotter.NewBarChart(WorkerDeaths, w)
	if err != nil {
		panic(err)
	}
	barsE.LineStyle.Width = vg.Length(0)
	barsE.Color = plotutil.Color(4)
	barsE.Offset = 2 * w

	p.Add(barsA, barsB, barsC, barsD, barsE)
	p.Legend.Add("QueenKills", barsA)
	p.Legend.Add("WarriorKills", barsB)
	p.Legend.Add("WorkerKills", barsC)
	p.Legend.Add("Q/W-Deaths", barsD)
	p.Legend.Add("WorkerDeaths", barsE)
	p.Legend.Top = true
	p.NominalX(labels...)

	// Draw the plot to an in-memory image.
	img := image.NewRGBA(image.Rect(0, 0, 16*dpi, 8*dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	// if err := p.Save(16*vg.Inch, 8*vg.Inch, "barchart.png"); err != nil {
	// 	panic(err)
	// }

	return c.Image()
}
