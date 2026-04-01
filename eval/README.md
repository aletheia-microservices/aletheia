# Aletheia

## Evaluation

Configure temporary and cache directories for Go necessary for the largest synthetic applications:

```zsh
mkdir -p $HOME/tmp-go
mkdir -p $HOME/go-big-cache

TMPDIR=$HOME/tmp-go
GOCACHE=$HOME/go-big-cache
GOGC=20
GOFLAGS="-p=1"
```

Generate the synthetic applications based on graph graph characteristics (call depth, fan-out, and request volume) defined in `blueprint/examples/scripts/config.yaml` whose values were derived from [Alibaba's production microservice traces](https://dl.acm.org/doi/10.1145/3472883.3487003):

```zsh
cd ~/aletheia/eval/gen_synthetic/
go run main.go
```

Run the experiments for **realistic applications** (digota, sockshop, eshopmicroservices, postnotification, dsb_socialnetwork, dsb_mediamicroservices, trainticket) and **synthetic applications** (app1, app2, app3, app4, app5):

```zsh
cd ~/aletheia/eval/
./run.sh --cache
./run.sh --eval
```

[Optional] Configure Python environment to collect results:
```zsh
cd ~/aletheia/eval/
python3 -m venv ~/.venv
source ~/.venv/bin/activate
pip3 install -r requirements.txt
```

Collect and summarize average **memory** results for all experiments ran on the current date:

```zsh
python3 collect_memory.py
python3 collect_memory.py --synthetic
```

Collect and plot average **time** results for all experiments ran on the current date:

```zsh
python3 collect_metrics.py
python3 collect_metrics.py --synthetic
python3 plot_times.py
```
