#!/bin/bash
set -euo pipefail

# Keycloak n'accepte pas de joker sur l'hote des redirect URIs : le domaine
# reseau (sslip.io) est donc injecte dans le realm avant l'import.
domaine="${HKDOMAINERESEAU:-historykanban.127.0.0.1.sslip.io}"
mkdir -p /opt/keycloak/data/import
for modele in /opt/keycloak/realm-modeles/*.json; do
  contenu=$(<"$modele")
  printf '%s' "${contenu//__DOMAINERESEAU__/$domaine}" > "/opt/keycloak/data/import/${modele##*/}"
done

exec /opt/keycloak/bin/kc.sh "$@"
