package visualization

import (
	"image"
	"image/color"
	"math"

	"github.com/Qucanft/Qucanft/pkg/astrology"
)

// Drawing utility functions

func (ag *ArtGenerator) drawCircle(img *image.RGBA, centerX, centerY, radius int, c color.Color) {
	for x := centerX - radius; x <= centerX+radius; x++ {
		for y := centerY - radius; y <= centerY+radius; y++ {
			if x >= 0 && x < ag.config.Width && y >= 0 && y < ag.config.Height {
				dx := x - centerX
				dy := y - centerY
				if dx*dx+dy*dy <= radius*radius {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func (ag *ArtGenerator) drawSquare(img *image.RGBA, centerX, centerY, size int, c color.Color) {
	for x := centerX - size/2; x <= centerX+size/2; x++ {
		for y := centerY - size/2; y <= centerY+size/2; y++ {
			if x >= 0 && x < ag.config.Width && y >= 0 && y < ag.config.Height {
				img.Set(x, y, c)
			}
		}
	}
}

func (ag *ArtGenerator) drawTriangle(img *image.RGBA, centerX, centerY, size int, c color.Color) {
	// Draw a simple triangle
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			x := centerX - i/2 + j
			y := centerY - size/2 + i
			if x >= 0 && x < ag.config.Width && y >= 0 && y < ag.config.Height {
				img.Set(x, y, c)
			}
		}
	}
}

func (ag *ArtGenerator) drawOrganicShape(img *image.RGBA, centerX, centerY, size int, c color.Color) {
	// Create an organic, slightly irregular shape
	for angle := 0; angle < 360; angle += 5 {
		rad := float64(angle) * math.Pi / 180
		// Add some variation to radius for organic feel
		variation := 0.3 * math.Sin(float64(angle)*0.1) * float64(size)
		radius := float64(size) + variation
		
		x := centerX + int(radius*math.Cos(rad))
		y := centerY + int(radius*math.Sin(rad))
		
		ag.drawCircle(img, x, y, 2, c)
	}
}

func (ag *ArtGenerator) drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	// Bresenham's line algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	sy := 1
	
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	
	err := dx - dy
	x, y := x1, y1
	
	for {
		if x >= 0 && x < ag.config.Width && y >= 0 && y < ag.config.Height {
			img.Set(x, y, c)
		}
		
		if x == x2 && y == y2 {
			break
		}
		
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func (ag *ArtGenerator) drawCurvedLine(img *image.RGBA, angle1, angle2 float64, centerX, centerY int, c color.Color) {
	// Draw a curved line between two angles
	steps := 50
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		
		// Interpolate between angles
		angle := angle1 + t*(angle2-angle1)
		
		// Add curve by varying radius
		baseRadius := 100.0
		curveRadius := baseRadius + 20*math.Sin(t*math.Pi)
		
		rad := angle * math.Pi / 180
		x := centerX + int(curveRadius*math.Cos(rad))
		y := centerY + int(curveRadius*math.Sin(rad))
		
		ag.drawCircle(img, x, y, 2, c)
	}
}

// Color functions

func (ag *ArtGenerator) getPlanetColor(planet astrology.Planet) color.RGBA {
	switch ag.config.ColorScheme {
	case Cosmic:
		return ag.getCosmicPlanetColor(planet)
	case Earthy:
		return ag.getEarthyPlanetColor(planet)
	case Oceanic:
		return ag.getOceanicPlanetColor(planet)
	case Sunset:
		return ag.getSunsetPlanetColor(planet)
	default:
		return ag.getCosmicPlanetColor(planet)
	}
}

func (ag *ArtGenerator) getCosmicPlanetColor(planet astrology.Planet) color.RGBA {
	colors := map[astrology.Planet]color.RGBA{
		astrology.Sun:     {255, 215, 0, 255},   // Gold
		astrology.Moon:    {192, 192, 192, 255}, // Silver
		astrology.Mercury: {255, 165, 0, 255},   // Orange
		astrology.Venus:   {255, 20, 147, 255},  // Deep Pink
		astrology.Mars:    {255, 0, 0, 255},     // Red
		astrology.Jupiter: {138, 43, 226, 255},  // Blue Violet
		astrology.Saturn:  {128, 128, 128, 255}, // Gray
		astrology.Uranus:  {0, 255, 255, 255},   // Cyan
		astrology.Neptune: {0, 0, 255, 255},     // Blue
		astrology.Pluto:   {128, 0, 128, 255},   // Purple
	}
	
	if color, exists := colors[planet]; exists {
		return color
	}
	return color.RGBA{255, 255, 255, 255}
}

func (ag *ArtGenerator) getEarthyPlanetColor(planet astrology.Planet) color.RGBA {
	colors := map[astrology.Planet]color.RGBA{
		astrology.Sun:     {218, 165, 32, 255},  // Goldenrod
		astrology.Moon:    {245, 245, 220, 255}, // Beige
		astrology.Mercury: {210, 180, 140, 255}, // Tan
		astrology.Venus:   {154, 205, 50, 255},  // Yellow Green
		astrology.Mars:    {160, 82, 45, 255},   // Saddle Brown
		astrology.Jupiter: {107, 142, 35, 255},  // Olive Drab
		astrology.Saturn:  {105, 105, 105, 255}, // Dim Gray
		astrology.Uranus:  {95, 158, 160, 255},  // Cadet Blue
		astrology.Neptune: {72, 61, 139, 255},   // Dark Slate Blue
		astrology.Pluto:   {85, 107, 47, 255},   // Dark Olive Green
	}
	
	if color, exists := colors[planet]; exists {
		return color
	}
	return color.RGBA{139, 69, 19, 255}
}

func (ag *ArtGenerator) getOceanicPlanetColor(planet astrology.Planet) color.RGBA {
	colors := map[astrology.Planet]color.RGBA{
		astrology.Sun:     {255, 215, 0, 255},   // Gold
		astrology.Moon:    {175, 238, 238, 255}, // Pale Turquoise
		astrology.Mercury: {64, 224, 208, 255},  // Turquoise
		astrology.Venus:   {0, 206, 209, 255},   // Dark Turquoise
		astrology.Mars:    {70, 130, 180, 255},  // Steel Blue
		astrology.Jupiter: {25, 25, 112, 255},   // Midnight Blue
		astrology.Saturn:  {112, 128, 144, 255}, // Slate Gray
		astrology.Uranus:  {176, 224, 230, 255}, // Powder Blue
		astrology.Neptune: {0, 0, 139, 255},     // Dark Blue
		astrology.Pluto:   {72, 61, 139, 255},   // Dark Slate Blue
	}
	
	if color, exists := colors[planet]; exists {
		return color
	}
	return color.RGBA{0, 191, 255, 255}
}

func (ag *ArtGenerator) getSunsetPlanetColor(planet astrology.Planet) color.RGBA {
	colors := map[astrology.Planet]color.RGBA{
		astrology.Sun:     {255, 140, 0, 255},   // Dark Orange
		astrology.Moon:    {255, 160, 122, 255}, // Light Salmon
		astrology.Mercury: {255, 165, 0, 255},   // Orange
		astrology.Venus:   {255, 20, 147, 255},  // Deep Pink
		astrology.Mars:    {220, 20, 60, 255},   // Crimson
		astrology.Jupiter: {255, 69, 0, 255},    // Orange Red
		astrology.Saturn:  {128, 0, 128, 255},   // Purple
		astrology.Uranus:  {255, 105, 180, 255}, // Hot Pink
		astrology.Neptune: {138, 43, 226, 255},  // Blue Violet
		astrology.Pluto:   {75, 0, 130, 255},    // Indigo
	}
	
	if color, exists := colors[planet]; exists {
		return color
	}
	return color.RGBA{255, 69, 0, 255}
}

func (ag *ArtGenerator) getZodiacColor(sign astrology.ZodiacSign) color.RGBA {
	colors := map[astrology.ZodiacSign]color.RGBA{
		astrology.Aries:       {255, 0, 0, 200},     // Red
		astrology.Taurus:      {0, 128, 0, 200},     // Green
		astrology.Gemini:      {255, 255, 0, 200},   // Yellow
		astrology.Cancer:      {192, 192, 192, 200}, // Silver
		astrology.Leo:         {255, 215, 0, 200},   // Gold
		astrology.Virgo:       {128, 128, 0, 200},   // Olive
		astrology.Libra:       {255, 20, 147, 200},  // Deep Pink
		astrology.Scorpio:     {128, 0, 0, 200},     // Maroon
		astrology.Sagittarius: {128, 0, 128, 200},   // Purple
		astrology.Capricorn:   {0, 0, 128, 200},     // Navy
		astrology.Aquarius:    {0, 255, 255, 200},   // Cyan
		astrology.Pisces:      {0, 128, 128, 200},   // Teal
	}
	
	if color, exists := colors[sign]; exists {
		return color
	}
	return color.RGBA{255, 255, 255, 200}
}

func (ag *ArtGenerator) getAspectColor(aspect astrology.AspectType) color.RGBA {
	switch aspect {
	case astrology.Conjunction:
		return color.RGBA{255, 255, 255, 150} // White
	case astrology.Sextile:
		return color.RGBA{0, 255, 0, 150} // Green
	case astrology.Square:
		return color.RGBA{255, 0, 0, 150} // Red
	case astrology.Trine:
		return color.RGBA{0, 0, 255, 150} // Blue
	case astrology.Opposition:
		return color.RGBA{255, 165, 0, 150} // Orange
	default:
		return color.RGBA{128, 128, 128, 150} // Gray
	}
}

func (ag *ArtGenerator) getPlanetSize(planet astrology.Planet) int {
	sizes := map[astrology.Planet]int{
		astrology.Sun:     12,
		astrology.Moon:    10,
		astrology.Mercury: 6,
		astrology.Venus:   8,
		astrology.Mars:    7,
		astrology.Jupiter: 15,
		astrology.Saturn:  13,
		astrology.Uranus:  9,
		astrology.Neptune: 9,
		astrology.Pluto:   5,
	}
	
	if size, exists := sizes[planet]; exists {
		return size
	}
	return 8
}

// Utility functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}