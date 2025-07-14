# Qucanft - Go Astrological Calculations Library

Qucanft is a comprehensive Go library for astrological calculations, providing utilities for time and date conversions, coordinate system transformations, planetary positions, zodiac signs, aspects, and house systems.

## Features

### Core Modules

- **Time Utilities** (`pkg/time`)
  - Julian Day Number calculations
  - Sidereal time calculations  
  - Time system conversions
  - Delta T calculations

- **Coordinate Systems** (`pkg/coordinates`)
  - Equatorial coordinates
  - Ecliptic coordinates
  - Horizontal coordinates
  - Coordinate transformations
  - Angular calculations

- **Planetary Positions** (`pkg/planets`)
  - Planetary orbital calculations
  - Kepler's equation solver
  - Sun and planet positions
  - Multiple planet calculations

- **Zodiac Signs** (`pkg/zodiac`)
  - 12 zodiac signs with properties
  - Ecliptic to zodiac conversion
  - Sign compatibility analysis
  - Retrograde motion detection

- **Aspects** (`pkg/aspects`)
  - Major and minor aspects
  - Aspect calculations between planets
  - Aspect pattern detection (Grand Trine, T-Square, etc.)
  - Aspect strength analysis

- **House Systems** (`pkg/houses`)
  - Equal House system
  - Placidus system
  - Whole Sign system
  - Koch, Campanus, Regiomontanus systems
  - Planet-to-house assignments

## Installation

```bash
go get github.com/Qucanft/Qucanft
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/Qucanft/Qucanft/pkg/planets"
    timeutil "github.com/Qucanft/Qucanft/pkg/time"
    "github.com/Qucanft/Qucanft/pkg/zodiac"
)

func main() {
    // Create calculators
    pc := planets.NewPlanetaryCalculator()
    zc := zodiac.NewZodiacCalculator()
    tc := timeutil.NewTimeConverter()
    
    // Calculate for a specific date
    birthTime := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
    jd := tc.ToJulianDay(birthTime)
    
    // Get Sun position
    sunPos, err := pc.CalculatePosition("Sun", jd)
    if err != nil {
        panic(err)
    }
    
    // Convert to zodiac
    zodiacPos := zc.EclipticToZodiac(sunPos.Coordinates.Longitude)
    formatted := zc.FormatZodiacPosition(zodiacPos)
    
    fmt.Printf("Sun position: %s\n", formatted)
}
```

## Detailed Usage

### Time Calculations

```go
import timeutil "github.com/Qucanft/Qucanft/pkg/time"

tc := timeutil.NewTimeConverter()

// Convert time to Julian Day
jd := tc.ToJulianDay(time.Now())

// Calculate sidereal time
lst := tc.LocalSiderealTime(jd, longitude)

// Get Julian centuries since J2000
centuries := tc.JulianCenturies(jd)
```

### Coordinate Transformations

```go
import "github.com/Qucanft/Qucanft/pkg/coordinates"

ct := coordinates.NewCoordinateTransformer()

// Convert equatorial to ecliptic
eq := coordinates.EquatorialCoordinates{
    RightAscension: 180.0,
    Declination:    23.5,
    Distance:       1.0,
}
ec := ct.EquatorialToEcliptic(eq)

// Calculate angular separation
sep := ct.AngularSeparation(coord1, coord2)
```

### Planetary Positions

```go
import "github.com/Qucanft/Qucanft/pkg/planets"

pc := planets.NewPlanetaryCalculator()

// Calculate single planet position
marsPos, err := pc.CalculatePosition("Mars", jd)

// Calculate multiple planets
planetNames := []string{"Sun", "Moon", "Mercury", "Venus", "Mars"}
positions, err := pc.CalculateMultiplePositions(planetNames, jd)

// Calculate Sun position (optimized)
sunPos, err := pc.CalculateSunPosition(jd)
```

### Zodiac Signs

