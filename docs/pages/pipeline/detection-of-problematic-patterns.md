# Pattern Detection

## Overview

```mermaid
--8<-- "docs/snippets/pipeline.mmd"
    class detect current
    class ssa,acg,schema step
```

The final step traverses the [abstract call graph](inter-microservice-analysis.md) of each request, locating the database operations and checking them against the five patterns described in [What Aletheia Detects](../patterns.md), based on the constraints from [schema extraction](schema-extraction.md). Aletheia reports a warning for each match.

```mermaid
flowchart LR
    acg["Abstract call graph"]
    schemas["Database schemas<br/>with their keys"]
    detectors["Detectors<br/>RI-1, RI-2, RI-3, EI-1, Un-1"]
    warnings["Warnings<br/>output/{app}/analysis/"]
    acg --> detectors
    schemas --> detectors
    detectors --> warnings
    classDef step fill:#268bd214,stroke:#268bd280
    classDef io fill:#00000000,stroke:#80808099,stroke-dasharray:4 3
    class detectors step
    class acg,schemas,warnings io
```

See [Usage → Read Detection Warnings](../usage/reading-warnings.md) for how to interpret the warnings.
