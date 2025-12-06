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
cd ssa-analysis/analyzer/

# e-commerce
go run main.go digota
go run main.go sockshop3
# socialnetwork
go run main.go postnotification_simple
go run main.go dsb_sn2
# media
go run main.go dsb_media_nosql
# reservation
go run main.go train_ticket2

# large apps
go run main.go synthetic_app

# extras
go run main.go foobar
go run main.go dsb_media_sql
go run main.go dsb_hotel2
```

## Evaluation

```zsh
cd ms-consistency-analyzer/blueprint/examples/scripts/
go run main.go

cd ms-consistency-analyzer/ssa-analysis/analyzer/
./run.sh --eval

cd ms-consistency-analyzer/eval
python3 -m venv ~/.venv
source ~/.venv/bin/activate
pip3 install -r requirements.txt

python3 average_times.py
python3 average_times.py --synthetic
python3 plot.py
python3 plot.py --synthetic

python3 average_metrics.py
python3 average_metrics.py --synthetic
```
