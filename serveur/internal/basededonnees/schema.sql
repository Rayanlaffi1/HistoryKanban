CREATE TABLE IF NOT EXISTS utilisateurs (
    id UUID PRIMARY KEY,
    courriel TEXT UNIQUE NOT NULL,
    nom TEXT NOT NULL DEFAULT '',
    prenom TEXT NOT NULL DEFAULT '',
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
    utilisateur UUID PRIMARY KEY REFERENCES utilisateurs(id) ON DELETE CASCADE,
    session TEXT NOT NULL,
    connexion BIGINT NOT NULL DEFAULT 0,
    modification TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS groupes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nom TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    proprietaire UUID NOT NULL REFERENCES utilisateurs(id),
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS membres (
    groupe UUID NOT NULL REFERENCES groupes(id) ON DELETE CASCADE,
    utilisateur UUID NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'membre' CHECK (role IN ('proprietaire', 'administrateur', 'membre', 'lecteur')),
    fonction TEXT NOT NULL DEFAULT '',
    ajout TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (groupe, utilisateur)
);

ALTER TABLE membres ADD COLUMN IF NOT EXISTS fonction TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS projets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    groupe UUID NOT NULL REFERENCES groupes(id) ON DELETE CASCADE,
    nom TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    couleur TEXT NOT NULL DEFAULT '#737373',
    archive BOOLEAN NOT NULL DEFAULT false,
    createur UUID NOT NULL REFERENCES utilisateurs(id),
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS colonnes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    projet UUID NOT NULL REFERENCES projets(id) ON DELETE CASCADE,
    nom TEXT NOT NULL,
    couleur TEXT NOT NULL DEFAULT '#a3a3a3',
    position INTEGER NOT NULL DEFAULT 0,
    limite INTEGER
);

CREATE TABLE IF NOT EXISTS lots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    projet UUID NOT NULL REFERENCES projets(id) ON DELETE CASCADE,
    nom TEXT NOT NULL,
    couleur TEXT NOT NULL DEFAULT '#525252',
    echeance TIMESTAMPTZ,
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS etiquettes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    projet UUID NOT NULL REFERENCES projets(id) ON DELETE CASCADE,
    nom TEXT NOT NULL,
    couleur TEXT NOT NULL DEFAULT '#737373'
);

CREATE TABLE IF NOT EXISTS taches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    projet UUID NOT NULL REFERENCES projets(id) ON DELETE CASCADE,
    colonne UUID NOT NULL REFERENCES colonnes(id) ON DELETE CASCADE,
    lot UUID REFERENCES lots(id) ON DELETE SET NULL,
    titre TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    points INTEGER NOT NULL DEFAULT 0,
    urgence TEXT NOT NULL DEFAULT 'normale' CHECK (urgence IN ('faible','normale','elevee','urgente')),
    echeance TIMESTAMPTZ,
    urls TEXT[] NOT NULL DEFAULT '{}',
    position INTEGER NOT NULL DEFAULT 0,
    suppression TIMESTAMPTZ,
    terminee TIMESTAMPTZ,
    createur UUID NOT NULL REFERENCES utilisateurs(id),
    creation TIMESTAMPTZ NOT NULL DEFAULT now(),
    modification TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS affectations (
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    utilisateur UUID NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    PRIMARY KEY (tache, utilisateur)
);

CREATE TABLE IF NOT EXISTS etiquetages (
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    etiquette UUID NOT NULL REFERENCES etiquettes(id) ON DELETE CASCADE,
    PRIMARY KEY (tache, etiquette)
);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    chemin TEXT NOT NULL,
    nom TEXT NOT NULL,
    taille BIGINT NOT NULL DEFAULT 0,
    typecontenu TEXT NOT NULL DEFAULT 'image/png',
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS soustaches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    libelle TEXT NOT NULL,
    faite BOOLEAN NOT NULL DEFAULT false,
    position INTEGER NOT NULL DEFAULT 0,
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS soustaches_tache ON soustaches(tache);

