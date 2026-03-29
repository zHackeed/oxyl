#!/bin/sh¡

set -e

#This script inits the required data for oxyl to work.

# Only init the keys if they are not already present.
if [ ! -f /data/keys/ed25519-priv.pem ]; then
    mkdir -p /data/keys
    openssl genpkey -algorithm ED25519 -out /data/keys/ed25519-priv.pem
    openssl pkey -in /data/keys/ed25519-priv.pem -pubout -out /data/keys/ed25519-pub.pem
fi

# Todo: might have to create more delegations for the agent backend logic?