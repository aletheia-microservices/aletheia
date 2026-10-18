# Schema Diagram

Draws the data schema of an app from the `output/{app}/schema.json` that Aletheia writes when it analyzes the app. From the repository root:

```zsh
go run ./scripts/schema-diagram trainticket   # writes output/trainticket/diagrams/schema.png
```

See [Generating Diagrams](https://aletheia-microservices.github.io/aletheia/3-usage/diagrams.html#schema-diagram) for an example diagram, the options, and how to update the diagrams in [examples](examples/).
