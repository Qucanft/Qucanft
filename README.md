# Astrological Art Generator 🌟

A powerful Go application that generates beautiful artistic visualizations based on astrological data. Create stunning visual representations of planetary positions, zodiac signs, aspects, and house calculations.

## Features ✨

- **Multiple Art Styles**: Choose from Mandala, Geometric, Organic, and Minimalist styles
- **Color Schemes**: Cosmic, Earthy, Oceanic, and Sunset palettes
- **Astrological Data**: Full support for planets, zodiac signs, aspects, and houses
- **Customizable Output**: Configure size, colors, and visual elements
- **Command-Line Interface**: Easy to use with extensive customization options

## Installation 🚀

```bash
# Clone the repository
git clone https://github.com/Qucanft/Qucanft.git
cd Qucanft

# Build the application
go build -o astroart cmd/astroart/main.go
```

## Usage 📖

### Basic Usage

```bash
# Generate a basic mandala-style chart
./astroart

# Generate with custom style and colors
./astroart -style geometric -colors sunset -output my-chart.png

# Generate with custom size
./astroart -width 1200 -height 800 -style organic

# Generate for specific date and time
./astroart -datetime "2023-12-25 12:00:00" -style minimalist
```

### Command-Line Options

- `-width` - Width of the generated image (default: 800)
- `-height` - Height of the generated image (default: 600)
- `-style` - Art style: mandala, geometric, organic, minimalist (default: mandala)
- `-colors` - Color scheme: cosmic, earthy, oceanic, sunset (default: cosmic)
- `-output` - Output filename (default: astro-art.png)
- `-aspects` - Show aspects in visualization (default: true)
- `-houses` - Show houses in visualization (default: true)
- `-labels` - Show labels in visualization (default: true)
- `-bg` - Background color as R,G,B (default: 20,20,40)
- `-datetime` - Date and time for chart generation (format: YYYY-MM-DD HH:MM:SS)

## Art Styles 🎨

### Mandala
Traditional circular astrological chart with zodiac wheel, planetary positions, and aspect lines.

### Geometric
Modern geometric interpretation using shapes and patterns based on elemental distribution.

### Organic
Flowing, natural forms that represent astrological energies through organic shapes and curves.

### Minimalist
Clean, simple design focusing on essential elements with minimal visual noise.

## Color Schemes 🌈

- **Cosmic**: Deep space colors with gold, silver, and vibrant planetary hues
- **Earthy**: Natural earth tones with browns, greens, and warm colors
- **Oceanic**: Ocean-inspired blues, teals, and aquatic colors
- **Sunset**: Warm sunset colors with oranges, pinks, and purples

## Astrological Data 🔮

The application generates comprehensive astrological charts including:

- **Planets**: Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto
- **Zodiac Signs**: All 12 signs with elemental and modal energies
- **Aspects**: Conjunction, Sextile, Square, Trine, Opposition
- **Houses**: Full 12-house system
- **Additional Features**: Retrograde motion, aspect orbs, elemental distribution

## Example Commands 💡

```bash
# Create a cosmic mandala
./astroart -style mandala -colors cosmic -output cosmic-mandala.png

# Generate geometric art with sunset colors
./astroart -style geometric -colors sunset -width 1000 -height 1000

# Create organic art for a specific date
./astroart -style organic -datetime "2024-01-01 00:00:00" -colors oceanic

# Minimalist chart with custom background
./astroart -style minimalist -bg "10,10,20" -colors earthy
```

## Project Structure 📁

```
pkg/
├── astrology/       # Astrological data structures and logic
│   ├── types.go     # Core astrological types
│   └── generator.go # Chart generation logic
└── visualization/   # Artistic visualization engine
    ├── art.go       # Main art generation logic
    └── drawing.go   # Drawing utilities and color schemes

cmd/
└── astroart/        # Command-line application
    └── main.go      # Main application entry point
```

## Contributing 🤝

Contributions are welcome! Please feel free to submit pull requests or open issues.

## License 📄

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.