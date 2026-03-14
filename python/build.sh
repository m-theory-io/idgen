#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="${ROOT_DIR}/python/build"
CGO_OUT_DIR="${ROOT_DIR}/python/.cgo"

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required" >&2
  exit 1
fi

if ! command -v cmake >/dev/null 2>&1; then
  echo "cmake is required" >&2
  exit 1
fi

if ! PYBIND11_CMAKE_DIR="$(python3 -m pybind11 --cmakedir 2>/dev/null)"; then
  echo "pybind11 is required. Install with: python3 -m pip install pybind11" >&2
  exit 1
fi

mkdir -p "${BUILD_DIR}" "${CGO_OUT_DIR}"

echo "[1/2] Building Go shared library..."
cd "${ROOT_DIR}"
go build -buildmode=c-shared -o "${CGO_OUT_DIR}/libidgen.dylib" ./cgo

if [[ "$(uname -s)" == "Darwin" ]]; then
  install_name_tool -id "@rpath/libidgen.dylib" "${CGO_OUT_DIR}/libidgen.dylib"
fi

echo "[2/2] Building pybind11 module..."
cd "${BUILD_DIR}"
cmake .. \
	-Dpybind11_DIR="${PYBIND11_CMAKE_DIR}" \
  -DIDGEN_CGO_INCLUDE_DIR="${CGO_OUT_DIR}" \
  -DIDGEN_CGO_LIBRARY="${CGO_OUT_DIR}/libidgen.dylib"
cmake --build .

echo "Build complete. Python extension is in: ${BUILD_DIR}"
