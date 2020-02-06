#!/bin/bash

# Script to validate a YAML file
python3 -m pip install --user yamllint > /dev/null
python3 check_lint.py $1

