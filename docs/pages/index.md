# Aletheia

Aletheia's source code is available on [GitHub](https://github.com/aletheia-microservices/aletheia).

## Overview

Aletheia uses static analysis to find service interactions and database operations that can break data integrity constraints (see [What Aletheia Detects](patterns.md)):

- Entity integrity
- Referential integrity
- Uniqueness

Aletheia analyzes applications targeting the [Blueprint](https://github.com/Blueprint-uServices/blueprint) compiler.

The code of all applications is located in `blueprint/examples/`, which includes a README summarizing the source repositories and versions used for each application.

The analysis has four steps:

1. [Intra-procedural analysis](pipeline/intra-procedural-analysis.md) tracks values within each function using Static Single-Assignment (SSA) graphs and taint propagation.
2. [Inter-microservice analysis](pipeline/inter-microservice-analysis.md) connects these values across services through an _abstract call graph_ of service calls and database operations.
3. [Schema extraction](pipeline/schema-extraction.md) builds the database schemas and infers relationships between fields stored in different databases.
4. [Pattern detection](pipeline/detection-of-problematic-patterns.md) checks the operations against patterns that can violate the extracted constraints.

```mermaid
--8<-- "docs/snippets/pipeline.mmd"
    class ssa,acg,schema,detect step
```

The detection step searches for five patterns of operations that break these constraints, described in [What Aletheia Detects](patterns.md). To build Aletheia and analyze your first application, see [Getting Started → Installation](getting-started/installation.md).

## Citation

If you use Aletheia in your work, please cite the paper:


```bibtex title="aletheia-osdi26.bib"
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
