# HistoryKanban

Application de gestion de projets sous forme de tableaux kanban personnalisables, avec groupes, rôles, temps réel et notifications.

## Stack

| Couche | Technologies |
| --- | --- |
| Interface | Vue 3, TypeScript, Vite, Vue Router, Pinia, TanStack Query, TailwindCSS, style shadcn, VueUse, vue-draggable-plus, Axios, Zod, Day.js |
| Serveur | Go, Gin, JWT (Keycloak), PostgreSQL, RabbitMQ, MinIO, WebSocket |
| Infrastructure | Docker Compose, Traefik (HTTPS), Keycloak (thème custom, magic link, events), MinIO, RabbitMQ, MailHog, Jenkins |
| Sauvegarde | Service `pg_dumpall` periodique compresse, avec retention |

## Demarrage

Generer d'abord les identifiants et le certificat TLS locaux. Les fichiers produits sont ignores par Git.

```powershell
# Windows PowerShell
.\scripts\initialiser-secrets.ps1
```

```bash
# Linux et macOS
./scripts/initialiser-secrets.sh
```

Pour renouveler une configuration existante, utiliser `-Force` sous PowerShell ou `--force` sous Linux et macOS. Si des volumes Docker existent deja, leurs comptes doivent aussi etre renouveles dans PostgreSQL, Keycloak et RabbitMQ avant le redemarrage.

```bash
docker compose up -d --build
```

L'acces se fait en HTTPS via Traefik (certificat auto-signe, une exception de securite est a accepter dans le navigateur).

| Service | URL |
| --- | --- |
| Interface | https://historykanban.localhost |
| API | https://historykanban.localhost/api (aussi : https://api.historykanban.localhost, direct : http://localhost:8085) |
| Keycloak | https://auth.historykanban.localhost |
| Images MinIO | https://images.historykanban.localhost |
| MailHog | https://courriel.historykanban.localhost |
| Console MinIO | http://localhost:9001 |
| RabbitMQ | http://localhost:15673 |
| Jenkins | http://localhost:8082 |
| Tableau de bord Traefik | http://localhost:8083 |

## Acces depuis le reseau local (sslip.io)

L'application est aussi accessible depuis une autre machine du reseau grace aux domaines [sslip.io](https://sslip.io), qui resolvent `historykanban.<IP>.sslip.io` vers `<IP>` sans configuration DNS. Les deux machines doivent pouvoir resoudre les DNS publics (sslip.io est un service DNS public).

1. Trouver l'adresse IP locale de la machine qui heberge la pile : `ipconfig` sous Windows (champ « Adresse IPv4 »), `ip addr` ou `ifconfig` sous Linux et macOS. Exemple : `192.168.1.42`.
2. Renseigner cette adresse dans `.env` avec la ligne `IPRESEAU=192.168.1.42` (ou generer les secrets avec `.\scripts\initialiser-secrets.ps1 -IpReseau 192.168.1.42` / `./scripts/initialiser-secrets.sh --ip 192.168.1.42`, ce qui ajoute aussi le domaine au certificat).
3. Redemarrer la pile : `docker compose up -d --build`.
4. Depuis n'importe quelle machine du reseau, ouvrir :

| Service | URL |
| --- | --- |
| Interface | https://historykanban.192.168.1.42.sslip.io |
| API | https://historykanban.192.168.1.42.sslip.io/api |
| Keycloak | https://auth.historykanban.192.168.1.42.sslip.io |
| Images MinIO | https://images.historykanban.192.168.1.42.sslip.io |
| MailHog | https://courriel.historykanban.192.168.1.42.sslip.io |

Le certificat reste auto-signe : accepter l'exception de securite pour l'interface puis pour Keycloak (ouvrir une fois `https://auth.historykanban.<IP>.sslip.io` dans un onglet). Le routage Traefik accepte n'importe quelle IP dans le domaine, mais la connexion Keycloak et la validation des jetons ne fonctionnent que pour l'adresse declaree dans `IPRESEAU`.

Le realm Keycloak n'est importe qu'au premier demarrage : si la pile a deja tourne avant le changement d'`IPRESEAU`, ajouter `https://historykanban.<IP>.sslip.io/*` aux « Valid redirect URIs » du client `interface` dans la console d'administration Keycloak (ou reinitialiser les volumes avec `docker compose down -v`, ce qui supprime toutes les donnees).

Le fonctionnement via `https://historykanban.localhost` reste inchange : les deux acces cohabitent, l'interface derive ses URL publiques de l'adresse consultee.

## Premier compte

Aucun compte applicatif n'est preconfigure dans le realm Keycloak. Creer le premier compte depuis la page d'inscription. Le nom de l'administrateur Keycloak est defini dans `.env` et son mot de passe aleatoire est genere localement par le script d'initialisation.

## Fonctionnalites

