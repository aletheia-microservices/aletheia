#!/bin/bash

# Usage:
#   ./validate.sh <appname> [schema|analysis]
#
# Examples:
#   ./validate.sh train_ticket2 schema
#   ./validate.sh train_ticket2 analysis

if [ "$#" -ne 2 ]; then
  echo "usage: $0 <appname> [schema|analysis]"
  echo "available apps:"
  echo "- foobar"
  echo "- digota"
  echo "- sockshop3"
  echo "- eshopmicroservices"
  echo "- postnotification"
  echo "- dsb_sn2"
  echo "- dsb_mediamicroservices"
  echo "- train_ticket2"
  echo "- dsb_hotel2"
  exit 1
fi

appname=$1
type=$2
basedir="."
outputdir="${basedir}/ssa_analysis/analyzer/output/${appname}"

if [ "$type" = "analysis" ]; then
  echo
  echo "Select analysis kind:"
  echo "  1) foreign-key-cascade"
  echo "  2) foreign-key-concurrency"
  echo "  3) foreign-key-coordination"
  echo "  4) primary-key-coordination"
  echo "  5) unicity-concurrency"
  echo

  read -rp "Enter a number [1-5]: " choice
  case "$choice" in
    1) kind="foreign-key-cascade" ;;
    2) kind="foreign-key-concurrency" ;;
    3) kind="foreign-key-coordination" ;;
    4) kind="primary-key-coordination" ;;
    5) kind="unicity-concurrency" ;;
    *)
      echo "Invalid choice."
      exit 1
      ;;
  esac
else
  kind=""
fi

if [ "$type" = "schema" ]; then
  expected="${outputdir}/expected-schema.json"
  actual="${outputdir}/schema.json"
else
  expected="${outputdir}/analysis/expected-${kind}.txt"
  actual="${outputdir}/analysis/${kind}.txt"
fi

if [ ! -f "$expected" ]; then
  echo "error: expected file not found: $expected"
  exit 1
fi

if [ ! -f "$actual" ]; then
  echo "error: actual file not found: $actual"
  exit 1
fi

# ignore any comma at the end of each line 
# important when results change and item is not last in array and loses the comma
jsondiff() {
  diff --color=always -uw <(sed 's/,\s*$//' "$1") <(sed 's/,\s*$//' "$2")
}

# ignore any enumeration #NUMBER_WARNINGS
# important when results change and order changes 
txtdiff() {
  diff --color=always -uw \
    <(sed -E 's/#([0-9]+)/#?/g; s/,\s*$//' "$1") \
    <(sed -E 's/#([0-9]+)/#?/g; s/,\s*$//' "$2")
}

if [ "$type" = "schema" ]; then
  jsondiff "$expected" "$actual"
else
  txtdiff "$expected" "$actual"
fi
