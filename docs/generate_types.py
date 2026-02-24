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
from html.parser import HTMLParser
from typing import Dict, List, Optional, Tuple


class MethodInfo:
    """Information about a method or property."""
    def __init__(self):
        self.name = ""
        self.params: List[Tuple[str, str, str]] = []  # (name, type, description)
        self.return_type = "void"
        self.return_desc = ""  # Detailed return type description from HTML
        self.description = ""
        self.inherited_from: Optional[str] = None
        self.is_static = False
        self.is_property = False  # True for constants/properties, False for methods


class ClassInfo:
    """Information about a class."""
    def __init__(self, name: str):
        self.name = name
        self.methods: Dict[str, MethodInfo] = {}  # Use dict to avoid duplicates
        self.description = ""
        self.extends: Optional[str] = None


class APIDocParser(HTMLParser):
    """HTML parser for extracting method information from API documentation."""

    def __init__(self):
        super().__init__()
        self.methods: Dict[str, MethodInfo] = {}
        self.current_method: Optional[MethodInfo] = None
        self.current_method_id: Optional[str] = None

        # State tracking
        self.in_method_header = False
        self.in_method_name = False
        self.in_signature = False
        self.in_type_signature = False
        self.in_description = False
        self.in_param_table = False
        self.in_param_row = False
        self.in_param_cell = False
        self.param_cell_type = None  # 'name', 'type', or 'description'
        self.in_code = False
        self.in_inherited_from = False
        self.in_returns_section = False
        self.in_return_desc = False

        # Current parameter being parsed
        self.current_param_name = ""
        self.current_param_type = ""
        self.current_param_desc = ""
        self.current_return_desc = ""

        # Buffer for collecting text
        self.text_buffer = []

    def handle_starttag(self, tag, attrs):
        attrs_dict = dict(attrs)

        # Detect method header: <h4 class="name" id="methodName">
        if tag == "h4" and attrs_dict.get("class") == "name":
            method_id = attrs_dict.get("id")
            if method_id:
                # Save previous method if exists
                if self.current_method:
                    self.methods[self.current_method.name] = self.current_method

                # Start new method
                self.current_method = MethodInfo()
                self.current_method_id = method_id
                self.in_method_header = True
                self.in_method_name = True
                self.text_buffer = []

        # Inside method header, detect signature: <span class="signature">
        elif self.in_method_header and tag == "span" and attrs_dict.get("class") == "signature":
            self.in_signature = True
            self.in_method_name = False

        # Inside method header, detect return type: <span class="type-signature">
        elif self.in_method_header and tag == "span" and attrs_dict.get("class") == "type-signature":
            self.in_type_signature = True
            self.in_method_name = False  # Stop capturing method name
            self.text_buffer = []

        # Detect description: <div class="description usertext">
        elif tag == "div" and "description" in attrs_dict.get("class", "") and "usertext" in attrs_dict.get("class", ""):
            self.in_description = True
            self.text_buffer = []

        # Inside description, detect paragraph
        elif self.in_description and tag == "p":
            self.text_buffer = []

        # Inside description, detect code: <code>
        elif self.in_description and tag == "code":
            self.in_code = True

        # Detect parameter table: <table class="params">
        elif tag == "table" and attrs_dict.get("class") == "params":
            self.in_param_table = True

        # Inside param table, detect row
        elif self.in_param_table and tag == "tr":
            self.in_param_row = True
            self.current_param_name = ""
            self.current_param_type = ""
            self.current_param_desc = ""
            self.param_cell_type = None

        # Inside param row, detect cells
        elif self.in_param_row and tag == "td":
            self.in_param_cell = True
            self.text_buffer = []
            # Determine cell type from class
            cell_class = attrs_dict.get("class", "")
            if "name" in cell_class:
                self.param_cell_type = "name"
            elif "type" in cell_class:
                self.param_cell_type = "type"
            elif "description" in cell_class:
                self.param_cell_type = "description"

        # Detect inherited from: <dt class="inherited-from">
        elif tag == "dt" and attrs_dict.get("class") == "inherited-from":
            self.in_inherited_from = True

        # Inside inherited from, detect link: <a href="ClassName.html#methodName">
        elif self.in_inherited_from and tag == "a":
            href = attrs_dict.get("href", "")
            # Extract class name from href like "ClassName.html#methodName"
            match = re.match(r'([^.]+)\.html#', href)
            if match and self.current_method:
                self.current_method.inherited_from = match.group(1)

        # Detect Returns section: <h5>Returns:</h5>
        elif tag == "h5":
            self.text_buffer = []

        # Detect return description: <div class="param-desc"> after Returns section
        elif self.in_returns_section and tag == "div" and attrs_dict.get("class") == "param-desc":
            self.in_return_desc = True
            self.text_buffer = []

        # Inside return desc, detect paragraph
        elif self.in_return_desc and tag == "p":
            self.text_buffer = []

        # Inside return desc, detect code: <code>
        elif self.in_return_desc and tag == "code":
            self.in_code = True

    def handle_endtag(self, tag):
        # End of method header
        if tag == "h4" and self.in_method_header:
            self.in_method_header = False
            self.in_method_name = False
            self.in_signature = False
            self.in_type_signature = False

        # End of h5 tag - check if it's "Returns:"
        elif tag == "h5":
            text = " ".join(self.text_buffer).strip()
            if text == "Returns:":
                self.in_returns_section = True
                self.current_return_desc = ""
            self.text_buffer = []

        # End of description
        elif tag == "div" and self.in_description:
            self.in_description = False

        # End of return description div
        elif tag == "div" and self.in_return_desc:
            self.in_return_desc = False
            self.in_returns_section = False
            # Save the return description to the current method
            if self.current_method and self.current_return_desc:
                self.current_method.return_desc = self.current_return_desc

        # End of return description paragraph - save description
        elif self.in_return_desc and tag == "p":
            if self.current_method and self.text_buffer:
                desc = " ".join(self.text_buffer).strip()
                desc = re.sub(r'\s+', ' ', desc)
                self.current_return_desc = desc
            self.text_buffer = []

        # End of description paragraph - save description
        elif self.in_description and tag == "p":
            if self.current_method and self.text_buffer:
                # Only save the first paragraph as description
                if not self.current_method.description:
                    desc = " ".join(self.text_buffer).strip()
                    desc = re.sub(r'\s+', ' ', desc)
                    self.current_method.description = desc
            self.text_buffer = []

        # End of code tag
        elif tag == "code":
            self.in_code = False

        # End of parameter table
        elif tag == "table" and self.in_param_table:
            self.in_param_table = False

        # End of parameter row - save parameter
        elif tag == "tr" and self.in_param_row:
            self.in_param_row = False
            if self.current_method and self.current_param_name:
                self.current_method.params.append((
                    self.current_param_name,
                    self.current_param_type,
                    self.current_param_desc
                ))

        # End of parameter cell
        elif tag == "td" and self.in_param_cell:
            self.in_param_cell = False
            text = " ".join(self.text_buffer).strip()
            if self.param_cell_type == "name":
                self.current_param_name = text
            elif self.param_cell_type == "type":
                self.current_param_type = text
            elif self.param_cell_type == "description":
                self.current_param_desc = text
            self.text_buffer = []

        # End of inherited from
        elif tag == "dd" and self.in_inherited_from:
            self.in_inherited_from = False

    def handle_data(self, data):
        data = data.strip()
        if not data:
            return

        # Capture method name
        if self.in_method_name and self.current_method:
            if not self.current_method.name:
                self.current_method.name = data

        # Capture return type from type signature
        elif self.in_type_signature and self.current_method:
            self.text_buffer.append(data)
            full_text = " ".join(self.text_buffer)

            # Check if it's a property type ":Type" or method return "→ {Type}"
            # Property pattern: ":type"
            property_match = re.search(r':\s*(\w+)', full_text)
            if property_match:
                self.current_method.is_property = True
                self.current_method.return_type = property_match.group(1)
            else:
                # Method return type pattern: "→ {Type}"
                method_match = re.search(r'→\s*\{([^}]+)\}', full_text)
                if method_match:
                    self.current_method.return_type = method_match.group(1)

        # Capture h5 text (to detect "Returns:")
        elif self.text_buffer is not None and not self.in_description and not self.in_return_desc and not self.in_param_cell:
            self.text_buffer.append(data)

        # Capture return description text
        elif self.in_return_desc:
            # Add backticks around code
            if self.in_code:
                self.text_buffer.append(f"`{data}`")
            else:
                self.text_buffer.append(data)

        # Capture description text
        elif self.in_description:
            # Add backticks around code
            if self.in_code:
                self.text_buffer.append(f"`{data}`")
            else:
                self.text_buffer.append(data)

        # Capture parameter cell data
        elif self.in_param_cell:
            self.text_buffer.append(data)

    def get_methods(self) -> Dict[str, MethodInfo]:
        """Get all parsed methods. Call this after parsing is complete."""
        # Add the last method if exists
        if self.current_method and self.current_method.name:
            self.methods[self.current_method.name] = self.current_method
        return self.methods


