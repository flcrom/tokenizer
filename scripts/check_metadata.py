#!/usr/bin/env python3
"""
@file scripts/check_metadata.py
@description Validates required metadata, repository-relative paths, and reciprocal dependency edges.
@depends_on MEMORY.md
@depended_on_by AGENTS.md
"""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent
IGNORED_DIRECTORIES = {
    ".git",
    ".local",
    ".venv",
    "build",
    "coverage",
    "dist",
    "node_modules",
    "playwright-report",
    "test-results",
    "vendor",
}
IGNORED_FILENAMES = {".DS_Store", "package-lock.json", "pnpm-lock.yaml", "yarn.lock"}
IGNORED_SUFFIXES = {
    ".avif",
    ".gif",
    ".ico",
    ".jpeg",
    ".jpg",
    ".json",
    ".lock",
    ".pdf",
    ".png",
    ".pyc",
    ".svg",
    ".webp",
    ".zip",
}
REQUIRED_FIELDS = ("file", "description", "depends_on", "depended_on_by")
FIELD_PATTERN = re.compile(
    r"@(?P<field>file|description|depends_on|depended_on_by)\s+(?P<value>.+?)\s*(?:\*/|-->)?$"
)


def candidate_files() -> list[Path]:
    candidates: list[Path] = []
    for path in ROOT.rglob("*"):
        if not path.is_file():
            continue
        relative = path.relative_to(ROOT)
        if any(part in IGNORED_DIRECTORIES for part in relative.parts):
            continue
        if path.name in IGNORED_FILENAMES or path.suffix.lower() in IGNORED_SUFFIXES:
            continue
        candidates.append(path)
    return sorted(candidates)


def read_header(path: Path) -> dict[str, str]:
    try:
        lines = path.read_text(encoding="utf-8").splitlines()[:10]
    except UnicodeDecodeError:
        return {}

    metadata: dict[str, str] = {}
    for line in lines:
        match = FIELD_PATTERN.search(line.strip())
        if match:
            metadata[match.group("field")] = match.group("value").strip()
    return metadata


def parse_paths(value: str) -> set[str]:
    if value.lower() == "none":
        return set()
    return {item.strip() for item in value.split(",") if item.strip()}


def main() -> int:
    errors: list[str] = []
    records: dict[str, dict[str, str]] = {}

    for path in candidate_files():
        relative = path.relative_to(ROOT).as_posix()
        metadata = read_header(path)
        missing = [field for field in REQUIRED_FIELDS if field not in metadata]
        if missing:
            errors.append(f"{relative}: missing metadata field(s): {', '.join(missing)}")
            continue
        if metadata["file"] != relative:
            errors.append(
                f"{relative}: @file declares {metadata['file']!r}; expected {relative!r}"
            )
        records[relative] = metadata

    for relative, metadata in records.items():
        dependencies = parse_paths(metadata["depends_on"])
        consumers = parse_paths(metadata["depended_on_by"])

        for dependency in dependencies:
            target = records.get(dependency)
            if target is None:
                errors.append(f"{relative}: dependency does not exist or lacks metadata: {dependency}")
                continue
            reverse_edges = parse_paths(target["depended_on_by"])
            if relative not in reverse_edges:
                errors.append(
                    f"{relative}: depends on {dependency}, but its @depended_on_by omits {relative}"
                )

        for consumer in consumers:
            target = records.get(consumer)
            if target is None:
                errors.append(f"{relative}: consumer does not exist or lacks metadata: {consumer}")
                continue
            forward_edges = parse_paths(target["depends_on"])
            if relative not in forward_edges:
                errors.append(
                    f"{relative}: lists consumer {consumer}, but its @depends_on omits {relative}"
                )

    if errors:
        print("Metadata validation failed:")
        for error in errors:
            print(f"- {error}")
        return 1

    print(f"Metadata validation passed for {len(records)} files.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
