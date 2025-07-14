"""
Visualization helper module for formatting astronomical and astrological data.

This module provides utilities for creating visual representations of
planetary positions, aspects, and other astrological data.
"""

from typing import Dict, List, Optional, Tuple, Union
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import matplotlib.patches as patches
from matplotlib.patches import Circle, Wedge
import math


class VisualizationHelper:
    """
    A class providing visualization utilities for astronomical and astrological data.
    
    This class formats data for visual representation and creates charts
    for planetary positions, aspects, and other astrological information.
    """
    
    # Zodiac sign colors for visualization
    ZODIAC_COLORS = {
        'Aries': '#FF6B6B',
        'Taurus': '#4ECDC4',
        'Gemini': '#FFE66D',
        'Cancer': '#95E1D3',
        'Leo': '#F8B500',
        'Virgo': '#C7CEEA',
        'Libra': '#FFB6C1',
        'Scorpio': '#8B0000',
        'Sagittarius': '#9B59B6',
        'Capricorn': '#2C3E50',
        'Aquarius': '#3498DB',
        'Pisces': '#1ABC9C'
    }
    
    # Planet colors and symbols
    PLANET_COLORS = {
        'Sun': '#FFD700',
        'Moon': '#C0C0C0',
        'Mercury': '#FFA500',
        'Venus': '#FF69B4',
        'Mars': '#FF4500',
        'Jupiter': '#1E90FF',
        'Saturn': '#8B4513',
        'Uranus': '#00CED1',
        'Neptune': '#4169E1',
        'Pluto': '#8B008B'
    }
    
    PLANET_SYMBOLS = {
        'Sun': '☉',
        'Moon': '☽',
        'Mercury': '☿',
        'Venus': '♀',
        'Mars': '♂',
        'Jupiter': '♃',
        'Saturn': '♄',
        'Uranus': '♅',
        'Neptune': '♆',
        'Pluto': '♇'
    }
    
    def __init__(self):
        """Initialize the VisualizationHelper."""
        pass
    
    def format_planetary_table(self, planetary_data: pd.DataFrame) -> pd.DataFrame:
        """
        Format planetary data for display in a table.
        
        Args:
            planetary_data: DataFrame with planetary positions
            
        Returns:
            Formatted DataFrame suitable for display
        """
        if planetary_data.empty:
            return pd.DataFrame()
        
        # Select and format columns for display
        display_columns = ['Planet', 'Zodiac_Sign', 'Position_String']
        
        # Add house information if available
        if 'House' in planetary_data.columns:
            display_columns.append('House')
        
        # Add aspect information if available  
        if 'Ecliptic_Longitude' in planetary_data.columns:
            display_columns.append('Ecliptic_Longitude')
        
        # Create formatted DataFrame
        formatted_df = planetary_data[display_columns].copy()
        
        # Add symbols
        if 'Planet' in formatted_df.columns:
            formatted_df['Symbol'] = formatted_df['Planet'].map(self.PLANET_SYMBOLS)
        
        return formatted_df
    
    def format_aspects_table(self, aspects_data: pd.DataFrame) -> pd.DataFrame:
        """
        Format aspects data for display in a table.
        
        Args:
            aspects_data: DataFrame with aspects
            
        Returns:
            Formatted DataFrame suitable for display
        """
        if aspects_data.empty:
            return pd.DataFrame()
        
        # Create formatted DataFrame
        formatted_df = aspects_data.copy()
        
        # Create a combined planet column for easier reading
        formatted_df['Planets'] = (
            formatted_df['planet1'] + ' ' + 
            formatted_df['symbol'] + ' ' + 
            formatted_df['planet2']
        )
        
        # Format orb as readable text
        formatted_df['Orb'] = formatted_df['orb_difference'].apply(
            lambda x: f"{x:.1f}°"
        )
        
        # Select display columns
        display_columns = ['Planets', 'aspect', 'nature', 'Orb', 'exactness']
        
        return formatted_df[display_columns]
    
    def create_natal_chart(self, 
                          planetary_data: pd.DataFrame,
                          house_cusps: Optional[Dict[int, float]] = None,
                          aspects_data: Optional[pd.DataFrame] = None,
                          chart_size: Tuple[int, int] = (12, 12)) -> plt.Figure:
        """
        Create a basic natal chart visualization.
        
        Args:
            planetary_data: DataFrame with planetary positions
            house_cusps: Dictionary with house cusps
            aspects_data: DataFrame with aspects
            chart_size: Size of the chart (width, height)
            
        Returns:
            Matplotlib figure object
        """
        fig, ax = plt.subplots(figsize=chart_size)
        ax.set_xlim(-1.5, 1.5)
        ax.set_ylim(-1.5, 1.5)
        ax.set_aspect('equal')
        ax.axis('off')
        
        # Draw outer circle (zodiac wheel)
        outer_circle = Circle((0, 0), 1.2, fill=False, color='black', linewidth=2)
        ax.add_patch(outer_circle)
        
        # Draw zodiac signs around the circle
        zodiac_signs = ['Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo',
                       'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces']
        
        for i, sign in enumerate(zodiac_signs):
            angle = (i * 30 - 90) * math.pi / 180  # Start from Aries at top
            x = 1.35 * math.cos(angle)
            y = 1.35 * math.sin(angle)
            ax.text(x, y, sign[:3], ha='center', va='center', fontsize=10, weight='bold')
        
        # Draw house cusps if provided
        if house_cusps:
            for house_num, cusp_deg in house_cusps.items():
                angle = (cusp_deg - 90) * math.pi / 180
                x1 = 0.8 * math.cos(angle)
                y1 = 0.8 * math.sin(angle)
                x2 = 1.2 * math.cos(angle)
                y2 = 1.2 * math.sin(angle)
                ax.plot([x1, x2], [y1, y2], 'gray', alpha=0.7)
                
                # House number
                x_text = 0.9 * math.cos(angle)
                y_text = 0.9 * math.sin(angle)
                ax.text(x_text, y_text, str(house_num), ha='center', va='center', 
                       fontsize=8, color='gray')
        
        # Draw planets
        if 'Ecliptic_Longitude' in planetary_data.columns:
            for _, planet in planetary_data.iterrows():
                angle = (planet['Ecliptic_Longitude'] - 90) * math.pi / 180
                x = 1.05 * math.cos(angle)
                y = 1.05 * math.sin(angle)
                
                # Planet symbol
                symbol = self.PLANET_SYMBOLS.get(planet['Planet'], planet['Planet'][:2])
                color = self.PLANET_COLORS.get(planet['Planet'], 'black')
                
                ax.scatter(x, y, s=200, color=color, alpha=0.8, edgecolors='black')
                ax.text(x, y, symbol, ha='center', va='center', fontsize=12, 
                       color='white', weight='bold')
        
        # Draw aspects if provided
        if aspects_data is not None and not aspects_data.empty:
            for _, aspect in aspects_data.iterrows():
                # Get planet positions
                planet1_data = planetary_data[planetary_data['Planet'] == aspect['planet1']]
                planet2_data = planetary_data[planetary_data['Planet'] == aspect['planet2']]
                
                if not planet1_data.empty and not planet2_data.empty:
                    pos1 = planet1_data.iloc[0]['Ecliptic_Longitude']
                    pos2 = planet2_data.iloc[0]['Ecliptic_Longitude']
                    
                    angle1 = (pos1 - 90) * math.pi / 180
                    angle2 = (pos2 - 90) * math.pi / 180
                    
                    x1 = 1.05 * math.cos(angle1)
                    y1 = 1.05 * math.sin(angle1)
                    x2 = 1.05 * math.cos(angle2)
                    y2 = 1.05 * math.sin(angle2)
                    
                    # Color based on aspect nature
                    if aspect['nature'] == 'Harmonious':
                        color = 'blue'
                        alpha = 0.6
                    elif aspect['nature'] == 'Challenging':
                        color = 'red'
                        alpha = 0.6
                    else:
                        color = 'gray'
                        alpha = 0.4
                    
                    ax.plot([x1, x2], [y1, y2], color=color, alpha=alpha, linewidth=1)
        
        plt.title('Natal Chart', fontsize=16, weight='bold', pad=20)
        return fig
    
    def create_planetary_positions_chart(self, planetary_data: pd.DataFrame) -> plt.Figure:
        """
        Create a chart showing planetary positions across zodiac signs.
        
        Args:
            planetary_data: DataFrame with planetary positions
            
        Returns:
            Matplotlib figure object
        """
        if 'Zodiac_Sign' not in planetary_data.columns:
            raise ValueError("DataFrame must contain 'Zodiac_Sign' column")
        
        fig, ax = plt.subplots(figsize=(14, 8))
        
        # Count planets in each sign
        sign_counts = planetary_data['Zodiac_Sign'].value_counts()
        
        # Prepare data for plotting
        zodiac_order = ['Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo',
                       'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces']
        
        counts = [sign_counts.get(sign, 0) for sign in zodiac_order]
        colors = [self.ZODIAC_COLORS[sign] for sign in zodiac_order]
        
        # Create bar chart
        bars = ax.bar(zodiac_order, counts, color=colors, alpha=0.8, edgecolor='black')
        
        # Add planet names on bars
        for i, (sign, count) in enumerate(zip(zodiac_order, counts)):
            if count > 0:
                planets_in_sign = planetary_data[planetary_data['Zodiac_Sign'] == sign]['Planet'].tolist()
                planet_text = ', '.join(planets_in_sign)
                ax.text(i, count + 0.1, planet_text, ha='center', va='bottom', 
                       fontsize=9, rotation=45)
        
        ax.set_xlabel('Zodiac Signs', fontsize=12, weight='bold')
        ax.set_ylabel('Number of Planets', fontsize=12, weight='bold')
        ax.set_title('Planetary Distribution Across Zodiac Signs', fontsize=14, weight='bold')
        ax.set_ylim(0, max(counts) + 2)
        
        plt.xticks(rotation=45, ha='right')
        plt.tight_layout()
        
        return fig
    
    def create_aspects_summary_chart(self, aspects_data: pd.DataFrame) -> plt.Figure:
        """
        Create a chart summarizing aspects by type and nature.
        
        Args:
            aspects_data: DataFrame with aspects
            
        Returns:
            Matplotlib figure object
        """
        if aspects_data.empty:
            fig, ax = plt.subplots(figsize=(10, 6))
            ax.text(0.5, 0.5, 'No aspects found', ha='center', va='center', 
                   fontsize=16, transform=ax.transAxes)
            ax.set_title('Aspects Summary', fontsize=14, weight='bold')
            return fig
        
        fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))
        
        # Chart 1: Aspects by type
        aspect_counts = aspects_data['aspect'].value_counts()
        colors1 = plt.cm.Set3(np.linspace(0, 1, len(aspect_counts)))
        
        wedges, texts, autotexts = ax1.pie(aspect_counts.values, labels=aspect_counts.index,
                                          autopct='%1.1f%%', colors=colors1, startangle=90)
        ax1.set_title('Aspects by Type', fontsize=12, weight='bold')
        
        # Chart 2: Aspects by nature
        nature_counts = aspects_data['nature'].value_counts()
        nature_colors = {'Harmonious': 'lightblue', 'Challenging': 'lightcoral', 
                        'Neutral': 'lightgray', 'Mild': 'lightgreen', 
                        'Adjusting': 'lightyellow', 'Creative': 'plum',
                        'Spiritual': 'lavender'}
        colors2 = [nature_colors.get(nature, 'lightgray') for nature in nature_counts.index]
        
        wedges2, texts2, autotexts2 = ax2.pie(nature_counts.values, labels=nature_counts.index,
                                             autopct='%1.1f%%', colors=colors2, startangle=90)
        ax2.set_title('Aspects by Nature', fontsize=12, weight='bold')
        
        plt.tight_layout()
        return fig
    
    def export_data_summary(self, 
                           planetary_data: pd.DataFrame,
                           aspects_data: Optional[pd.DataFrame] = None,
                           house_cusps: Optional[Dict[int, float]] = None) -> Dict[str, any]:
        """
        Create a comprehensive data summary for export.
        
        Args:
            planetary_data: DataFrame with planetary positions
            aspects_data: DataFrame with aspects (optional)
            house_cusps: Dictionary with house cusps (optional)
            
        Returns:
            Dictionary with formatted data summary
        """
        summary = {
            'planetary_positions': {},
            'zodiac_distribution': {},
            'aspects_summary': {},
            'house_cusps': house_cusps or {}
        }
        
        # Planetary positions
        for _, planet in planetary_data.iterrows():
            summary['planetary_positions'][planet['Planet']] = {
                'zodiac_sign': planet.get('Zodiac_Sign', 'Unknown'),
                'position_string': planet.get('Position_String', 'Unknown'),
                'house': planet.get('House', 'Unknown'),
                'ecliptic_longitude': planet.get('Ecliptic_Longitude', 'Unknown')
            }
        
        # Zodiac distribution
        if 'Zodiac_Sign' in planetary_data.columns:
            sign_counts = planetary_data['Zodiac_Sign'].value_counts()
            summary['zodiac_distribution'] = sign_counts.to_dict()
        
        # Aspects summary
        if aspects_data is not None and not aspects_data.empty:
            summary['aspects_summary'] = {
                'total_aspects': len(aspects_data),
                'by_type': aspects_data['aspect'].value_counts().to_dict(),
                'by_nature': aspects_data['nature'].value_counts().to_dict(),
                'strongest_aspects': aspects_data.nsmallest(5, 'orb_difference')[
                    ['planet1', 'planet2', 'aspect', 'orb_difference']
                ].to_dict('records')
            }
        
        return summary
    
    def create_ephemeris_plot(self, 
                            ephemeris_data: pd.DataFrame,
                            planet: str = 'Sun') -> plt.Figure:
        """
        Create a plot showing planetary motion over time.
        
        Args:
            ephemeris_data: DataFrame with ephemeris data over time
            planet: Name of the planet to plot
            
        Returns:
            Matplotlib figure object
        """
        if ephemeris_data.empty:
            fig, ax = plt.subplots(figsize=(12, 6))
            ax.text(0.5, 0.5, 'No ephemeris data available', ha='center', va='center',
                   fontsize=16, transform=ax.transAxes)
            return fig
        
        fig, ax = plt.subplots(figsize=(12, 6))
        
        # Assume the ephemeris data has datetime and position columns
        if 'datetime_jd' in ephemeris_data.columns and 'RA' in ephemeris_data.columns:
            ax.plot(ephemeris_data['datetime_jd'], ephemeris_data['RA'], 
                   label='Right Ascension', color='blue', linewidth=2)
            
            if 'DEC' in ephemeris_data.columns:
                ax2 = ax.twinx()
                ax2.plot(ephemeris_data['datetime_jd'], ephemeris_data['DEC'], 
                        label='Declination', color='red', linewidth=2)
                ax2.set_ylabel('Declination (degrees)', color='red')
                ax2.tick_params(axis='y', labelcolor='red')
        
        ax.set_xlabel('Time (Julian Date)')
        ax.set_ylabel('Right Ascension (degrees)', color='blue')
        ax.tick_params(axis='y', labelcolor='blue')
        ax.set_title(f'{planet} Motion Over Time', fontsize=14, weight='bold')
        
        plt.tight_layout()
        return fig