- Projets avec tableaux kanban personnalisables (colonnes, couleurs, limites)
- Taches avec pieces jointes (MinIO) : images, PDF, documents bureautiques, archives et fichiers texte
- Taches avec points, niveau d'urgence, echeances, etiquettes, affectations
- Sous-taches cochables avec indicateur d'avancement sur la carte
- Lots de taches avec echeance commune
- Groupes d'utilisateurs avec roles et droits (proprietaire, administrateur, membre, lecteur)
- Filtres des taches (texte, membre, etiquette, echeance, urgence, points)
- Recherche globale depuis l'en-tete : titres, descriptions et commentaires de tous les groupes accessibles
- Corbeille : une tache supprimee reste restaurable trente jours avant effacement definitif
- Sauvegarde automatique et compressee de la base, avec retention et procedure de restauration
- Temps reel via WebSocket : creation et modification instantanees, notifications
- Courriels configurables a chaque changement de tache (RabbitMQ + SMTP)
- Session unique : une nouvelle connexion deconnecte instantanement la precedente
- Mode sombre, palette gris et blanc
- Connexion par mot de passe ou lien magique (Keycloak)
- Acces IA : cle d'API par projet et par utilisateur pour piloter le tableau depuis un agent

## API agent (acces IA)

Chaque membre d'un projet peut generer une cle d'API personnelle depuis le bouton « Acces IA » de la page projet. Cette cle permet a un agent (IA, script, integration continue) d'agir sur ce projet en son nom, avec ses droits.

- Adresse de base : `https://historykanban.localhost/api/agent`
- Authentification : en-tete `X-Cle-API: <cle>` ou `Authorization: Bearer <cle>`

| Methode | Chemin | Description |
| --- | --- | --- |
| GET | `/projet` | Colonnes, membres, etiquettes et lots du projet |
| GET | `/membres` | Membres du groupe : identifiant, nom, courriel, role et fonction |
| GET | `/taches` | Toutes les taches du tableau (filtrables par `?urgence=`) |
| GET | `/taches/{id}/soustaches` | Sous-taches d'une tache |
| POST | `/taches/{id}/soustaches` | Ajouter une sous-tache `{ libelle }` |
| PUT | `/soustaches/{id}` | Cocher ou decocher une sous-tache `{ faite }`, renommer avec `{ libelle }` |
| GET | `/corbeille` | Taches supprimees encore restaurables |
| PUT | `/taches/{id}/restaurer` | Sortir une tache de la corbeille |
| POST | `/taches` | Creer une tache `{ titre, colonne, description?, points?, urgence?, echeance?, lot?, affectations?, etiquettes? }` |
| GET | `/taches/{id}` | Detail d'une tache |
| PUT | `/taches/{id}` | Modifier une tache `{ titre, description?, points?, urgence?, echeance?, lot?, urls? }` |
| PUT | `/taches/{id}/deplacer` | Deplacer une tache `{ colonne, position }` |
| PUT | `/taches/{id}/affectations` | Remplacer les personnes affectees `{ affectations: [identifiants] }` |
| PUT | `/taches/{id}/etiquettes` | Remplacer les etiquettes `{ etiquettes: [identifiants] }` |
| PUT | `/taches/{id}/urls` | Renseigner les liens du code d'une tache terminee `{ urls: [liens complets] }` |
| GET | `/taches/{id}/commentaires` | Lire les commentaires |
| POST | `/taches/{id}/commentaires` | Commenter `{ contenu }` |
| GET | `/taches/{id}/activites` | Journal d'activite de la tache |
| POST | `/etiquettes` | Creer une etiquette `{ nom, couleur? }` |
| PUT | `/etiquettes/{id}` | Renommer une etiquette `{ nom, couleur? }` |
| DELETE | `/etiquettes/{id}` | Supprimer une etiquette |
| POST | `/lots` | Creer un lot `{ nom, couleur?, echeance? }` |
| PUT | `/lots/{id}` | Modifier un lot `{ nom, couleur?, echeance? }` |
| DELETE | `/lots/{id}` | Supprimer un lot |

Une affectation n'est acceptee que si la personne appartient au groupe du projet, et une etiquette que si elle appartient au projet : sinon l'API repond 400 avec le motif.

Le champ `urgence` vaut `faible`, `normale` (defaut), `elevee` ou `urgente`. Toute autre valeur est refusee en 400. Sur `PUT /taches/{id}`, l'omettre conserve le niveau existant : l'agent peut ainsi trier son travail en lisant `urgence` sans jamais l'ecraser par megarde.

### Accents et encodage

Les corps sont attendus en UTF-8. Sous Windows, ecrire des accents directement dans `-d "..."` les corrompt (la console utilise Windows-1252) : passer par un fichier UTF-8 avec `--data-binary @fichier.json`.

```bash
curl -k -X POST https://historykanban.localhost/api/agent/taches \
  -H "X-Cle-API: <cle>" \
  -H "Content-Type: application/json; charset=utf-8" \
  --data-binary @tache.json
```

Par tolerance, l'API agent retire le BOM ajoute par `Set-Content -Encoding utf8` et reinterprete en UTF-8 les corps envoyes en Windows-1252, ce qui evite les caracteres corrompus en base.

