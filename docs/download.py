#!/usr/bin/env python3
"""
Download all Dreamtonics Synthesizer V Studio Scripting API documentation.

This script downloads the main index page and all individual class documentation pages
from the official Dreamtonics scripting documentation website.
"""

import os
import sys
import time
from pathlib import Path
from urllib.request import urlopen, Request
from urllib.error import HTTPError, URLError


# Base URL for the documentation
BASE_URL = "https://resource.dreamtonics.com/scripting/"

# List of all documented classes
CLASSES = [
    "ArrangementSelectionState",
    "ArrangementView",
    "Automation",
    "CoordinateSystem",
    "GroupSelection",
    "MainEditorView",
    "NestedObject",
    "Note",
    "NoteGroup",
    "NoteGroupReference",
    "PitchControlCurve",
    "PitchControlPoint",
    "PlaybackControl",
    "Project",
    "RetakeList",
    "SV",
    "ScriptableNestedObject",
    "SelectionStateBase",
    "TimeAxis",
    "Track",
    "TrackInnerSelectionState",
    "TrackMixer",
    "WidgetValue",
]

# Additional pages to download
ADDITIONAL_PAGES = [
    "index.html",
]


def download_page(url, output_path, delay=0.5):
    """
    Download a single page from the URL to the output path.

    Args:
        url: Full URL to download from
        output_path: Path where the file should be saved
        delay: Delay in seconds between requests (to be respectful to the server)
    """
    try:
        print(f"Downloading: {url}")

        # Create a request with a user agent
        headers = {
            'User-Agent': 'Mozilla/5.0 (compatible; DocDownloader/1.0)'
        }
        req = Request(url, headers=headers)

        # Download the content
        with urlopen(req, timeout=30) as response:
            content = response.read()

        # Save to file
        with open(output_path, 'wb') as f:
            f.write(content)

        print(f"  ✓ Saved to: {output_path}")

        # Be respectful to the server
        time.sleep(delay)
        return True

    except HTTPError as e:
        print(f"  ✗ HTTP Error {e.code}: {e.reason}")
        return False
    except URLError as e:
        print(f"  ✗ URL Error: {e.reason}")
        return False
    except Exception as e:
        print(f"  ✗ Error: {e}")
        return False


def main():
    """Main function to download all documentation."""
    # Create output directory
    output_dir = Path("dreamtonics-api")
    output_dir.mkdir(parents=True, exist_ok=True)

    print(f"Downloading Dreamtonics Scripting API documentation to: {output_dir}")
    print(f"Base URL: {BASE_URL}")
    print("-" * 60)

    success_count = 0
    fail_count = 0

    # Download additional pages (like index.html)
    for page in ADDITIONAL_PAGES:
        url = BASE_URL + page
        output_path = output_dir / page
        if download_page(url, output_path):
            success_count += 1
        else:
            fail_count += 1

    # Download each class documentation page
    for class_name in CLASSES:
        url = BASE_URL + class_name + ".html"
        output_path = output_dir / f"{class_name}.html"
        if download_page(url, output_path):
            success_count += 1
        else:
            fail_count += 1

    # Summary
    print("-" * 60)
    print(f"Download complete!")
    print(f"  ✓ Successful: {success_count}")
    print(f"  ✗ Failed: {fail_count}")
    print(f"  Total: {success_count + fail_count}")
    print(f"\nDocumentation saved to: {output_dir.absolute()}")

    return 0 if fail_count == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
