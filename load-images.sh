#!/usr/bin/env bash

function load { docker load <$(nix build .#$1 --no-link --print-out-paths) | cut -d' ' -f 3; }

IMAGE_DB=$(load db)
IMAGE_BACKEND=$(load backend)
IMAGE_FRONTEND=$(load frontend)

docker compose up -d
