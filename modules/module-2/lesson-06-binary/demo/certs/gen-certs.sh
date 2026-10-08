#!/usr/bin/env bash
# Generates a demo CA and a TLS certificate for harbor.local, then regenerates
# the Kubernetes secret manifest (../manifests/harbor-tls-secret.yaml).
#
# Run from anywhere:  bash certs/gen-certs.sh
#
# The generated CA (ca.crt) must be trusted by the docker daemon:
#   sudo mkdir -p /etc/docker/certs.d/harbor.local
#   sudo cp certs/ca.crt /etc/docker/certs.d/harbor.local/ca.crt
set -euo pipefail
cd "$(dirname "$0")"

DOMAIN="${1:-harbor.local}"
CA_DAYS="${CA_DAYS:-3650}"      # 10 years, demo CA
CERT_DAYS="${CERT_DAYS:-825}"   # max accepted by modern clients

# --- CA ---------------------------------------------------------------------
openssl req -x509 -newkey rsa:2048 -nodes -days "$CA_DAYS" \
  -subj "/CN=Harbor Demo CA" \
  -keyout ca.key -out ca.crt

# --- server certificate for $DOMAIN ------------------------------------------
openssl req -newkey rsa:2048 -nodes \
  -subj "/CN=${DOMAIN}" \
  -keyout "${DOMAIN}.key" -out "${DOMAIN}.csr" \
  -addext "subjectAltName=DNS:${DOMAIN}"

printf "subjectAltName=DNS:%s\n" "$DOMAIN" > san.ext
openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial \
  -in "${DOMAIN}.csr" -out "${DOMAIN}.crt" \
  -days "$CERT_DAYS" -extfile san.ext

# --- Kubernetes secret manifest ----------------------------------------------
kubectl create secret tls harbor-tls -n harbor \
  --cert="${DOMAIN}.crt" --key="${DOMAIN}.key" \
  --dry-run=client -o yaml > ../manifests/harbor-tls-secret.yaml

echo "OK: certs/ca.crt (CA), ${DOMAIN}.crt/.key (server)"
echo "K8s secret manifest updated: manifests/harbor-tls-secret.yaml"
