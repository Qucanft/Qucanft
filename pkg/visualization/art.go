package visualization

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/Qucanft/Qucanft/pkg/astrology"
)

// ArtConfig holds configuration for artistic visualization
type ArtConfig struct {
	Width        int
	Height       int
	Background   color.Color
	Style        ArtStyle
	ShowLabels   bool
	ShowAspects  bool
	ShowHouses   bool
	ColorScheme  ColorScheme
}

// ArtStyle defines different artistic styles
type ArtStyle int

const (
	Mandala ArtStyle = iota
	Geometric
	Organic
	Minimalist
)

// ColorScheme defines color palettes
type ColorScheme int

const (
	Cosmic ColorScheme = iota
	Earthy
	Oceanic
	Sunset
)

// ArtGenerator creates artistic visualizations from astrological data
type ArtGenerator struct {
	config ArtConfig
}

// NewArtGenerator creates a new art generator with default configuration
func NewArtGenerator(config ArtConfig) *ArtGenerator {
	return &ArtGenerator{config: config}
}

// GenerateVisualization creates an artistic visualization from a chart
func (ag *ArtGenerator) GenerateVisualization(chart *astrology.Chart) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, ag.config.Width, ag.config.Height))
	
	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{ag.config.Background}, image.Point{}, draw.Src)

	switch ag.config.Style {
	case Mandala:
		ag.drawMandala(img, chart)
	case Geometric:
		ag.drawGeometric(img, chart)
	case Organic:
		ag.drawOrganic(img, chart)
	case Minimalist:
		ag.drawMinimalist(img, chart)
	}

	return img, nil
}

// drawMandala creates a mandala-style visualization
func (ag *ArtGenerator) drawMandala(img *image.RGBA, chart *astrology.Chart) {
	centerX := ag.config.Width / 2
	centerY := ag.config.Height / 2
	radius := float64(min(ag.config.Width, ag.config.Height)) / 2 * 0.8

	// Draw zodiac circle
	ag.drawZodiacCircle(img, centerX, centerY, radius)

	// Draw planets
	for _, planet := range chart.Planets {
		ag.drawPlanet(img, planet, centerX, centerY, radius*0.7)
	}

	// Draw aspects if enabled
	if ag.config.ShowAspects {
		for _, aspect := range chart.Aspects {
			ag.drawAspect(img, aspect, chart, centerX, centerY, radius*0.7)
		}
	}
}

// drawGeometric creates a geometric visualization
func (ag *ArtGenerator) drawGeometric(img *image.RGBA, chart *astrology.Chart) {
	centerX := ag.config.Width / 2
	centerY := ag.config.Height / 2
	
	// Draw concentric shapes based on elements
	elements := map[string]int{"Fire": 0, "Earth": 0, "Air": 0, "Water": 0}
	for _, planet := range chart.Planets {
		elements[planet.GetElementalEnergy()]++
	}

	// Create geometric patterns based on elemental distribution
	ag.drawElementalGeometry(img, elements, centerX, centerY)

	// Draw planet positions as geometric shapes
	for _, planet := range chart.Planets {
		ag.drawGeometricPlanet(img, planet, centerX, centerY)
	}
}

// drawOrganic creates an organic, flowing visualization
func (ag *ArtGenerator) drawOrganic(img *image.RGBA, chart *astrology.Chart) {
	centerX := ag.config.Width / 2
	centerY := ag.config.Height / 2

	// Draw flowing energy lines based on aspects
	for _, aspect := range chart.Aspects {
		ag.drawEnergyFlow(img, aspect, chart, centerX, centerY)
	}

	// Draw planets as organic shapes
	for _, planet := range chart.Planets {
		ag.drawOrganicPlanet(img, planet, centerX, centerY)
	}
}

// drawMinimalist creates a clean, minimalist visualization
func (ag *ArtGenerator) drawMinimalist(img *image.RGBA, chart *astrology.Chart) {
	centerX := ag.config.Width / 2
	centerY := ag.config.Height / 2

	// Draw simple circle for zodiac
	ag.drawSimpleCircle(img, centerX, centerY, 200)

	// Draw planets as simple dots
	for _, planet := range chart.Planets {
		ag.drawMinimalPlanet(img, planet, centerX, centerY)
	}

	// Draw aspects as simple lines
	if ag.config.ShowAspects {
		for _, aspect := range chart.Aspects {
			ag.drawSimpleAspect(img, aspect, chart, centerX, centerY)
		}
	}
}

// Helper functions for drawing specific elements

func (ag *ArtGenerator) drawZodiacCircle(img *image.RGBA, centerX, centerY int, radius float64) {
	for i := 0; i < 12; i++ {
		angle := float64(i) * 30 * math.Pi / 180
		x := centerX + int(radius*math.Cos(angle))
		y := centerY + int(radius*math.Sin(angle))
		
		// Draw zodiac sign markers
		ag.drawCircle(img, x, y, 5, ag.getZodiacColor(astrology.ZodiacSign(i)))
	}
}

func (ag *ArtGenerator) drawPlanet(img *image.RGBA, planet astrology.PlanetPosition, centerX, centerY int, radius float64) {
	angle := planet.Degree * math.Pi / 180
	x := centerX + int(radius*math.Cos(angle))
	y := centerY + int(radius*math.Sin(angle))
	
	planetColor := ag.getPlanetColor(planet.Planet)
	size := ag.getPlanetSize(planet.Planet)
	
	ag.drawCircle(img, x, y, size, planetColor)
	
	// Draw retrograde indicator
	if planet.Retrograde {
		ag.drawCircle(img, x, y, size+2, color.RGBA{255, 255, 255, 100})
	}
}

