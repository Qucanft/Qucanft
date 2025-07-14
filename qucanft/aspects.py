"""
Planetary aspects calculation module.

This module provides functionality to calculate astrological aspects between planets,
including major and minor aspects, and their interpretations.
"""

from typing import Dict, List, Optional, Tuple
import numpy as np
import pandas as pd
from datetime import datetime
import math


class AspectsCalculator:
    """
    A class for calculating astrological aspects between planets.
    
    This class provides methods to calculate various types of aspects,
    their orbs, and interpretations.
    """
    
    # Major aspects with their degrees and default orbs
    MAJOR_ASPECTS = {
        'Conjunction': {'degrees': 0, 'orb': 8, 'symbol': '☌', 'nature': 'Neutral'},
        'Sextile': {'degrees': 60, 'orb': 6, 'symbol': '⚹', 'nature': 'Harmonious'},
        'Square': {'degrees': 90, 'orb': 8, 'symbol': '□', 'nature': 'Challenging'},
        'Trine': {'degrees': 120, 'orb': 8, 'symbol': '△', 'nature': 'Harmonious'},
        'Opposition': {'degrees': 180, 'orb': 8, 'symbol': '☍', 'nature': 'Challenging'}
    }
    
    # Minor aspects with their degrees and default orbs
    MINOR_ASPECTS = {
        'Semisextile': {'degrees': 30, 'orb': 3, 'symbol': '⚺', 'nature': 'Mild'},
        'Semisquare': {'degrees': 45, 'orb': 3, 'symbol': '∠', 'nature': 'Challenging'},
        'Sesquiquadrate': {'degrees': 135, 'orb': 3, 'symbol': '⚼', 'nature': 'Challenging'},
        'Quincunx': {'degrees': 150, 'orb': 3, 'symbol': '⚻', 'nature': 'Adjusting'},
        'Quintile': {'degrees': 72, 'orb': 2, 'symbol': 'Q', 'nature': 'Creative'},
        'Biquintile': {'degrees': 144, 'orb': 2, 'symbol': 'bQ', 'nature': 'Creative'},
        'Septile': {'degrees': 51.4, 'orb': 1, 'symbol': 'S', 'nature': 'Spiritual'},
        'Novile': {'degrees': 40, 'orb': 1, 'symbol': 'N', 'nature': 'Spiritual'}
    }
    
    # Planet-specific orb adjustments (luminaries get larger orbs)
    PLANET_ORB_ADJUSTMENTS = {
        'Sun': 1.5,
        'Moon': 1.5,
        'Mercury': 1.0,
        'Venus': 1.0,
        'Mars': 1.0,
        'Jupiter': 1.2,
        'Saturn': 1.2,
        'Uranus': 0.8,
        'Neptune': 0.8,
        'Pluto': 0.8
    }
    
    def __init__(self, include_minor_aspects: bool = True):
        """
        Initialize the AspectsCalculator.
        
        Args:
            include_minor_aspects: Whether to include minor aspects in calculations
        """
        self.include_minor_aspects = include_minor_aspects
        self.aspects = self.MAJOR_ASPECTS.copy()
        if include_minor_aspects:
            self.aspects.update(self.MINOR_ASPECTS)
    
    def calculate_angular_distance(self, pos1: float, pos2: float) -> float:
        """
        Calculate the angular distance between two positions.
        
        Args:
            pos1: First position in degrees (0-360)
            pos2: Second position in degrees (0-360)
            
        Returns:
            Angular distance in degrees (0-180)
        """
        # Normalize positions to 0-360
        pos1 = pos1 % 360
        pos2 = pos2 % 360
        
        # Calculate difference
        diff = abs(pos2 - pos1)
        
        # Return the smaller angle
        return min(diff, 360 - diff)
    
    def find_aspects_between_planets(self, 
                                   planet1: str, 
                                   pos1: float,
                                   planet2: str,
                                   pos2: float) -> List[Dict[str, any]]:
        """
        Find all aspects between two planets.
        
        Args:
            planet1: Name of first planet
            pos1: Position of first planet in degrees
            planet2: Name of second planet
            pos2: Position of second planet in degrees
            
        Returns:
            List of dictionaries containing aspect information
        """
        aspects_found = []
        angular_distance = self.calculate_angular_distance(pos1, pos2)
        
        for aspect_name, aspect_info in self.aspects.items():
            aspect_degrees = aspect_info['degrees']
            base_orb = aspect_info['orb']
            
            # Adjust orb based on planets involved
            orb_adjustment = max(
                self.PLANET_ORB_ADJUSTMENTS.get(planet1, 1.0),
                self.PLANET_ORB_ADJUSTMENTS.get(planet2, 1.0)
            )
            adjusted_orb = base_orb * orb_adjustment
            
            # Check if angular distance is within orb of the aspect
            orb_difference = abs(angular_distance - aspect_degrees)
            
            if orb_difference <= adjusted_orb:
                # Calculate exact orb (how close to perfect aspect)
                exactness = (adjusted_orb - orb_difference) / adjusted_orb * 100
                
                # Determine if aspect is applying or separating
                # (This would require planetary motion data for accuracy)
                applying = "Unknown"  # Placeholder
                
                aspects_found.append({
                    'aspect': aspect_name,
                    'planet1': planet1,
                    'planet2': planet2,
                    'degrees': aspect_degrees,
                    'orb_used': adjusted_orb,
                    'orb_difference': orb_difference,
                    'exactness': exactness,
                    'symbol': aspect_info['symbol'],
                    'nature': aspect_info['nature'],
                    'applying': applying,
                    'angular_distance': angular_distance
                })
        
        return aspects_found
    
    def calculate_all_aspects(self, planetary_data: pd.DataFrame) -> pd.DataFrame:
        """
        Calculate all aspects between all planets in the dataset.
        
        Args:
            planetary_data: DataFrame with planetary positions
            
        Returns:
            DataFrame with all aspects found
        """
        if 'Ecliptic_Longitude' not in planetary_data.columns:
            raise ValueError("DataFrame must contain 'Ecliptic_Longitude' column")
        
        aspects_list = []
        planets = planetary_data['Planet'].tolist()
        positions = planetary_data['Ecliptic_Longitude'].tolist()
        
        # Calculate aspects between all planet pairs
        for i in range(len(planets)):
            for j in range(i + 1, len(planets)):
                planet1 = planets[i]
                planet2 = planets[j]
                pos1 = positions[i]
                pos2 = positions[j]
                
                aspects = self.find_aspects_between_planets(planet1, pos1, planet2, pos2)
                aspects_list.extend(aspects)
        
        return pd.DataFrame(aspects_list)
    
    def get_aspect_interpretation(self, aspect: str, planet1: str, planet2: str) -> str:
        """
        Get a basic interpretation of an aspect between two planets.
        
        Args:
            aspect: Name of the aspect
            planet1: First planet
            planet2: Second planet
            
        Returns:
            String with aspect interpretation
        """
        # Basic aspect interpretations
        interpretations = {
            'Conjunction': f"{planet1} and {planet2} are united, blending their energies.",
            'Sextile': f"{planet1} and {planet2} work together harmoniously, creating opportunities.",
            'Square': f"{planet1} and {planet2} create tension that demands action and growth.",
            'Trine': f"{planet1} and {planet2} flow together easily, bringing natural talents.",
            'Opposition': f"{planet1} and {planet2} are in opposition, creating awareness through contrast.",
            'Semisextile': f"{planet1} and {planet2} create a mild connection requiring adjustment.",
            'Semisquare': f"{planet1} and {planet2} create minor friction that stimulates action.",
            'Sesquiquadrate': f"{planet1} and {planet2} create tension requiring conscious effort.",
            'Quincunx': f"{planet1} and {planet2} require constant adjustment and adaptation.",
            'Quintile': f"{planet1} and {planet2} create opportunities for creative expression.",
            'Biquintile': f"{planet1} and {planet2} support creative and artistic endeavors.",
            'Septile': f"{planet1} and {planet2} create a subtle spiritual connection.",
            'Novile': f"{planet1} and {planet2} connect through spiritual completion."
        }
        
        return interpretations.get(aspect, f"{planet1} and {planet2} form a {aspect} aspect.")
    
    def get_strongest_aspects(self, 
                            aspects_df: pd.DataFrame,
                            limit: int = 10) -> pd.DataFrame:
        """
        Get the strongest aspects (those with smallest orbs).
        
        Args:
            aspects_df: DataFrame with aspects
            limit: Maximum number of aspects to return
            
        Returns:
            DataFrame with strongest aspects
        """
        if aspects_df.empty:
            return aspects_df
        
        # Sort by orb difference (smaller is stronger)
        sorted_aspects = aspects_df.sort_values('orb_difference')
        
        return sorted_aspects.head(limit)
    
    def get_aspects_by_nature(self, aspects_df: pd.DataFrame, nature: str) -> pd.DataFrame:
        """
        Filter aspects by their nature (Harmonious, Challenging, etc.).
        
        Args:
            aspects_df: DataFrame with aspects
            nature: Nature of aspects to filter for
            
        Returns:
            DataFrame with filtered aspects
        """
        if aspects_df.empty:
            return aspects_df
        
        return aspects_df[aspects_df['nature'] == nature]
    
    def get_planet_aspects(self, aspects_df: pd.DataFrame, planet: str) -> pd.DataFrame:
        """
        Get all aspects involving a specific planet.
        
        Args:
            aspects_df: DataFrame with aspects
            planet: Planet name to filter for
            
        Returns:
            DataFrame with aspects involving the specified planet
        """
        if aspects_df.empty:
            return aspects_df
        
        return aspects_df[
            (aspects_df['planet1'] == planet) | 
            (aspects_df['planet2'] == planet)
        ]
    
    def calculate_aspect_patterns(self, aspects_df: pd.DataFrame) -> List[Dict[str, any]]:
        """
        Identify common aspect patterns (Grand Trine, T-Square, etc.).
        
        Args:
            aspects_df: DataFrame with aspects
            
        Returns:
            List of dictionaries with pattern information
        """
        patterns = []
        
        if aspects_df.empty:
            return patterns
        
        # Group aspects by type
        trines = aspects_df[aspects_df['aspect'] == 'Trine']
        squares = aspects_df[aspects_df['aspect'] == 'Square']
        oppositions = aspects_df[aspects_df['aspect'] == 'Opposition']
        
        # Look for Grand Trine (3 planets in trine to each other)
        if len(trines) >= 3:
            trine_planets = set()
            for _, aspect in trines.iterrows():
                trine_planets.add(aspect['planet1'])
                trine_planets.add(aspect['planet2'])
            
            # Check if we have a triangle of trines
            if len(trine_planets) >= 3:
                # Simplified check - in practice, would verify exact triangle
                patterns.append({
                    'pattern': 'Grand Trine',
                    'planets': list(trine_planets)[:3],
                    'description': 'A harmonious triangle of energy flow'
                })
        
        # Look for T-Square (2 squares and 1 opposition)
        if len(squares) >= 2 and len(oppositions) >= 1:
            # Simplified pattern detection
            patterns.append({
                'pattern': 'T-Square',
                'planets': ['Unknown'],  # Would need more complex logic
                'description': 'A challenging pattern requiring action'
            })
        
        # Look for Grand Cross (4 squares and 2 oppositions)
        if len(squares) >= 4 and len(oppositions) >= 2:
            patterns.append({
                'pattern': 'Grand Cross',
                'planets': ['Unknown'],  # Would need more complex logic
                'description': 'A highly dynamic and challenging pattern'
            })
        
        return patterns
    
    def format_aspect_string(self, aspect_row: pd.Series) -> str:
        """
        Format an aspect as a readable string.
        
        Args:
            aspect_row: Row from aspects DataFrame
            
        Returns:
            Formatted aspect string
        """
        return (f"{aspect_row['planet1']} {aspect_row['symbol']} {aspect_row['planet2']} "
                f"({aspect_row['orb_difference']:.1f}° orb)")
    
    def export_aspects_summary(self, aspects_df: pd.DataFrame) -> Dict[str, any]:
        """
        Create a summary of all aspects.
        
        Args:
            aspects_df: DataFrame with aspects
            
        Returns:
            Dictionary with aspects summary
        """
        if aspects_df.empty:
            return {'total_aspects': 0, 'by_nature': {}, 'by_type': {}}
        
        summary = {
            'total_aspects': len(aspects_df),
            'by_nature': aspects_df['nature'].value_counts().to_dict(),
            'by_type': aspects_df['aspect'].value_counts().to_dict(),
            'strongest_aspect': {
                'aspect': aspects_df.loc[aspects_df['orb_difference'].idxmin(), 'aspect'],
                'planets': f"{aspects_df.loc[aspects_df['orb_difference'].idxmin(), 'planet1']} - {aspects_df.loc[aspects_df['orb_difference'].idxmin(), 'planet2']}",
                'orb': aspects_df['orb_difference'].min()
            }
        }
        
        return summary