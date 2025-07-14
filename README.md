# Qucanft Astronomical Library

A comprehensive Go library for astronomical calculations and astrological computations. This library provides modular and reusable components for planetary positioning, coordinate transformations, and astrological chart calculations.

## Features

- **Time Calculations**: Julian Date conversions and time utilities
- **Coordinate Systems**: Equatorial and ecliptic coordinate transformations
- **Planetary Positions**: Calculate positions for Sun, Moon, and planets
- **Zodiac Signs**: Determine zodiac signs and their properties
- **Astrological Aspects**: Calculate aspects between planets
- **Chart Generation**: Create complete astrological charts

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
    
    "github.com/Qucanft/Qucanft/pkg/astro"
    "github.com/Qucanft/Qucanft/pkg/planets"
)

func main() {
    // Create a chart for a specific date and location
    dateTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
    location := astro.Location{
        Latitude:  40.7128,  // New York City
        Longitude: -74.0060,
        Timezone:  "America/New_York",
    }
    
    chart := astro.NewChart(dateTime, location)
    
    // Get Sun position
    sunPos := chart.GetPlanetPosition(planets.Sun)
    sunSign := chart.GetPlanetSign(planets.Sun)
    
    fmt.Printf("Sun: %s at %s\n", sunSign.FormatPosition(), sunSign.Element)
    
    // Get major aspects
    aspects := chart.GetMajorAspects()
    fmt.Printf("Found %d major aspects\n", len(aspects))
}
```

## Package Structure

### Core Packages

- **`pkg/time`**: Julian Date calculations and time utilities
- **`pkg/coordinates`**: Coordinate system transformations
- **`pkg/planets`**: Planetary position calculations
- **`pkg/zodiac`**: Zodiac sign calculations and properties
- **`pkg/aspects`**: Astrological aspect calculations
- **`pkg/astro`**: Main chart generation and analysis

### Time Package

```go
import "github.com/Qucanft/Qucanft/pkg/time"

jd := time.ToJulianDate(time.Now())
daysSinceJ2000 := jd.DaysSinceJ2000()
```

### Coordinates Package

```go
import "github.com/Qucanft/Qucanft/pkg/coordinates"

// Convert between coordinate systems
obliquity := coordinates.GetObliquity(jd)
ecliptic := coordinates.EclipticCoordinates{Longitude: 45, Latitude: 0}
equatorial := ecliptic.ToEquatorial(obliquity)
```

### Planets Package

```go
import "github.com/Qucanft/Qucanft/pkg/planets"

// Calculate planetary positions
sunPos := planets.CalculatePosition(planets.Sun, jd)
fmt.Printf("Sun longitude: %.2f°\n", sunPos.Longitude)
```

### Zodiac Package

```go
import "github.com/Qucanft/Qucanft/pkg/zodiac"

// Get zodiac sign information
signInfo := zodiac.GetSignFromLongitude(sunPos.Longitude)
fmt.Printf("Sign: %s (%s, %s)\n", signInfo.Sign, signInfo.Element, signInfo.Quality)
```

### Aspects Package

```go
import "github.com/Qucanft/Qucanft/pkg/aspects"

// Calculate aspects between planets
aspects := aspects.CalculateAspects(planetPositions)
majorAspects := aspects.GetMajorAspects(aspects)
```

## Examples

See `cmd/example/main.go` for a complete example that demonstrates:

- Chart creation
- Planetary position calculations
- Zodiac sign determination
- Aspect calculations
- Coordinate transformations
- Elemental and quality balance analysis

Run the example:

```bash
go run cmd/example/main.go
```

## Testing

Run the test suite:

```bash
go test ./...
```

## Features Implemented

### ✅ Time Calculations
- Julian Date conversion
- J2000.0 epoch handling
- Time utilities for astronomical calculations

### ✅ Coordinate Systems
- Equatorial coordinates (Right Ascension, Declination)
- Ecliptic coordinates (Longitude, Latitude)
- Coordinate transformations
- Obliquity calculations

### ✅ Planetary Calculations
- Sun position (simplified solar ephemeris)
- Moon position (simplified lunar ephemeris)
- Planetary positions for Mercury through Pluto
- Orbital elements and Kepler's equation solving

### ✅ Zodiac System
- 12 zodiac signs with properties
- Element classification (Fire, Earth, Air, Water)
- Quality classification (Cardinal, Fixed, Mutable)
- Traditional rulerships
- Sign compatibility calculations

### ✅ Astrological Aspects
- Major aspects (conjunction, sextile, square, trine, opposition)
- Minor aspects (quincunx, semisextile, semisquare, sesquiquadrate)
- Orb calculations
- Applying/separating determination
- Aspect filtering and analysis

### ✅ Chart Generation
- Complete astrological chart creation
- Planetary position analysis
- Aspect calculations
- Elemental and quality balance
- Lunar phase calculations

## Architecture

The library is designed with modularity in mind:

- **Separation of Concerns**: Each package handles a specific domain
- **Reusable Components**: Functions can be used independently
- **Type Safety**: Strong typing for angles, coordinates, and astronomical objects
- **Extensible**: Easy to add new features or calculation methods

## Accuracy Note

This library uses simplified algorithms suitable for astrological purposes. For high-precision astronomical calculations, consider using more sophisticated ephemeris data and algorithms.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.