def parse_html_file(filepath: Path) -> Optional[ClassInfo]:
    """Parse a single HTML file and extract class information."""
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()

        # Extract class name from filename
        class_name = filepath.stem
        class_info = ClassInfo(class_name)

        # Parse the HTML
        parser = APIDocParser()
        parser.feed(content)

        # Get the parsed methods
        class_info.methods = parser.get_methods()

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


def parse_return_description(return_desc: str) -> Optional[str]:
    """
    Parse a return type description and convert it to TypeScript type.

    Examples:
    - "an `array` of `array` of `number`" -> "number[][]"
    - "an `array` of `object`" -> "object[]"
    - "An array of `PitchControlPoint` or `PitchControlCurve`" -> "(PitchControlPoint | PitchControlCurve)[]"
    - "a `string`" -> "string"
    - "The stored value, or `undefined` if..." -> None (skip, let type signature handle it)
    """
    if not return_desc:
        return None

    # Extract all types in backticks
    types = re.findall(r'`([^`]+)`', return_desc)
    if not types:
        return None

    # Check for pattern: "array of X or Y" (union types inside array)
    # Pattern: "array of `Type1` or `Type2`"
    if re.search(r'array\s+of.*\s+or\s+', return_desc, re.IGNORECASE):
        # Extract non-array types (all types except "array")
        element_types = [t for t in types if t != "array"]
        if len(element_types) > 1:
            # Build union type and wrap in array
            union = " | ".join(element_types)
            return f"({union})[]"

    # Check for union type patterns like "value or `undefined`" (not inside array)
    # In these cases, prefer the type signature instead
    if "or" in return_desc.lower() and "array" not in return_desc.lower():
        # This might be a union type description, skip parsing
        # Let the type signature from the HTML handle it
        return None

    # Common type mappings
    type_map = {
        "string": "string",
        "number": "number",
        "boolean": "boolean",
        "object": "any",
        "array": None,  # Special handling
        "function": "Function",
        "Function": "Function",
    }

    # Parse nested array structure
    # Look for patterns like "array of array of number"
    # Build from the innermost type outward
    result_type = None

    # Reverse the list to process from innermost to outermost
    for i in range(len(types) - 1, -1, -1):
        type_name = types[i].strip()

        if type_name == "array":
            # Wrap the current result in an array
            if result_type is None:
                result_type = "any[]"
            else:
                result_type = f"{result_type}[]"
        else:
            # This is the base type
            mapped_type = type_map.get(type_name, type_name)
            if mapped_type is None:
                continue
            # If it looks like a class name (starts with uppercase), keep it
            if type_name and type_name[0].isupper():
                result_type = type_name
            else:
                result_type = mapped_type

    return result_type


