#!/usr/bin/env python3
"""
Visualization demo for Qucanft library without external dependencies.

This demo shows how the library can generate visualizations using basic matplotlib.
"""

import matplotlib.pyplot as plt
import numpy as np
import math
import json

# Create a simple visualization demo
def create_demo_natal_chart():
    """Create a demo natal chart visualization"""
    fig, ax = plt.subplots(figsize=(10, 10))
    ax.set_xlim(-1.5, 1.5)
    ax.set_ylim(-1.5, 1.5)
    ax.set_aspect('equal')
    ax.axis('off')
    
    # Draw outer circle (zodiac wheel)
    outer_circle = plt.Circle((0, 0), 1.2, fill=False, color='black', linewidth=2)
    ax.add_patch(outer_circle)
    
    # Draw inner circle
    inner_circle = plt.Circle((0, 0), 0.8, fill=False, color='gray', linewidth=1)
    ax.add_patch(inner_circle)
    
    # Zodiac signs and symbols
    zodiac_data = [
        ('Aries', '♈', '#FF6B6B'),
        ('Taurus', '♉', '#4ECDC4'),
        ('Gemini', '♊', '#FFE66D'),
        ('Cancer', '♋', '#95E1D3'),
        ('Leo', '♌', '#F8B500'),
        ('Virgo', '♍', '#C7CEEA'),
        ('Libra', '♎', '#FFB6C1'),
        ('Scorpio', '♏', '#8B0000'),
        ('Sagittarius', '♐', '#9B59B6'),
        ('Capricorn', '♑', '#2C3E50'),
        ('Aquarius', '♒', '#3498DB'),
        ('Pisces', '♓', '#1ABC9C')
    ]
    
    # Draw zodiac signs around the circle
    for i, (sign_name, symbol, color) in enumerate(zodiac_data):
        angle = (i * 30 - 90) * math.pi / 180  # Start from Aries at top
        x = 1.35 * math.cos(angle)
        y = 1.35 * math.sin(angle)
        ax.text(x, y, symbol, ha='center', va='center', fontsize=14, 
               weight='bold', color=color)
        
        # Sign name
        x_name = 1.05 * math.cos(angle)
        y_name = 1.05 * math.sin(angle)
        ax.text(x_name, y_name, sign_name[:3], ha='center', va='center', 
               fontsize=8, color='black')
    
    # Draw house cusps (Equal House system, Ascendant at 75° - Gemini)
    ascendant = 75.0
    for house_num in range(1, 13):
        cusp_deg = (ascendant + (house_num - 1) * 30) % 360
        angle = (cusp_deg - 90) * math.pi / 180
        x1 = 0.8 * math.cos(angle)
        y1 = 0.8 * math.sin(angle)
        x2 = 1.2 * math.cos(angle)
        y2 = 1.2 * math.sin(angle)
        ax.plot([x1, x2], [y1, y2], 'gray', alpha=0.7, linewidth=1)
        
        # House number
        x_text = 0.9 * math.cos(angle)
        y_text = 0.9 * math.sin(angle)
        ax.text(x_text, y_text, str(house_num), ha='center', va='center', 
               fontsize=10, color='gray', weight='bold')
    
    # Sample planetary positions (degrees)
    planets = [
        ('Sun', '☉', 280.5, '#FFD700'),
        ('Moon', '☽', 45.2, '#C0C0C0'),
        ('Mercury', '☿', 290.1, '#FFA500'),
        ('Venus', '♀', 315.8, '#FF69B4'),
        ('Mars', '♂', 120.3, '#FF4500'),
        ('Jupiter', '♃', 30.7, '#1E90FF'),
        ('Saturn', '♄', 350.2, '#8B4513')
    ]
    
    # Draw planets
    for planet_name, symbol, longitude, color in planets:
        angle = (longitude - 90) * math.pi / 180
        x = 1.0 * math.cos(angle)
        y = 1.0 * math.sin(angle)
        
        # Planet circle
        circle = plt.Circle((x, y), 0.08, color=color, alpha=0.8, 
                          edgecolor='black', linewidth=2)
        ax.add_patch(circle)
        
        # Planet symbol
        ax.text(x, y, symbol, ha='center', va='center', fontsize=12, 
               color='white', weight='bold')
    
    # Draw some sample aspects
    aspects = [
        (280.5, 45.2, 'blue', 'Trine'),    # Sun trine Moon
        (290.1, 120.3, 'red', 'Square'),   # Mercury square Mars
        (315.8, 30.7, 'blue', 'Trine'),    # Venus trine Jupiter
    ]
    
    for pos1, pos2, color, aspect_name in aspects:
        angle1 = (pos1 - 90) * math.pi / 180
        angle2 = (pos2 - 90) * math.pi / 180
        
        x1 = 1.0 * math.cos(angle1)
        y1 = 1.0 * math.sin(angle1)
        x2 = 1.0 * math.cos(angle2)
        y2 = 1.0 * math.sin(angle2)
        
        ax.plot([x1, x2], [y1, y2], color=color, alpha=0.6, linewidth=2)
    
    plt.title('Qucanft Demo - Natal Chart\nJanuary 1, 2024 - New York, NY', 
              fontsize=16, weight='bold', pad=20)
    
    # Add legend
    legend_elements = [
        plt.Line2D([0], [0], color='blue', lw=2, alpha=0.6, label='Harmonious Aspects'),
        plt.Line2D([0], [0], color='red', lw=2, alpha=0.6, label='Challenging Aspects'),
        plt.Line2D([0], [0], color='gray', lw=1, alpha=0.7, label='House Cusps')
    ]
    ax.legend(handles=legend_elements, loc='upper left', bbox_to_anchor=(1.02, 1))
    
    plt.tight_layout()
    return fig

