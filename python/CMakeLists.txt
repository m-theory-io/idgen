cmake_minimum_required(VERSION 3.18)
project(idgen_py LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(pybind11 REQUIRED)

if(NOT DEFINED IDGEN_CGO_INCLUDE_DIR)
  message(FATAL_ERROR "IDGEN_CGO_INCLUDE_DIR is required (path containing libidgen.h)")
endif()

if(NOT DEFINED IDGEN_CGO_LIBRARY)
  message(FATAL_ERROR "IDGEN_CGO_LIBRARY is required (full path to libidgen.{dylib,so})")
endif()

get_filename_component(IDGEN_CGO_LIBRARY_DIR ${IDGEN_CGO_LIBRARY} DIRECTORY)

pybind11_add_module(_idgen idgen_pybind.cpp)
target_include_directories(_idgen PRIVATE ${IDGEN_CGO_INCLUDE_DIR})
target_link_libraries(_idgen PRIVATE ${IDGEN_CGO_LIBRARY})

# Ensure Python can load the Go shared library without manual DYLD_* setup.
set_target_properties(_idgen PROPERTIES
  BUILD_RPATH "${IDGEN_CGO_LIBRARY_DIR}"
  INSTALL_RPATH "${IDGEN_CGO_LIBRARY_DIR}"
)