```go
import "github.com/Qucanft/Qucanft/pkg/zodiac"

zc := zodiac.NewZodiacCalculator()

// Convert ecliptic longitude to zodiac position
zodiacPos := zc.EclipticToZodiac(longitude)

// Format zodiac position
formatted := zc.FormatZodiacPosition(zodiacPos)

// Get sign compatibility
aries, _ := zc.GetSignByName("Aries")
leo, _ := zc.GetSignByName("Leo")
compatibility := zc.GetSignCompatibility(aries, leo)
```

### Aspects

```go
import "github.com/Qucanft/Qucanft/pkg/aspects"

ac := aspects.NewAspectCalculator()

// Calculate aspect between two planets
aspect := ac.CalculateAspect(planetPos1, planetPos2)

// Calculate all aspects
aspects := ac.CalculateAllAspects(positions)

// Get aspects by planet
sunAspects := ac.GetAspectsByPlanet(aspects, "Sun")

// Get aspects by nature
harmonious := ac.GetAspectsByNature(aspects, "Harmonious")
```

### House Systems

```go
import "github.com/Qucanft/Qucanft/pkg/houses"

hc := houses.NewHouseCalculator(houses.Equal)

// Calculate house cusps
cusps, err := hc.CalculateHouseCusps(ascendant, midheaven, latitude)

// Calculate complete houses
houses, err := hc.CalculateHouses(ascendant, midheaven, latitude)

// Add planets to houses
housesWithPlanets := hc.AddPlanetsToHouses(houses, positions)
```

## Available Data

### Planets
- Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto
- Orbital elements and symbols included

### Zodiac Signs
- All 12 signs with elements, qualities, and rulers
- Element: Fire, Earth, Air, Water
- Quality: Cardinal, Fixed, Mutable

### Aspects
- **Major**: Conjunction (0°), Sextile (60°), Square (90°), Trine (120°), Opposition (180°)
- **Minor**: Semisextile (30°), Semisquare (45°), Sesquiquadrate (135°), Quincunx (150°), Quintile (72°), Biquintile (144°)

### House Systems
- Equal House (30° each)
- Placidus (most common)
- Whole Sign
- Koch, Campanus, Regiomontanus

## Examples

See the `examples/` directory for complete examples:

- `basic_usage.go` - Simple planetary calculations
- `cmd/main.go` - Comprehensive demonstration

Run the main example:
```bash
go run cmd/main.go
```

## Testing

The library includes comprehensive tests for all modules:

```bash
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
```

## API Reference

### Core Types

- `timeutil.JulianDay` - Julian Day Number representation
- `coordinates.EquatorialCoordinates` - RA/Dec coordinates
- `coordinates.EclipticCoordinates` - Longitude/Latitude coordinates
- `planets.PlanetaryPosition` - Planet position at specific time
- `zodiac.ZodiacPosition` - Position within zodiac sign
- `aspects.Aspect` - Planetary aspect with strength
- `houses.House` - Astrological house with planets

### Key Functions

- `timeutil.ToJulianDay()` - Convert time to Julian Day
- `coordinates.EquatorialToEcliptic()` - Convert coordinates
- `planets.CalculatePosition()` - Calculate planet position
- `zodiac.EclipticToZodiac()` - Convert to zodiac position
- `aspects.CalculateAllAspects()` - Find all aspects
- `houses.CalculateHouses()` - Calculate house cusps

## Mathematical Foundations

The library implements standard astronomical algorithms:

- **Kepler's Equation** for orbital mechanics
- **Coordinate Transformations** using spherical trigonometry
- **Sidereal Time** calculations
- **Aspect Detection** with customizable orbs
- **House Division** systems

## Accuracy

This library provides reasonable accuracy for astrological purposes. For precise astronomical calculations, consider using specialized ephemeris data (JPL, Swiss Ephemeris).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Astronomical algorithms based on standard references
- Inspired by the original Python implementation
- Mathematical foundations from astronomical textbooks