#!/usr/bin/env python3
"""Smoke test for the pybind11 idgen binding.

Run after building with ./python/build.sh:
    python3 python/smoke_test.py
"""

from __future__ import annotations

import importlib
import pathlib
import re
import sys
import uuid

ALLOWED = set("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")


def _fail(msg: str) -> None:
    print(f"FAIL: {msg}")
    sys.exit(1)


def _check_prefix(value: str, expected_prefix: str) -> str:
    if ":" not in value:
        _fail(f"missing ':' separator in {value!r}")

    prefix, suffix = value.split(":", 1)
    if prefix != expected_prefix:
        _fail(f"prefix mismatch: expected {expected_prefix!r}, got {prefix!r}")
    if not suffix:
        _fail("suffix is empty")
    return suffix


def main() -> None:
    root = pathlib.Path(__file__).resolve().parent
    build_dir = root / "build"
    if not build_dir.exists():
        _fail("python/build not found; run ./python/build.sh first")

    sys.path.insert(0, str(build_dir))

    try:
        mod = importlib.import_module("_idgen")
    except Exception as exc:  # pragma: no cover - smoke path
        _fail(f"unable to import _idgen from {build_dir}: {exc}")

    short_id = mod.doc_id("Invoice Data", "short")
    short_suffix = _check_prefix(short_id, "invoice_data")
    if len(short_suffix) != 10:
        _fail(f"short format length mismatch: expected 10, got {len(short_suffix)}")
    if any(ch not in ALLOWED for ch in short_suffix):
        _fail(f"short format contains invalid characters: {short_suffix!r}")

    uuid_id = mod.doc_id("order", "uuid")
    uuid_suffix = _check_prefix(uuid_id, "order")
    try:
        uuid.UUID(uuid_suffix)
    except ValueError as exc:
        _fail(f"uuid format invalid: {exc}")

    ts_id = mod.doc_id("report", "timestamp")
    ts_suffix = _check_prefix(ts_id, "report")
    if not re.fullmatch(r"\d+", ts_suffix):
        _fail(f"timestamp format should be numeric, got {ts_suffix!r}")

    print("PASS: pybind11 smoke test")


if __name__ == "__main__":
    main()