func (ag *ArtGenerator) drawAspect(img *image.RGBA, aspect astrology.Aspect, chart *astrology.Chart, centerX, centerY int, radius float64) {
	planet1Pos, _ := chart.GetPlanetPosition(aspect.Planet1)
	planet2Pos, _ := chart.GetPlanetPosition(aspect.Planet2)
	
	angle1 := planet1Pos.Degree * math.Pi / 180
	angle2 := planet2Pos.Degree * math.Pi / 180
	
	x1 := centerX + int(radius*math.Cos(angle1))
	y1 := centerY + int(radius*math.Sin(angle1))
	x2 := centerX + int(radius*math.Cos(angle2))
	y2 := centerY + int(radius*math.Sin(angle2))
	
	aspectColor := ag.getAspectColor(aspect.Type)
	ag.drawLine(img, x1, y1, x2, y2, aspectColor)
}

func (ag *ArtGenerator) drawElementalGeometry(img *image.RGBA, elements map[string]int, centerX, centerY int) {
	colors := map[string]color.RGBA{
		"Fire":  {255, 100, 100, 200},
		"Earth": {139, 69, 19, 200},
		"Air":   {173, 216, 230, 200},
		"Water": {100, 149, 237, 200},
	}
	
	i := 0
	for element, count := range elements {
		if count > 0 {
			radius := 50 + count*20
			ag.drawCircle(img, centerX+i*30, centerY+i*30, radius, colors[element])
			i++
		}
	}
}

func (ag *ArtGenerator) drawGeometricPlanet(img *image.RGBA, planet astrology.PlanetPosition, centerX, centerY int) {
	angle := planet.Degree * math.Pi / 180
	radius := 100.0
	x := centerX + int(radius*math.Cos(angle))
	y := centerY + int(radius*math.Sin(angle))
	
	// Draw different shapes for different planets
	planetColor := ag.getPlanetColor(planet.Planet)
	switch planet.Planet {
	case astrology.Sun:
		ag.drawSquare(img, x, y, 10, planetColor)
	case astrology.Moon:
		ag.drawCircle(img, x, y, 8, planetColor)
	default:
		ag.drawTriangle(img, x, y, 6, planetColor)
	}
}

func (ag *ArtGenerator) drawEnergyFlow(img *image.RGBA, aspect astrology.Aspect, chart *astrology.Chart, centerX, centerY int) {
	// Create flowing, organic lines for aspects
	planet1Pos, _ := chart.GetPlanetPosition(aspect.Planet1)
	planet2Pos, _ := chart.GetPlanetPosition(aspect.Planet2)
	
	// Draw curved line instead of straight line
	ag.drawCurvedLine(img, planet1Pos.Degree, planet2Pos.Degree, centerX, centerY, ag.getAspectColor(aspect.Type))
}

func (ag *ArtGenerator) drawOrganicPlanet(img *image.RGBA, planet astrology.PlanetPosition, centerX, centerY int) {
	angle := planet.Degree * math.Pi / 180
	radius := 120.0
	x := centerX + int(radius*math.Cos(angle))
	y := centerY + int(radius*math.Sin(angle))
	
	// Draw organic, blob-like shapes
	planetColor := ag.getPlanetColor(planet.Planet)
	ag.drawOrganicShape(img, x, y, ag.getPlanetSize(planet.Planet), planetColor)
}

func (ag *ArtGenerator) drawMinimalPlanet(img *image.RGBA, planet astrology.PlanetPosition, centerX, centerY int) {
	angle := planet.Degree * math.Pi / 180
	radius := 150.0
	x := centerX + int(radius*math.Cos(angle))
	y := centerY + int(radius*math.Sin(angle))
	
	ag.drawCircle(img, x, y, 3, color.RGBA{0, 0, 0, 255})
}

func (ag *ArtGenerator) drawSimpleCircle(img *image.RGBA, centerX, centerY, radius int) {
	for angle := 0; angle < 360; angle++ {
		rad := float64(angle) * math.Pi / 180
		x := centerX + int(float64(radius)*math.Cos(rad))
		y := centerY + int(float64(radius)*math.Sin(rad))
		
		if x >= 0 && x < ag.config.Width && y >= 0 && y < ag.config.Height {
			img.Set(x, y, color.RGBA{100, 100, 100, 255})
		}
	}
}

func (ag *ArtGenerator) drawSimpleAspect(img *image.RGBA, aspect astrology.Aspect, chart *astrology.Chart, centerX, centerY int) {
	planet1Pos, _ := chart.GetPlanetPosition(aspect.Planet1)
	planet2Pos, _ := chart.GetPlanetPosition(aspect.Planet2)
	
	angle1 := planet1Pos.Degree * math.Pi / 180
	angle2 := planet2Pos.Degree * math.Pi / 180
	
	x1 := centerX + int(150*math.Cos(angle1))
	y1 := centerY + int(150*math.Sin(angle1))
	x2 := centerX + int(150*math.Cos(angle2))
	y2 := centerY + int(150*math.Sin(angle2))
	
	ag.drawLine(img, x1, y1, x2, y2, color.RGBA{200, 200, 200, 100})
}

// SaveImage saves the generated image to a file
func (ag *ArtGenerator) SaveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return png.Encode(file, img)
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() ArtConfig {
	return ArtConfig{
		Width:       800,
		Height:      600,
		Background:  color.RGBA{20, 20, 40, 255},
		Style:       Mandala,
		ShowLabels:  true,
		ShowAspects: true,
		ShowHouses:  true,
		ColorScheme: Cosmic,
	}
}

// Color helper functions will be implemented in the next file