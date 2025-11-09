# Microservices Consistency Analyzer (MSCA)

## Getting Started

After cloning the repository, fetch submodules for `blueprint` repository

```zsh
cd ms-consistency-analyzer
git submodule update --init --recursive
```

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.22.2

## Prepare Environment

```zsh
cd ms-consistency-anayzer
vagrant up
vagrant ssh
```

## Running the Tool

```zsh
cd ssa_analysis/analyzer/

# testing purposes
go run main.go foobar
# e-commerce
go run main.go digota
go run main.go sockshop3
# socialnetwork
go run main.go postnotification_simple
go run main.go dsb_sn2
# media
go run main.go dsb_media_sql
# reservation
go run main.go dsb_hotel2
go run main.go train_ticket2
# eval
go run main.go large_scale_app
```
