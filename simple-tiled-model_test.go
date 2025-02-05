package wfc

import (
	"encoding/json"
	// "fmt"
	"image"
	"io/ioutil"
	"testing"

	"github.com/alexschoenwitz/wfc/internal/testutils"
)

// Parsed data supplied by user
type RawData struct {
	Path      string        `json:"path"`      // Path to tiles
	TileSize  int           `json:"tileSize"`  // Default to 16
	Tiles     []RawTile     `json:"tiles"`     //
	Neighbors []RawNeighbor `json:"neighbors"` //
}

// Raw information on a tile
type RawTile struct {
	Name     TileName `json:"name"`     // Name used to identify the tile
	Symmetry string   `json:"symmetry"` // Default to ""
	Weight   float64  `json:"weight"`   // Default to 1
}

// Information on which tiles can be neighbors
type RawNeighbor struct {
	Left     TileName `json:"left"`     // Matches Tile.Name
	LeftNum  int      `json:"leftNum"`  // Default to 0
	Right    TileName `json:"right"`    // Matches Tile.Name
	RightNum int      `json:"rightNum"` // Default to 0
}

func initiateData(dataFileName string) SimpleTiledData {
	// Load data file
	dataFile, err := ioutil.ReadFile("internal/input/" + dataFileName)
	if err != nil {
		panic(err)
	}

	// Parse rawData file
	var rawData RawData
	if err := json.Unmarshal(dataFile, &rawData); err != nil {
		panic(err)
	}

	// Marshal into data settings struct
	tiles := make([]Tile, len(rawData.Tiles))
	for i, rt := range rawData.Tiles {
		imgs := make([]image.Image, 0)
		img, err := testutils.LoadImage("internal/input/" + rawData.Path + string(rt.Name) + ".png")
		if err != nil {
			panic(err)
		}
		imgs = append(imgs, img)
		weight := rt.Weight
		if weight == 0 {
			weight = 1
		}
		tiles[i] = Tile{Name: rt.Name, Symmetry: rt.Symmetry, Weight: weight, Variants: imgs}
	}
	neighbors := make([]Neighbor, len(rawData.Neighbors))
	for i, rn := range rawData.Neighbors {
		neighbors[i] = Neighbor(rn)
	}
	return SimpleTiledData{TileSize: rawData.TileSize, Tiles: tiles, Neighbors: neighbors}
}

func simpleTiledTest(t *testing.T, dataFileName, snapshotFileName string, iterations int) {
	// Set test parameters
	periodic := false
	width := 20
	height := 20
	seed := int64(42)
	data := initiateData(dataFileName)

	// Generate output image
	var outputImg image.Image
	success, finished := false, false
	model := NewSimpleTiledModel(data, width, height, periodic)
	model.SetSeed(seed)
	if iterations == -1 {
		outputImg, success = model.Generate()
		if !success {
			t.Log("Failed to generate image on the first try.")
			t.FailNow()
		}
	} else {
		outputImg, finished, _ = model.Iterate(iterations)
		if finished {
			t.Log("Test for incomplete state actually finished.")
			t.FailNow()
		}
	}

	// Save output
	// err := testutils.SaveImage("internal/snapshots/"+snapshotFileName, outputImg)
	// if err != nil {
	// 	panic(err)
	// }

	// Test that files match
	snapshotImg, err := testutils.LoadImage("internal/snapshots/" + snapshotFileName)
	if err != nil {
		panic(err)
	}
	areEqual := testutils.CompareImages(outputImg, snapshotImg)
	if !areEqual {
		t.Log("Output image is not the same as the snapshot image.")
		t.FailNow()
	}
}

func TestSimpleTiledGenerationCompletes(t *testing.T) {
	simpleTiledTest(t, "castle_data.json", "castle.png", -1)
}

func TestSimpleTiledIterationIncomplete(t *testing.T) {
	simpleTiledTest(t, "castle_data.json", "castle_incomplete.png", 5)
}
