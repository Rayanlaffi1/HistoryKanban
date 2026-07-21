#!/bin/sh
set -eu

force=false
if [ "${1:-}" = "--force" ]; then
  force=true
elif [ "$#" -gt 0 ]; then
  echo "Usage: $0 [--force]" >&2
  exit 2
fi

command -v openssl >/dev/null 2>&1 || {
  echo "OpenSSL est requis." >&2
  exit 1
}

racine=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
fichier_env="$racine/.env"
dossier_certificats="$racine/traefik/certificats"
fichier_cle="$dossier_certificats/historykanban.key"
fichier_certificat="$dossier_certificats/historykanban.crt"

if [ "$force" = false ] && { [ -e "$fichier_env" ] || [ -e "$fichier_cle" ] || [ -e "$fichier_certificat" ]; }; then
  echo "Des secrets existent deja. Relancez avec --force pour les renouveler." >&2
  exit 1
fi

dossier_temporaire=$(mktemp -d)
case "$dossier_temporaire" in
  ""|/|.) echo "Dossier temporaire invalide." >&2; exit 1 ;;
esac
trap 'rm -rf -- "$dossier_temporaire"' EXIT HUP INT TERM

umask 077
cat >"$dossier_temporaire/.env" <<EOF
BDMDP=$(openssl rand -hex 24)
BDNOM=historykanban
BDUTILISATEUR=historykanban
KEYCLOAKADMIN=admin
KEYCLOAKADMINMDP=$(openssl rand -hex 24)
MINIOCLE=$(openssl rand -hex 10)
MINIOSECRET=$(openssl rand -hex 20)
RABBITUTILISATEUR=historykanban
RABBITMDP=$(openssl rand -hex 24)
SAUVEGARDEINTERVALLE=86400
SAUVEGARDERETENTION=14
EOF

openssl req -x509 -newkey rsa:4096 -sha256 -days 825 -nodes \
  -subj '/CN=historykanban.localhost' \
  -addext 'subjectAltName=DNS:historykanban.localhost,DNS:*.historykanban.localhost,DNS:localhost,IP:127.0.0.1,IP:::1' \
  -keyout "$dossier_temporaire/historykanban.key" \
  -out "$dossier_temporaire/historykanban.crt"

mkdir -p -- "$dossier_certificats"
install -m 600 "$dossier_temporaire/.env" "$fichier_env"
install -m 600 "$dossier_temporaire/historykanban.key" "$fichier_cle"
install -m 644 "$dossier_temporaire/historykanban.crt" "$fichier_certificat"

echo "Secrets et certificat locaux renouveles. Aucune valeur sensible n'a ete affichee."
