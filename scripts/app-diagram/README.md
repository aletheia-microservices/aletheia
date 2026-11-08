# App Diagram

Draws the service dependency graph of an app from the `output/{app}/app.json` that Aletheia writes when it analyzes the app. From the repository root:

```zsh
go run ./scripts/app-diagram trainticket   # writes output/trainticket/diagrams/app.png
```

See [Generate Diagrams](https://aletheia-microservices.github.io/aletheia/usage/diagrams/#app-diagram) for an example diagram, the options, and how to update the diagrams in [examples](examples/).
