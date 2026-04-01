# Aletheia

## Getting Started

After cloning the repository, fetch submodules for `blueprint` repository

```zsh
cd aletheia
git submodule update --init --recursive
```

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.24.5

## Analyzing Applications

By default, Aletheia supports analysis for the following applications
- digota
- sockshop
- eshopmicroservices
- postnotification
- dsb_socialnetwork
- dsb_mediamicroservices
- trainticket

If you want to analyze additional applications, follow these steps:


## Running Aletheia

Run Aletheia to analyze the application specified by the `app` parameter:
```zsh
cd ~/aletheia/ssa-analysis/analyzer/
go run main.go {app}
```
