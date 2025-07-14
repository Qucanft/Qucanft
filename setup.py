#!/usr/bin/env python3
"""
Setup script for Qucanft package.
"""

from setuptools import setup, find_packages
import os

# Read the README file for long description
def read_readme():
    with open("README.md", "r", encoding="utf-8") as fh:
        return fh.read()

# Read version from package
def read_version():
    version_file = os.path.join("qucanft", "__init__.py")
    with open(version_file, "r", encoding="utf-8") as fh:
        for line in fh:
            if line.startswith("__version__"):
                return line.split("=")[1].strip().strip('"\'')
    return "0.1.0"

setup(
    name="qucanft",
    version=read_version(),
    description="A Python package for NFT and quantum computing functionality",
    long_description=read_readme(),
    long_description_content_type="text/markdown",
    author="Qucanft Team",
    author_email="contact@qucanft.com",
    url="https://github.com/Qucanft/Qucanft",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Scientific/Engineering",
    ],
    python_requires=">=3.7",
    install_requires=[
        # Add dependencies here as needed
        # Example: "numpy>=1.20.0",
        # Example: "requests>=2.25.0",
    ],
    extras_require={
        "dev": [
            "pytest>=6.0",
            "pytest-cov>=2.0",
            "black>=21.0",
            "flake8>=3.8",
            "mypy>=0.800",
        ],
    },
    include_package_data=True,
    zip_safe=False,
    keywords="nft quantum computing blockchain",
    project_urls={
        "Bug Reports": "https://github.com/Qucanft/Qucanft/issues",
        "Source": "https://github.com/Qucanft/Qucanft",
    },
)