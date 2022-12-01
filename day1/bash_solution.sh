#!/bin/bash
set -euo pipefail

elfCalories=()

function readInput() {
  local idx=0
  local elfCal=0

  while read LINE; do
    if [[ -n $LINE ]]; then
      elfCal=$((elfCal + LINE))
    else
      elfCalories[$idx]=$elfCal
      idx=$((idx + 1))
      elfCal=0
    fi
  done

  elfCalories[$idx]=$elfCal
}

function solve1() {
  local max=0
  for cal in ${elfCalories[@]}; do
    [ $cal -gt $max ] && max=$cal || continue
  done

  echo "Solution 1: $max"
}

function solve2() {
  IFS=$'\n'

  local sum=0
  for cal in $(echo "${elfCalories[*]}" | sort -nr | head -n3); do
    sum=$((sum + cal))
  done

  echo "Solution 2: $sum"
}

function main() {
  readInput
  solve1
  solve2
}

main