def create_demo_planetary_distribution():
    """Create a demo planetary distribution chart"""
    fig, ax = plt.subplots(figsize=(12, 8))
    
    # Sample data: planets in zodiac signs
    zodiac_signs = ['Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo',
                   'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces']
    
    # Sample distribution (planets in signs)
    planet_counts = [0, 2, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0]  # Example distribution
    planet_lists = [
        [],
        ['Moon', 'Jupiter'],
        ['Mercury'],
        ['Mars'],
        ['Sun'],
        [],
        [],
        [],
        [],
        ['Venus'],
        ['Saturn'],
        []
    ]
    
    colors = ['#FF6B6B', '#4ECDC4', '#FFE66D', '#95E1D3', '#F8B500', '#C7CEEA',
             '#FFB6C1', '#8B0000', '#9B59B6', '#2C3E50', '#3498DB', '#1ABC9C']
    
    # Create bar chart
    bars = ax.bar(zodiac_signs, planet_counts, color=colors, alpha=0.8, 
                  edgecolor='black', linewidth=1)
    
    # Add planet names on bars
    for i, (count, planets) in enumerate(zip(planet_counts, planet_lists)):
        if count > 0:
            planet_text = ', '.join(planets)
            ax.text(i, count + 0.1, planet_text, ha='center', va='bottom', 
                   fontsize=10, weight='bold', rotation=0)
    
    ax.set_xlabel('Zodiac Signs', fontsize=14, weight='bold')
    ax.set_ylabel('Number of Planets', fontsize=14, weight='bold')
    ax.set_title('Qucanft Demo - Planetary Distribution\nJanuary 1, 2024 - New York, NY', 
                 fontsize=16, weight='bold', pad=20)
    ax.set_ylim(0, max(planet_counts) + 1)
    
    plt.xticks(rotation=45, ha='right')
    plt.grid(axis='y', alpha=0.3)
    plt.tight_layout()
    
    return fig

def create_demo_aspects_summary():
    """Create a demo aspects summary chart"""
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))
    
    # Sample aspects data
    aspect_types = ['Conjunction', 'Sextile', 'Square', 'Trine', 'Opposition']
    aspect_counts = [2, 1, 2, 2, 1]
    
    nature_types = ['Harmonious', 'Challenging', 'Neutral']
    nature_counts = [3, 3, 2]
    
    # Chart 1: Aspects by type
    colors1 = ['#FFD700', '#87CEEB', '#FF6B6B', '#90EE90', '#DDA0DD']
    wedges, texts, autotexts = ax1.pie(aspect_counts, labels=aspect_types,
                                      autopct='%1.1f%%', colors=colors1, 
                                      startangle=90)
    ax1.set_title('Aspects by Type', fontsize=14, weight='bold')
    
    # Chart 2: Aspects by nature
    colors2 = ['lightblue', 'lightcoral', 'lightgray']
    wedges2, texts2, autotexts2 = ax2.pie(nature_counts, labels=nature_types,
                                         autopct='%1.1f%%', colors=colors2, 
                                         startangle=90)
    ax2.set_title('Aspects by Nature', fontsize=14, weight='bold')
    
    plt.suptitle('Qucanft Demo - Aspects Summary\nJanuary 1, 2024', 
                 fontsize=16, weight='bold', y=1.02)
    plt.tight_layout()
    
    return fig

