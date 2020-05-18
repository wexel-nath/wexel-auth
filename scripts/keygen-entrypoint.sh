#!/bin/sh
set -e

PASSPHRASE="${PASSPHRASE:-}"
PRIVATE_KEY_PATH="${PRIVATE_KEY_PATH:-/keys/private.pem}"
PUBLIC_KEY_PATH="${PUBLIC_KEY_PATH:-/keys/public.pem}"

# delete old keys
rm -rf /keys/*

# generate private key
ssh-keygen \
	-t rsa \
	-b 2048 \
	-m PEM \
	-N "$PASSPHRASE" \
	-f "$PRIVATE_KEY_PATH" \

# generate public key
openssl rsa \
	-in "$PRIVATE_KEY_PATH" \
	-pubout \
	-outform PEM \
	-out "$PUBLIC_KEY_PATH"
