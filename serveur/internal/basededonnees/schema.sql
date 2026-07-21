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
    echeance TIMESTAMPTZ,
    position INTEGER NOT NULL DEFAULT 0,
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
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS commentaires (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    auteur UUID NOT NULL REFERENCES utilisateurs(id),
    contenu TEXT NOT NULL,
    creation TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS activites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tache UUID NOT NULL REFERENCES taches(id) ON DELETE CASCADE,
    utilisateur UUID REFERENCES utilisateurs(id),
    type TEXT NOT NULL,
    detail TEXT NOT NULL DEFAULT '',
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
