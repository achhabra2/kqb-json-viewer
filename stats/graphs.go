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
	for _, teamStats := range stats {
		for name, playerStats := range teamStats {
			labels[counter] = name
			// QueenKills = append(QueenKills, float64(playerStats["QueenKills"]))
			QueenKills[counter] = float64(playerStats["QueenKills"])
			WarriorKills[counter] = float64(playerStats["WarriorKills"])
			WorkerKills[counter] = float64(playerStats["WorkerKills"])
			WarriorDeaths[counter] = float64(playerStats["WarriorDeaths"])
			WorkerDeaths[counter] = float64(playerStats["WorkerDeaths"])
			counter++
		}
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Military Stats"
	p.Title.Font.Size = vg.Points(20)
	p.Y.Label.Text = "Count"
	p.Y.Label.Font.Size = vg.Points(16)

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
	p.Legend.Font.Size = vg.Points(20)
	p.NominalX(labels...)
	p.X.Tick.Label.Font.Size = vg.Points(12)
	p.Y.Tick.Label.Font.Size = vg.Points(16)
	// Draw the plot to an in-memory image.
	img := image.NewRGBA(image.Rect(0, 0, 16*dpi, 9*dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	// if err := p.Save(16*vg.Inch, 8*vg.Inch, "barchart.png"); err != nil {
	// 	panic(err)
	// }

	return c.Image()
}

func PlotObjectiveStats(data StatsJSON) image.Image {
	stats := data.AdvancedStats()

	labels := make([]string, 8)
	BerriesDunked := make(plotter.Values, 8)
	BerriesThrown := make(plotter.Values, 8)
	Snail := make(plotter.Values, 8)

	counter := 0
	for _, teamStats := range stats {
		for name, playerStats := range teamStats {
			labels[counter] = name
			// QueenKills = append(QueenKills, float64(playerStats["QueenKills"]))
			BerriesDunked[counter] = float64(playerStats["BerryDunks"])
			BerriesThrown[counter] = float64(playerStats["BerryThrows"])
			Snail[counter] = float64(playerStats["Snail"]) / 50
			counter++
		}
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Objective Stats"
	p.Title.Font.Size = vg.Points(20)
	p.Y.Label.Text = "Count"
	p.Y.Label.Font.Size = vg.Points(16)

	w := vg.Points(20)

	barsA, err := plotter.NewBarChart(BerriesDunked, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	barsB, err := plotter.NewBarChart(BerriesThrown, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)

	barsC, err := plotter.NewBarChart(Snail, w)
	if err != nil {
		panic(err)
	}
	barsC.LineStyle.Width = vg.Length(0)
	barsC.Color = plotutil.Color(2)
	barsC.Offset = w

	p.Add(barsA, barsB, barsC)
	p.Legend.Add("BerriesDunked", barsA)
	p.Legend.Add("BerriesThrown", barsB)
	p.Legend.Add("Snail/50", barsC)
	p.Legend.Top = true
	p.Legend.Font.Size = vg.Points(20)
	p.NominalX(labels...)
	p.X.Tick.Label.Font.Size = vg.Points(12)
	p.Y.Tick.Label.Font.Size = vg.Points(16)

	// Draw the plot to an in-memory image.
	img := image.NewRGBA(image.Rect(0, 0, 16*dpi, 9*dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	// if err := p.Save(16*vg.Inch, 8*vg.Inch, "barchart.png"); err != nil {
	// 	panic(err)
	// }

	return c.Image()
}
