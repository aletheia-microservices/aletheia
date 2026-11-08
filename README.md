# Aletheia

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/aletheia-microservices/aletheia/blob/main/LICENSE)
[![Documentation](https://img.shields.io/badge/Docs-site-green.svg)](https://aletheia-microservices.github.io/aletheia/)
[![Go](https://img.shields.io/github/go-mod/go-version/aletheia-microservices/aletheia)](https://github.com/aletheia-microservices/aletheia/blob/main/go.mod)
[![USENIX OSDI 2026](https://img.shields.io/badge/USENIX-OSDI_%2726-green)](https://www.usenix.org/conference/osdi26/presentation/ferreira)
[![Paper](https://img.shields.io/badge/Paper-PDF-b31b1b)](https://www.usenix.org/system/files/osdi26-ferreira.pdf)
[![Cite](https://img.shields.io/badge/Cite-BibTeX-6c757d)](#citation)

### Overview

Aletheia uses static analysis to find service interactions and database operations in microservice applications that break data integrity, for entity integrity (enforced through primary keys), referential integrity (enforced through foreign keys), and uniqueness constraints. It analyzes applications written with the [Blueprint](https://github.com/Blueprint-uServices/blueprint) compiler.

### Quick Start

Requires [Golang](https://go.dev/doc/install) >= 1.26. Clone the repository, build Aletheia, and analyze one of the included applications:

```zsh
git clone --recurse-submodules https://github.com/aletheia-microservices/aletheia.git
cd aletheia
make build
./bin/aletheia postnotification
```

The warnings are printed and saved in `output/{app}/analysis/`.

### Documentation

The documentation is available at [https://aletheia-microservices.github.io/aletheia](https://aletheia-microservices.github.io/aletheia/).

### Experiments

For instructions on how to reproduce the experiments from the paper, see the [Aletheia Artifact OSDI 2026](https://github.com/aletheia-microservices/aletheia-artifact-osdi26) repository.

### Citation

If you use Aletheia in your work, please cite the paper:

```bibtex
@inproceedings{aletheia-osdi26,
  author    = {Mafalda Sofia Ferreira and Jo{\~a}o Ferreira Loff and Jo{\~a}o Garcia and Rodrigo Rodrigues},
  title     = {Aletheia: Automated Detection of Data Integrity Violations in Microservices},
  booktitle = {Proceedings of the 20th USENIX Symposium on Operating Systems Design and Implementation (OSDI 26)},
  year      = {2026},
  isbn      = {978-1-939133-55-7},
  address   = {Seattle, WA},
  pages     = {721--737},
  url       = {https://www.usenix.org/conference/osdi26/presentation/ferreira},
  publisher = {USENIX Association},
  month     = jul
}
```
