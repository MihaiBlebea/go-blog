package page

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

type Position struct {
	Title       string
	StartDate   string
	EndDate     string
	Description string
}

type Home struct {
	*Page
	Positions []Position
}

func (h *Home) AddPositions(positions []Position) {
	h.Positions = positions
}

func (h *Home) Render(w io.Writer) error {
	err := h.Template.Execute(w, h)
	return err
}

func NewHome() *Home {
	page := New()
	positions, err := loadPositionsFromFile("./content/pages/positions.json")
	if err != nil {
		log.Panic(err)
	}

	home := &Home{
		page,
		positions,
	}

	home.AddTitle("Home")
	return home
}

func loadPositionsFromFile(filePath string) ([]Position, error) {
	posJson, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var positions []Position
	err = json.Unmarshal(posJson, &positions)
	if err != nil {
		return nil, err
	}

	return positions, nil
}