Les actions de l'agent apparaissent en temps reel sur le tableau et alimentent notifications, courriels et journal d'activite comme toute action humaine.

## Sauvegarde et restauration

Le service `sauvegarde` archive automatiquement l'instance Postgres complete (base `historykanban` **et** base `keycloak`, roles compris) avec `pg_dumpall` compresse en gzip. Les archives sont ecrites dans `./sauvegardes` sur la machine hote, sous la forme `historykanban-AAAAMMJJ-HHMMSS.sql.gz`.

| Variable `.env` | Defaut | Role |
| --- | --- | --- |
| `SAUVEGARDEINTERVALLE` | `86400` | Delai entre deux sauvegardes, en secondes |
| `SAUVEGARDERETENTION` | `14` | Age maximal d'une archive, en jours, avant suppression |

Une sauvegarde part au demarrage du service, puis a chaque intervalle. L'archive n'est renommee sous son nom definitif qu'une fois `pg_dumpall` termine avec succes : une sauvegarde interrompue ne laisse donc jamais de fichier tronque qu'on croirait valide. Les archives sont exclues du depot par `.gitignore`.

### Restaurer

Le dump contient les instructions de creation des bases : il se rejoue sur une instance vierge.

```bash
# Repartir d'un volume neuf
docker compose down
docker volume rm historykanban_basededonnees
docker compose up -d basededonnees

# Rejouer l'archive choisie
gunzip -c sauvegardes/historykanban-AAAAMMJJ-HHMMSS.sql.gz \
  | docker compose exec -T basededonnees psql -U historykanban -d postgres

docker compose up -d
```

Pour verifier une archive sans toucher a l'installation, la rejouer dans un conteneur jetable :

```bash
docker run -d --name verification -e POSTGRES_PASSWORD=test postgres:16-alpine
docker cp sauvegardes/historykanban-AAAAMMJJ-HHMMSS.sql.gz verification:/archive.sql.gz
docker exec verification sh -c "gunzip -c /archive.sql.gz | psql -U postgres -d postgres"
docker exec verification psql -U postgres -d historykanban -c "SELECT count(*) FROM taches"
docker rm -f verification
```

Les images des taches vivent dans MinIO, pas dans Postgres : sauvegarder aussi le volume `minio` si les pieces jointes doivent etre restaurees.

## Mise a jour automatique (Jenkins)

La page Parametres affiche une section « Mise a jour de l'application » reservee aux utilisateurs portant le role realm Keycloak `administrateur` (les roles de groupe lecteur/membre n'y donnent pas acces). Le bouton declenche le pipeline Jenkins qui reconstruit les images et redeploie la stack via `docker compose up -d --build`.

Fonctionnement :

1. Le serveur interroge l'API GitHub (`GITHUBDEPOT`, forme `proprietaire/depot`) pour recuperer le tag de la derniere release et le compare a `VERSIONAPPLICATION`. Si le tag differe, l'interface signale « Une nouvelle release est disponible » : c'est la condition de mise a jour.
2. `GET /api/systeme/maj` renvoie versions et disponibilite ; `POST /api/systeme/maj` declenche le job Jenkins. Les deux exigent le role `administrateur` et le serveur n'execute aucune commande : il appelle uniquement l'API Jenkins.
3. Le declenchement utilise `POST {JENKINSURL}/job/{JENKINSJOB}/build` authentifie par jeton API Jenkins (utilisateur + jeton, exempt de crumb CSRF). Le job execute le `Jenkinsfile` du depot : verification, construction des images serveur et interface, puis deploiement.

Variables a renseigner dans `.env` (voir `.env.example`, ne jamais commiter le jeton) :

| Variable | Role |
| --- | --- |
| `JENKINSURL` | URL de Jenkins vue par le serveur (defaut `http://jenkins:8080`) |
| `JENKINSJOB` | Nom du job (dossiers acceptes, ex. `historykanban/main`) |
| `JENKINSUTILISATEUR` | Utilisateur Jenkins proprietaire du jeton API |
| `JENKINSJETON` | Jeton API Jenkins (profil utilisateur > Security > API Token) |
| `GITHUBDEPOT` | Depot GitHub `proprietaire/depot` pour detecter les releases |
| `VERSIONAPPLICATION` | Version deployee, a aligner sur le tag de la release installee |

Tant que `JENKINSUTILISATEUR` et `JENKINSJETON` sont vides, le bouton reste desactive et l'endpoint repond 503 : rien n'est declenchable par defaut.

## Depannage

L'API et le WebSocket sont servis sous le meme domaine que l'interface (https://historykanban.localhost/api et /ws) : accepter le certificat de l'interface et celui de Keycloak lors de la connexion suffit. En cas d'erreur reseau, l'interface affiche un message explicite avec l'URL a verifier.

## Developpement de l'interface

```bash
cd interface
npm install
npm run dev
```

L'interface de developpement tourne sur http://localhost:5173 et consomme l'API sur http://localhost:8080.