CREATE TABLE IF NOT EXISTS commentaires (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    auteur UUID NOT NULL REFERENCES utilisateurs(id),
    contenu TEXT NOT NULL,
    agent BOOLEAN NOT NULL DEFAULT false,
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS activites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    utilisateur UUID REFERENCES utilisateurs(id),
    type TEXT NOT NULL,
    detail TEXT NOT NULL DEFAULT '',
    agent BOOLEAN NOT NULL DEFAULT false,
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    utilisateur UUID NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    contenu JSONB NOT NULL DEFAULT '{}',
    lue BOOLEAN NOT NULL DEFAULT false,
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE taches ADD COLUMN IF NOT EXISTS urls TEXT[] NOT NULL DEFAULT '{}';
ALTER TABLE taches ADD COLUMN IF NOT EXISTS urgence TEXT NOT NULL DEFAULT 'normale';
ALTER TABLE images ADD COLUMN IF NOT EXISTS typecontenu TEXT NOT NULL DEFAULT 'image/png';
ALTER TABLE taches ADD COLUMN IF NOT EXISTS suppression TIMESTAMPTZ;
ALTER TABLE taches ADD COLUMN IF NOT EXISTS terminee TIMESTAMPTZ;
ALTER TABLE activites ADD COLUMN IF NOT EXISTS agent BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE commentaires ADD COLUMN IF NOT EXISTS agent BOOLEAN NOT NULL DEFAULT false;

UPDATE taches t SET terminee = t.modification
FROM colonnes c
WHERE c.id = t.colonne
  AND t.terminee IS NULL
  AND c.position = (SELECT max(cc.position) FROM colonnes cc WHERE cc.projet = t.projet);

CREATE INDEX IF NOT EXISTS taches_suppression ON taches(suppression);
CREATE INDEX IF NOT EXISTS taches_terminee ON taches(terminee);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'taches_urgence_valide') THEN
        ALTER TABLE taches ADD CONSTRAINT taches_urgence_valide
            CHECK (urgence IN ('faible','normale','elevee','urgente'));
    END IF;
END $$;

-- Migration : l'ancien numero de commit devient une URL complete rattachee a la tache,
-- construite depuis le depot du projet quand il etait renseigne, puis les colonnes
-- commit (taches) et depot (projets) sont supprimees.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
        WHERE table_name = 'taches' AND column_name = 'commit') THEN
        UPDATE taches t SET urls = ARRAY[rtrim(p.depot, '/') || '/commit/' || t.commit]
        FROM projets p
        WHERE p.id = t.projet AND t.commit <> '' AND p.depot <> '' AND t.urls = '{}';
        UPDATE taches SET urls = ARRAY[taches.commit]
        WHERE taches.commit <> '' AND urls = '{}';
        ALTER TABLE taches DROP COLUMN "commit";
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns
        WHERE table_name = 'projets' AND column_name = 'depot') THEN
        ALTER TABLE projets DROP COLUMN depot;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS cles (
    projet UUID NOT NULL REFERENCES projets(id) ON DELETE CASCADE,
    utilisateur UUID NOT NULL REFERENCES utilisateurs(id) ON DELETE CASCADE,
    cle TEXT NOT NULL UNIQUE,
    creation TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (projet, utilisateur)
);

UPDATE cles
SET cle = encode(sha256(cle::bytea), 'hex')
WHERE cle !~ '^[a-f0-9]{64}$';

CREATE TABLE IF NOT EXISTS preferences (
    utilisateur UUID PRIMARY KEY REFERENCES utilisateurs(id) ON DELETE CASCADE,
    courriels BOOLEAN NOT NULL DEFAULT true,
    types JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS indexmembresutilisateur ON membres(utilisateur);
CREATE INDEX IF NOT EXISTS indexprojetsgroupe ON projets(groupe);
CREATE INDEX IF NOT EXISTS indexcolonnesprojet ON colonnes(projet);
CREATE INDEX IF NOT EXISTS indextachesprojet ON taches(projet);
CREATE INDEX IF NOT EXISTS indextachescolonne ON taches(colonne);
CREATE INDEX IF NOT EXISTS indexnotificationsutilisateur ON notifications(utilisateur, lue);
CREATE INDEX IF NOT EXISTS indexcommentairestache ON commentaires(tache);
CREATE INDEX IF NOT EXISTS indexactivitestache ON activites(tache);
