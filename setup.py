from setuptools import setup, find_packages

setup(
    name="qucanft",
    version="0.1.0",
    description="Astronomical data fetching and astrological calculations using Astroquery",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    author="Qucanft",
    packages=find_packages(),
    install_requires=[
        "astroquery>=0.4.6",
        "numpy>=1.21.0",
        "pandas>=1.3.0",
        "matplotlib>=3.5.0",
        "astropy>=5.0.0",
        "pytz>=2021.1",
        "swisseph>=2.10.0",
    ],
    python_requires=">=3.7",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Science/Research",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Topic :: Scientific/Engineering :: Astronomy",
    ],
)