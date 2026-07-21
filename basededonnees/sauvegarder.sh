#!/bin/sh
set -u

DESTINATION=/sauvegardes
INTERVALLE=${INTERVALLE:-86400}
RETENTION=${RETENTION:-14}

mkdir -p "$DESTINATION"

journaliser() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') sauvegarde : $1"
}

sauvegarder() {
    horodatage=$(date '+%Y%m%d-%H%M%S')
    fichier="$DESTINATION/historykanban-$horodatage.sql.gz"
    temporaire="$fichier.encours"
    if pg_dumpall -h "$BDHOTE" -U "$BDUTILISATEUR" | gzip -6 > "$temporaire"; then
        mv "$temporaire" "$fichier"
        journaliser "$fichier ($(du -h "$fichier" | cut -f1))"
    else
        rm -f "$temporaire"
        journaliser "echec de la sauvegarde"
        return 1
    fi
}

purger() {
    supprimees=$(find "$DESTINATION" -name 'historykanban-*.sql.gz' -mtime "+$RETENTION" -print -delete | wc -l)
    if [ "$supprimees" -gt 0 ]; then
        journaliser "$supprimees archive(s) de plus de $RETENTION jours supprimee(s)"
    fi
}

journaliser "service demarre, une sauvegarde toutes les $INTERVALLE secondes, retention $RETENTION jours"

while true; do
    until pg_isready -h "$BDHOTE" -U "$BDUTILISATEUR" >/dev/null 2>&1; do
        journaliser "base indisponible, nouvelle tentative dans 10 secondes"
        sleep 10
    done
    sauvegarder
    purger
    sleep "$INTERVALLE"
done
