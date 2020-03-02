#!/bin/bash

export BASEURI=http://localhost:38413

curl -o messages.template.json $BASEURI/download
