#!/bin/bash
set -euo pipefail

# Keycloak n'accepte pas de joker sur l'hote des redirect URIs : le domaine
# reseau (sslip.io) est donc injecte dans le realm avant l'import.
domaine="${HKDOMAINERESEAU:-historykanban.127.0.0.1.sslip.io}"
realm="${KEYCLOAKREALM:-historykanban}"
client="${KEYCLOAKCLIENTINTERFACE:-interface}"
mkdir -p /opt/keycloak/data/import
for modele in /opt/keycloak/realm-modeles/*.json; do
  contenu=$(<"$modele")
  printf '%s' "${contenu//__DOMAINERESEAU__/$domaine}" > "/opt/keycloak/data/import/${modele##*/}"
done

synchroniser_redirects() {
  if [[ -z "${KC_BOOTSTRAP_ADMIN_USERNAME:-}" || -z "${KC_BOOTSTRAP_ADMIN_PASSWORD:-}" ]]; then
    echo "Synchronisation des redirect URIs ignoree : administrateur Keycloak absent."
    return 0
  fi

  for tentative in $(seq 1 60); do
    if /opt/keycloak/bin/kcadm.sh config credentials \
      --server http://localhost:8080 \
      --realm master \
      --user "$KC_BOOTSTRAP_ADMIN_USERNAME" \
      --password "$KC_BOOTSTRAP_ADMIN_PASSWORD" >/dev/null 2>&1; then
      break
    fi
    if [[ "$tentative" -eq 60 ]]; then
      echo "Synchronisation des redirect URIs impossible : Keycloak non pret."
      return 0
    fi
    sleep 2
  done

  identifiant=$(/opt/keycloak/bin/kcadm.sh get clients -r "$realm" -q clientId="$client" --fields id --format csv --noquotes 2>/dev/null | tail -n 1 | tr -d '\r')

  if [[ -z "$identifiant" ]]; then
    echo "Synchronisation des redirect URIs ignoree : client $client introuvable."
    return 0
  fi

  cat > /tmp/client-interface-redirects.json <<JSON
{
  "redirectUris": [
    "https://historykanban.localhost/*",
    "http://localhost:5173/*",
    "https://$domaine/*"
  ],
  "webOrigins": ["+"],
  "attributes": {
    "pkce.code.challenge.method": "S256",
    "post.logout.redirect.uris": "+"
  }
}
JSON

  /opt/keycloak/bin/kcadm.sh update "clients/$identifiant" -r "$realm" -f /tmp/client-interface-redirects.json
  echo "Redirect URIs Keycloak synchronisees pour https://$domaine/*"
}

/opt/keycloak/bin/kc.sh "$@" &
pid_keycloak=$!
trap 'kill -TERM "$pid_keycloak" 2>/dev/null || true; wait "$pid_keycloak"' TERM INT

synchroniser_redirects &
wait "$pid_keycloak"
