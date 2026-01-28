#! /usr/bin/env bash

go build -o pilot main.go
sudo mv pilot /usr/local/bin

COMPLETION_DEST="/usr/local/share/bash-completion/completions"
sudo mkdir -p "$COMPLETION_DEST"

/usr/local/bin/pilot completion bash | sudo tee "$COMPLETION_DEST/pilot" > /dev/null

sudo chmod 644 "$COMPLETION_DEST/pilot"

echo "Successfully installed pilot and completions."