def convert_type_to_typescript(doc_type: str, class_name: str = "", return_desc: str = "", method_name: str = "", param_name: str = "") -> str:
    """Convert documentation type to TypeScript type."""
    # Parameter-specific type overrides
    param_type_overrides = {
        ("SV", "showCustomDialog", "form"): "Form",
        ("SV", "showCustomDialogAsync", "form"): "Form",

        ("TrackInnerSelectionState", "selectPitchControls", "controls"): "(PitchControlPoint | PitchControlCurve)[]",
        ("TrackInnerSelectionState", "unselectPitchControls", "controls"): "(PitchControlPoint | PitchControlCurve)[]",
        ("TrackInnerSelectionState", "unselectPoints", "positions"): "number[]",

        ("Project", "getNoteGroup", "id"): "number",
    }

    # Check for parameter-specific override first
    if class_name and method_name and param_name:
        override_key = (class_name, method_name, param_name)
        if override_key in param_type_overrides:
            return param_type_overrides[override_key]
        # If we're processing a parameter and no override found, continue with normal processing
        # Do NOT check method overrides for parameters

    # Method-specific type overrides for complex object types (return types only, not parameters)
    method_type_overrides = {
        ("NoteGroupReference", "getVoice"): "VoiceParameters",

        ("SV", "getComputedAttributesForGroup"): "ComputedAttributes[]",
        ("SV", "getComputedPitchForGroup"): "(number|null)[]",
        ("SV", "getHostInfo"): "HostInfo",
        ("SV", "getPlayback"): "PlaybackControl",

        ("TimeAxis", "getAllMeasureMarks"): "MeasureMark[]",
        ("TimeAxis", "getMeasureMarkAt"): "MeasureMark",
        ("TimeAxis", "getMeasureMarkAtBlick"): "MeasureMark",

        ("TimeAxis", "getAllTempoMarks"): "TempoMark[]",
        ("TimeAxis", "getTempoMarkAt"): "TempoMark",

        ("Note", "getAttributes"): "NoteAttributes",

        ("TrackInnerSelectionState", "getSelectedPoints"): "number[]",
        ("TrackInnerSelectionState", "getSelectedNotes"): "Note[]",
    }

    # Check for method-specific override (only for return types, not parameters)
    if class_name and method_name and not param_name:
        override_key = (class_name, method_name)
        if override_key in method_type_overrides:
            return method_type_overrides[override_key]

    # First, try to parse the return description if available
    if return_desc:
        parsed = parse_return_description(return_desc)
        if parsed:
            return parsed

    if not doc_type or doc_type == "void":
        return "void"

    # Clean up the type string
    doc_type = doc_type.strip()

    # Normalize union types: "Type|undefined" -> "Type | undefined"
    doc_type = re.sub(r'\s*\|\s*', ' | ', doc_type)

    # If it's already a TypeScript array type (e.g., "number[]"), return as-is
    if doc_type.endswith("[]"):
        return doc_type

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
        result = type_map[doc_type]
    # Handle nested Array.<Type> notation (including HTML entities)
    # Convert Array.<Array.<number>> to number[][]
    elif "Array.<" in doc_type or "Array.&lt;" in doc_type:
        # Handle HTML entities
        doc_type = doc_type.replace("&lt;", "<").replace("&gt;", ">")

        # Recursively convert nested Array.<Type> to Type[]
        while "Array.<" in doc_type:
            # Find innermost Array.<Type>
            match = re.search(r'Array\.<([^<>]+)>', doc_type)
            if match:
                inner_type = match.group(1).strip()
                # Recursively convert the inner type
                converted_inner = convert_type_to_typescript(inner_type, class_name)
                # Replace this Array.<Type> with Type[]
                doc_type = doc_type[:match.start()] + converted_inner + "[]" + doc_type[match.end():]
            else:
                break

        result = doc_type
    # If it looks like a class name (starts with uppercase), keep it
    elif doc_type and doc_type[0].isupper():
        result = doc_type
    else:
        # Default to any for unknown types
        result = "any"

    # Check if the return description mentions that undefined can be returned
    # Pattern: "... or `undefined` ..." or "returns `undefined` if..."
    if return_desc and "`undefined`" in return_desc and result != "undefined":
        # Add undefined to the union type if not already void or undefined
        if result != "void":
            # Remove trailing space if present and add proper spacing
            result = result.rstrip() + " | undefined"

    return result


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

    # Add interface definitions for complex return types
    lines.append("/**")
    lines.append(" * Voice parameters object returned by NoteGroupReference.getVoice")
    lines.append(" */")
    lines.append("interface VoiceParameters {")
    lines.append("  paramLoudness: number;")
    lines.append("  paramTension: number;")
    lines.append("  paramBreathiness: number;")
    lines.append("  paramGender: number;")
    lines.append("  paramToneShift: number;")
    lines.append("  vocalModeParams: {")
    lines.append("    [vocalModeName: string]: {")
    lines.append("      pitch: number;")
    lines.append("      timbre: number;")
    lines.append("      pronunciation: number;")
    lines.append("    };")
    lines.append("  };")
    lines.append("}")
    lines.append("")

    lines.append("/**")
    lines.append(" * Computed attributes object returned by SV.getComputedAttributesForGroup")
    lines.append(" */")
    lines.append("interface ComputedAttributes {")
    lines.append("  accent: string;")
    lines.append("  rapTone: number | null;")
    lines.append("  rapIntonation: number | null;")
    lines.append("  phonemes: {")
    lines.append("    symbol: string;")
    lines.append("    language: string;")
    lines.append("    activity: number | null;")
    lines.append("    position: number | null;")
    lines.append("  }[];")
    lines.append("}")
    lines.append("")

    lines.append("/**")
    lines.append(" * Note attributes object returned by Note.getAttributes")
    lines.append(" */")
    lines.append("interface NoteAttributes {")
    lines.append("  rTone: number;")
    lines.append("  rIntonation: number;")
    lines.append("  dF0VbrMod: number;")
    lines.append("  expValueX: number;")
    lines.append("  expValueY: number;")
    lines.append("  phonemes: {;")
    lines.append("    leftOffset: number;")
    lines.append("    position: number;")
    lines.append("    activity: number;")
    lines.append("    strength: number;")
    lines.append("  }[];")
    lines.append("  muted: boolean;")
    lines.append("  evenSyllableDuration: boolean;")
    lines.append("  languageOverride: string;")
    lines.append("  phonesetOverride: string;")
    lines.append("}")

    lines.append("/**")
    lines.append(" * Host information object returned by SV.getHostInfo")
    lines.append(" */")
    lines.append("interface HostInfo {")
    lines.append("  osType: string;")
    lines.append("  osName: string;")
    lines.append("  hostName: string;")
    lines.append("  hostVersion: string;")
    lines.append("  hostVersionNumber: number;")
    lines.append("  languageCode: string;")
    lines.append("}")

    lines.append("/**")
    lines.append(" * TempoMark returned by TimeAxis")
    lines.append(" */")
    lines.append("interface TempoMark {")
    lines.append("  position: number;")
    lines.append("  positionSeconds: number;")
    lines.append("  bpm: number;")
    lines.append("}")

    lines.append("/**")
    lines.append(" * MeasureMark returned by TimeAxis")
    lines.append(" */")
    lines.append("interface MeasureMark {")
    lines.append("  position: number;")
    lines.append("  positionBlick: number;")
    lines.append("  numerator: number;")
    lines.append("  denominator: number;")
    lines.append("}")

    lines.append("/**")
    lines.append(" * Form is argument for showDialog methods")
    lines.append(" */")
    lines.append("interface Form {")
    lines.append("  title: string;")
    lines.append("  message: string;")
    lines.append('  buttons: "YesNoCancel"|"OkCancel";')
    lines.append("  widgets: Widget[];")
    lines.append("}")

    lines.append("/**")
    lines.append(" * Widget is definition of a form element")
    lines.append(" */")
    lines.append("type Widget = Slider | CheckBox | ComboBox | TextBox | TextArea;")

    lines.append('interface Slider {')
    lines.append('    name: string;')
    lines.append('    type: "Slider";')
    lines.append('    label: string;')
    lines.append('    format: string;')
    lines.append('    minValue: number;')
    lines.append('    maxValue: number;')
    lines.append('    interval: number;')
    lines.append('    default: number;')
    lines.append('}')
    lines.append('')
    lines.append('interface CheckBox {')
    lines.append('    name: string;')
    lines.append('    type: "CheckBox";')
    lines.append('    text: string;')
    lines.append('    default: boolean;')
    lines.append('}')
    lines.append('')
    lines.append('interface ComboBox {')
    lines.append('    name: string;')
    lines.append('    type: "ComboBox";')
    lines.append('    label: string;')
    lines.append('    choices: string[];')
    lines.append('    default: number;')
    lines.append('}')
    lines.append('')
    lines.append('interface TextBox { ')
    lines.append('    name: string;')
    lines.append('    type: "TextBox";')
    lines.append('    label: string;')
    lines.append('    default: string;')
    lines.append('}')
    lines.append('')
    lines.append('interface TextArea {')
    lines.append('    name: string;')
    lines.append('    type: "TextArea";')
    lines.append('    label: string;')
    lines.append('    height: number;')
    lines.append('    default: string;')
    lines.append('}')
    lines.append('')

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
                    ts_type = convert_type_to_typescript(param_type, class_info.name, "", method.name, param_name)
                    if param_desc:
                        lines.append(f"   * @param {param_name} {param_desc}")
                    else:
                        lines.append(f"   * @param {param_name}")

                # Return documentation
                if method.return_type != "void":
                    ts_return = convert_type_to_typescript(method.return_type, class_info.name, method.return_desc, method.name)
                    lines.append(f"   * @returns {ts_return}")

                lines.append("   */")

            # Generate signature based on whether it's a property or method
            static_keyword = "static " if method.is_static else ""

            if method.is_property:
                # Property syntax: static readonly PROPERTY: type;
                readonly_keyword = "readonly " if static_keyword else ""
                return_type = convert_type_to_typescript(method.return_type, class_info.name, method.return_desc, method.name)
                lines.append(f"  {static_keyword}{readonly_keyword}{method.name}: {return_type};")
            else:
                # Method syntax: methodName(params): returnType;
                params_str = ", ".join([
                    f"{name}: {convert_type_to_typescript(ptype, class_info.name, '', method.name, name)}"
                    for name, ptype, _ in method.params
                ])
                return_type = convert_type_to_typescript(method.return_type, class_info.name, method.return_desc, method.name)
                lines.append(f"  {static_keyword}{method.name}({params_str}): {return_type};")

            lines.append("")

        lines.append("}")
        lines.append("")

    return "\n".join(lines)


def main():
    """Main function to generate TypeScript definitions."""
    # Input and output paths
    docs_dir = Path("dreamtonics-api")
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
