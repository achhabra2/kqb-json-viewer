package stats

import (
	"path/filepath"
	"testing"
)

func TestListStats(t *testing.T) {
	files := ListStatFiles()
	for _, name := range files {
		if filepath.Ext(name) != ".json" {
			t.Error("Directory contains non json files")
			t.Log(files)
		}
	}
}

func TestAdvancedStats(t *testing.T) {
	names := ListStatFiles()
	data := ReadJson(names[0])
	output := data.AdvancedStats()
	t.Log(output)
}
