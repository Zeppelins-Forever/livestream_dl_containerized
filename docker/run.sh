#!/bin/sh

# Puts our environment into our Python virtual environment (venv):
source myenv/bin/activate

# Potential alternate command. Currently, it just outputs videos as "[date] video name (ID).mp4.mp4". 
# Not ideal, but keeping it around for reference.
# python3 runner.py --log-level DEBUG --output "/out/[%(upload_date)s] %(title)s (%(id)s).%(ext)s"  --resolution best $1

# This runs CanOfSocks's livestream_dl program.
python3 runner.py --log-level DEBUG --output "/out/[%(upload_date)s] %(title)s (%(id)s)" --resolution best $1