def create_demo_data_summary():
    """Create a demo data summary"""
    summary = {
        'date': '2024-01-01T12:00:00',
        'location': 'New York, NY',
        'planetary_positions': {
            'Sun': {
                'zodiac_sign': 'Capricorn',
                'position_string': '10°30\' ♑',
                'house': 5,
                'ecliptic_longitude': 280.5
            },
            'Moon': {
                'zodiac_sign': 'Taurus',
                'position_string': '15°12\' ♉',
                'house': 11,
                'ecliptic_longitude': 45.2
            },
            'Mercury': {
                'zodiac_sign': 'Capricorn',
                'position_string': '20°06\' ♑',
                'house': 5,
                'ecliptic_longitude': 290.1
            },
            'Venus': {
                'zodiac_sign': 'Aquarius',
                'position_string': '15°48\' ♒',
                'house': 6,
                'ecliptic_longitude': 315.8
            },
            'Mars': {
                'zodiac_sign': 'Cancer',
                'position_string': '00°18\' ♋',
                'house': 1,
                'ecliptic_longitude': 120.3
            },
            'Jupiter': {
                'zodiac_sign': 'Taurus',
                'position_string': '00°42\' ♉',
                'house': 10,
                'ecliptic_longitude': 30.7
            },
            'Saturn': {
                'zodiac_sign': 'Pisces',
                'position_string': '20°12\' ♓',
                'house': 8,
                'ecliptic_longitude': 350.2
            }
        },
        'aspects_summary': {
            'total_aspects': 8,
            'by_type': {
                'Conjunction': 2,
                'Sextile': 1,
                'Square': 2,
                'Trine': 2,
                'Opposition': 1
            },
            'by_nature': {
                'Harmonious': 3,
                'Challenging': 3,
                'Neutral': 2
            }
        },
        'house_system': 'Equal House',
        'ascendant': '15°00\' ♊ (Gemini)'
    }
    
    return summary

# Generate demo outputs
def generate_demo_outputs():
    """Generate all demo outputs"""
    print("Generating Qucanft demo visualizations...")
    
    # Create natal chart
    natal_fig = create_demo_natal_chart()
    natal_fig.savefig('/tmp/demo_natal_chart.png', dpi=150, bbox_inches='tight')
    print("✓ Demo natal chart saved to /tmp/demo_natal_chart.png")
    
    # Create planetary distribution chart
    dist_fig = create_demo_planetary_distribution()
    dist_fig.savefig('/tmp/demo_planetary_distribution.png', dpi=150, bbox_inches='tight')
    print("✓ Demo planetary distribution chart saved to /tmp/demo_planetary_distribution.png")
    
    # Create aspects summary
    aspects_fig = create_demo_aspects_summary()
    aspects_fig.savefig('/tmp/demo_aspects_summary.png', dpi=150, bbox_inches='tight')
    print("✓ Demo aspects summary saved to /tmp/demo_aspects_summary.png")
    
    # Create data summary
    summary = create_demo_data_summary()
    with open('/tmp/demo_data_summary.json', 'w') as f:
        json.dump(summary, f, indent=2)
    print("✓ Demo data summary saved to /tmp/demo_data_summary.json")
    
    plt.close('all')  # Close all figures
    
    print("\nDemo visualizations completed!")
    print("Files created:")
    print("- /tmp/demo_natal_chart.png")
    print("- /tmp/demo_planetary_distribution.png")
    print("- /tmp/demo_aspects_summary.png")
    print("- /tmp/demo_data_summary.json")
    
    return True

if __name__ == "__main__":
    generate_demo_outputs()