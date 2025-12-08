#!/usr/bin/env python3
"""
Parse Dreamtonics Synthesizer V Studio Scripting API documentation and generate TypeScript type definitions.

This script parses the HTML documentation files and extracts:
- Class names
- Method signatures (parameters and return types)
- Method descriptions
- Inheritance relationships
"""

import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple


class MethodInfo:
    """Information about a method."""
    def __init__(self):
        self.name = ""
        self.params: List[Tuple[str, str, str]] = []  # (name, type, description)
        self.return_type = "void"
        self.description = ""
        self.inherited_from: Optional[str] = None
        self.is_static = False


class ClassInfo:
    """Information about a class."""
    def __init__(self, name: str):
        self.name = name
        self.methods: Dict[str, MethodInfo] = {}  # Use dict to avoid duplicates
        self.description = ""
        self.extends: Optional[str] = None


def clean_html(text: str) -> str:
    """Remove HTML tags and clean up text."""
    # Remove HTML tags but keep the content
    text = re.sub(r'<[^>]+>', '', text)
    # Clean up whitespace
    text = re.sub(r'\s+', ' ', text)
    return text.strip()


def parse_html_file(filepath: Path) -> Optional[ClassInfo]:
    """Parse a single HTML file and extract class information."""
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()

        # Extract class name from filename
        class_name = filepath.stem
        class_info = ClassInfo(class_name)

        # Find all method sections by looking for <h4 class="name" id="...">
        # Match the entire method section until the next <h4>, <hr>, or end
        method_sections = re.findall(
            r'<h4[^>]*class="name"[^>]*id="([^"]+)"[^>]*>(.*?)(?=<h4[^>]*class="name"|<hr|$)',
            content,
            re.DOTALL
        )

        for method_id, section_content in method_sections:
            # Extract method name from the beginning of section_content
            name_match = re.match(r'^([^<]+)<', section_content)
            if not name_match:
                continue

            method_name = name_match.group(1).strip()

            # Skip if we already have this method
            if method_name in class_info.methods:
                continue

            method = MethodInfo()
            method.name = method_name

            # Extract signature (parameters)
            sig_match = re.search(r'<span class="signature">\(([^)]*)\)</span>', section_content)
            params_str = sig_match.group(1).strip() if sig_match else ""

            # Extract return type from type-signature span
            return_match = re.search(r'<span class="type-signature">→\s*\{([^}]+)\}</span>', section_content, re.DOTALL)
            if return_match:
                return_type_html = return_match.group(1)
                method.return_type = clean_html(return_type_html)
            else:
                method.return_type = "void"

            # Extract description (first <p> in description usertext div)
            desc_match = re.search(r'<div[^>]*class="description\s+usertext"[^>]*>\s*<p>(.*?)</p>', section_content, re.DOTALL)
            if desc_match:
                desc = desc_match.group(1)
                desc = re.sub(r'<code>([^<]+)</code>', r'`\1`', desc)  # Convert code tags to backticks
                desc = clean_html(desc)
                method.description = desc

            # Extract parameters from table
            param_table_match = re.search(r'<h5>Parameters:</h5>.*?<table[^>]*class="params"[^>]*>(.*?)</table>', section_content, re.DOTALL)
            if param_table_match:
                param_table = param_table_match.group(1)

                # Extract parameter rows
                param_rows = re.findall(
                    r'<tr>\s*<td[^>]*class="name"[^>]*><code>([^<]+)</code></td>\s*<td[^>]*class="type"[^>]*><span[^>]*class="param-type"[^>]*>([^<]+)</span></td>\s*<td[^>]*class="description[^"]*"[^>]*>([^<]*)</td>',
                    param_table,
                    re.DOTALL
                )

                for param_name, param_type, param_desc in param_rows:
                    method.params.append((
                        param_name.strip(),
                        param_type.strip(),
                        param_desc.strip()
                    ))

            # Check if method is inherited
            inherited_match = re.search(r'<dt[^>]*class="inherited-from"[^>]*>.*?<a[^>]*href="([^"#]+)\.html#', section_content, re.DOTALL)
            if inherited_match:
                method.inherited_from = inherited_match.group(1)

            # Add method to class
            class_info.methods[method_name] = method

        # Check for class inheritance (extends)
        # Look for any inherited method to determine parent class
        for method in class_info.methods.values():
            if method.inherited_from and not class_info.extends:
                class_info.extends = method.inherited_from
                break

        return class_info

    except Exception as e:
        print(f"Error parsing {filepath}: {e}")
        import traceback
        traceback.print_exc()
        return None


