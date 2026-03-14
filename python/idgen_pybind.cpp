#include <pybind11/pybind11.h>

extern "C" {
char* DocID(const char* prefix, const char* format);
void FreeCString(char* s);
}

namespace py = pybind11;

std::string py_doc_id(const std::string& prefix, const std::string& format) {
  char* raw = DocID(prefix.c_str(), format.c_str());
  if (raw == nullptr) {
    throw std::runtime_error("DocID returned null");
  }

  std::string out(raw);
  FreeCString(raw);
  return out;
}

PYBIND11_MODULE(_idgen, m) {
  m.doc() = "Python bindings for github.com/m-theory-io/idgen";
  m.def("doc_id", &py_doc_id, py::arg("prefix") = "", py::arg("format") = "short",
        "Generate a document ID using the Go idgen library.");
}
