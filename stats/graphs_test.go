package stats

import (
	"image/png"
	"os"
	"testing"
)

func TestGraph(t *testing.T) {
	data := ReadJson("../fixtures/sample.json")
	img := PlotStats(data)

	// // Save the image.
	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}