def convert_type_to_typescript(doc_type: str, class_name: str = "") -> str:
    """Convert documentation type to TypeScript type."""
    if not doc_type or doc_type == "void":
        return "void"

    # Clean up the type string
    doc_type = doc_type.strip()

    # Common type mappings
    type_map = {
        "string": "string",
        "number": "number",
        "boolean": "boolean",
        "object": "any",
        "array": "any[]",
        "Array": "any[]",
        "function": "Function",
        "Function": "Function",
        "undefined": "undefined",
        "null": "null",
    }

    # Check if it's in the map
    if doc_type in type_map:
        return type_map[doc_type]

    # Handle array types like "Array.<Type>"
    array_match = re.match(r'Array\.<(\w+)>', doc_type)
    if array_match:
        inner_type = convert_type_to_typescript(array_match.group(1))
        return f"{inner_type}[]"

    # If it looks like a class name (starts with uppercase), keep it
    if doc_type and doc_type[0].isupper():
        return doc_type

    # Default to any for unknown types
    return "any"


def generate_typescript_definitions(classes: List[ClassInfo]) -> str:
    """Generate TypeScript definition file content."""
    lines = []

    # Header
    lines.append("/**")
    lines.append(" * Type definitions for Dreamtonics Synthesizer V Studio Scripting API")
    lines.append(" * Generated from official documentation")
    lines.append(" * https://resource.dreamtonics.com/scripting/index.html")
    lines.append(" */")
    lines.append("")

    # Sort classes by name
    classes_sorted = sorted(classes, key=lambda c: c.name)

    # Generate each class
    for class_info in classes_sorted:
        if class_info.description:
            lines.append("/**")
            lines.append(f" * {class_info.description}")
            lines.append(" */")

        # Class declaration
        extends_clause = f" extends {class_info.extends}" if class_info.extends else ""
        lines.append(f"declare class {class_info.name}{extends_clause} {{")

        # Sort methods alphabetically
        sorted_methods = sorted(class_info.methods.values(), key=lambda m: m.name)

        # Generate methods
        for method in sorted_methods:
            # Skip inherited methods (they'll be in the parent class)
            if method.inherited_from:
                continue

            # Method documentation
            if method.description or method.params or method.return_type != "void":
                lines.append("  /**")

                if method.description:
                    lines.append(f"   * {method.description}")

                # Parameter documentation
                for param_name, param_type, param_desc in method.params:
                    ts_type = convert_type_to_typescript(param_type, class_info.name)
                    if param_desc:
                        lines.append(f"   * @param {param_name} {param_desc}")
                    else:
                        lines.append(f"   * @param {param_name}")

                # Return documentation
                if method.return_type != "void":
                    ts_return = convert_type_to_typescript(method.return_type, class_info.name)
                    lines.append(f"   * @returns {ts_return}")

                lines.append("   */")

            # Method signature
            static_keyword = "static " if method.is_static else ""
            params_str = ", ".join([
                f"{name}: {convert_type_to_typescript(ptype, class_info.name)}"
                for name, ptype, _ in method.params
            ])
            return_type = convert_type_to_typescript(method.return_type, class_info.name)

            lines.append(f"  {static_keyword}{method.name}({params_str}): {return_type};")
            lines.append("")

        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def main():
    """Main function to generate TypeScript definitions."""
    # Input and output paths
    docs_dir = Path("reamtonics-api")
    output_file = Path("synthesizer-v-api.d.ts")

    if not docs_dir.exists():
        print(f"Error: Documentation directory not found: {docs_dir}")
        print("Please run download_docs.py first.")
        return 1

    print(f"Parsing documentation from: {docs_dir}")
    print("-" * 60)

    # Parse all HTML files
    classes: List[ClassInfo] = []
    html_files = sorted(docs_dir.glob("*.html"))

    for html_file in html_files:
        # Skip index.html
        if html_file.name == "index.html":
            continue

        print(f"Parsing: {html_file.name}")
        class_info = parse_html_file(html_file)
        if class_info and class_info.methods:
            classes.append(class_info)
            print(f"  Found {len(class_info.methods)} methods")
        else:
            print(f"  No methods found")

    print("-" * 60)
    print(f"Parsed {len(classes)} classes")

    # Generate TypeScript definitions
    print(f"\nGenerating TypeScript definitions...")
    ts_content = generate_typescript_definitions(classes)

    # Write to file
    with open(output_file, "w", encoding="utf-8") as f:
        f.write(ts_content)

    print(f"✓ TypeScript definitions written to: {output_file.absolute()}")

    # Statistics
    total_methods = sum(len(c.methods) for c in classes)
    print(f"\nStatistics:")
    print(f"  Classes: {len(classes)}")
    print(f"  Total methods: {total_methods}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
