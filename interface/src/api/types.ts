import { z } from "zod"

export const schemaUtilisateur = z.object({
  id: z.string(),
  courriel: z.string(),
  nom: z.string(),
  prenom: z.string(),
})

export const schemaProfil = schemaUtilisateur.extend({
  roles: z.array(z.string()).nullish(),
})

export const schemaGroupe = z.object({
  id: z.string(),
  nom: z.string(),
  description: z.string(),
  proprietaire: z.string(),
  creation: z.string(),
  role: z.string(),
  nbmembres: z.number(),
  nbprojets: z.number(),
})

export const schemaMembre = z.object({
  groupe: z.string(),
  utilisateur: z.string(),
  role: z.string(),
  fonction: z.string().default(""),
  ajout: z.string(),
  courriel: z.string(),
  nom: z.string(),
  prenom: z.string(),
})

export const schemaProjet = z.object({
  id: z.string(),
  groupe: z.string(),
  nom: z.string(),
  description: z.string(),
  couleur: z.string(),
  archive: z.boolean(),
  depot: z.string().default(""),
  createur: z.string(),
  creation: z.string(),
  nbtaches: z.number().optional().default(0),
})

export const schemaColonne = z.object({
  id: z.string(),
  projet: z.string(),
  nom: z.string(),
  couleur: z.string(),
  position: z.number(),
  limite: z.number().nullable(),
})

export const schemaLot = z.object({
  id: z.string(),
  projet: z.string(),
  nom: z.string(),
  couleur: z.string(),
  echeance: z.string().nullable(),
  creation: z.string(),
})

export const schemaEtiquette = z.object({
  id: z.string(),
  projet: z.string(),
  nom: z.string(),
  couleur: z.string(),
})

export const schemaImage = z.object({
  id: z.string(),
  tache: z.string(),
  chemin: z.string(),
  nom: z.string(),
  taille: z.number(),
  url: z.string(),
  typecontenu: z.string().default(""),
  creation: z.string(),
})

export const schemaSousTache = z.object({
  id: z.string(),
  tache: z.string(),
  libelle: z.string(),
  faite: z.boolean(),
  position: z.number(),
  creation: z.string(),
})

export const schemaTache = z.object({
  id: z.string(),
  projet: z.string(),
  colonne: z.string(),
  lot: z.string().nullable(),
  titre: z.string(),
  description: z.string(),
  points: z.number(),
  urgence: z.enum(["faible", "normale", "elevee", "urgente"]).default("normale"),
  echeance: z.string().nullable(),
  commit: z.string().default(""),
  position: z.number(),
  suppression: z.string().nullable().default(null),
  createur: z.string(),
  creation: z.string(),
  modification: z.string(),
  affectations: z.array(z.string()),
  etiquettes: z.array(z.string()),
  images: z.array(schemaImage),
  soustaches: z.array(schemaSousTache).default([]),
})

export const schemaResultat = z.object({
  tache: z.string(),
  titre: z.string(),
  urgence: z.string().default("normale"),
  echeance: z.string().nullable(),
  projet: z.string(),
  projetnom: z.string(),
  projetcouleur: z.string(),
  groupe: z.string(),
  groupenom: z.string(),
  colonne: z.string(),
  origine: z.string(),
  extrait: z.string().default(""),
  modification: z.string(),
})

export const schemaCommentaire = z.object({
  id: z.string(),
  tache: z.string(),
  auteur: z.string(),
  contenu: z.string(),
  creation: z.string(),
  nom: z.string(),
  prenom: z.string(),
})

export const schemaActivite = z.object({
  id: z.string(),
  tache: z.string(),
  utilisateur: z.string(),
  type: z.string(),
  detail: z.string(),
  creation: z.string(),
  nom: z.string(),
  prenom: z.string(),
})

export const schemaNotification = z.object({
  id: z.string(),
  utilisateur: z.string(),
  type: z.string(),
  contenu: z.object({
    titre: z.string().optional(),
    projet: z.string().optional(),
    acteurnom: z.string().optional(),
  }),
  lue: z.boolean(),
  creation: z.string(),
})

export const schemaPreferences = z.object({
  utilisateur: z.string(),
  courriels: z.boolean(),
  types: z.record(z.boolean()),
})

export const schemaLigneClassement = z.object({
  utilisateur: z.string(),
  nom: z.string(),
  prenom: z.string(),
  points: z.number(),
  taches: z.number(),
  creees: z.number(),
  commentaires: z.number(),
})

export const schemaStatistiques = z.object({
  classement: z.array(schemaLigneClassement),
  serie: z.array(
    z.object({
      periode: z.string(),
      points: z.number(),
      taches: z.number(),
    }),
  ),
  totaux: z.object({
    points: z.number(),
    terminees: z.number(),
    enretard: z.number(),
    total: z.number(),
    totalpoints: z.number(),
  }),
})

export const schemaCle = z.object({
  projet: z.string(),
  utilisateur: z.string(),
  cle: z.string(),
  creation: z.string(),
})

export const schemaReponseCle = z.object({
  cle: schemaCle.nullable(),
})

export const schemaDetailProjet = z.object({
  projet: schemaProjet,
  colonnes: z.array(schemaColonne),
  etiquettes: z.array(schemaEtiquette),
  lots: z.array(schemaLot),
  membres: z.array(schemaMembre),
  role: z.string(),
})

export type Utilisateur = z.infer<typeof schemaUtilisateur>
export type Profil = z.infer<typeof schemaProfil>
export type Groupe = z.infer<typeof schemaGroupe>
export type Membre = z.infer<typeof schemaMembre>
export type Projet = z.infer<typeof schemaProjet>
export type Colonne = z.infer<typeof schemaColonne>
export type Lot = z.infer<typeof schemaLot>
export type Etiquette = z.infer<typeof schemaEtiquette>
export type Image = z.infer<typeof schemaImage>
export type Tache = z.infer<typeof schemaTache>
export type Commentaire = z.infer<typeof schemaCommentaire>
export type Activite = z.infer<typeof schemaActivite>
export type Cle = z.infer<typeof schemaCle>
export type Statistiques = z.infer<typeof schemaStatistiques>
export type LigneClassement = z.infer<typeof schemaLigneClassement>

export interface ParametresStatistiques {
  debut: string
  fin: string
  granularite: "jour" | "semaine" | "mois"
  projet: string
}
export type SousTache = z.infer<typeof schemaSousTache>
export type Resultat = z.infer<typeof schemaResultat>
export type Notification = z.infer<typeof schemaNotification>
export type Preferences = z.infer<typeof schemaPreferences>
export type DetailProjet = z.infer<typeof schemaDetailProjet>

export interface FiltreTaches {
  texte?: string
  membre?: string
  etiquette?: string
  lot?: string
  echeance?: string
  urgence?: string
  pointsmin?: number
  pointsmax?: number
}

export const urgences = ["faible", "normale", "elevee", "urgente"] as const

export type Urgence = (typeof urgences)[number]

export const libellesUrgences: Record<Urgence, string> = {
  faible: "Faible",
  normale: "Normale",
  elevee: "Élevée",
  urgente: "Urgente",
}

export const couleursUrgences: Record<Urgence, string> = {
  faible: "bg-neutral-100 text-neutral-500 dark:bg-neutral-800 dark:text-neutral-400",
  normale: "bg-neutral-200 text-neutral-700 dark:bg-neutral-700 dark:text-neutral-200",
  elevee: "bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300",
  urgente: "bg-red-100 text-red-800 dark:bg-red-900/40 dark:text-red-300",
}

export const roles = ["administrateur", "membre", "lecteur"] as const

export const libellesRoles: Record<string, string> = {
  proprietaire: "Propriétaire",
  administrateur: "Administrateur",
  membre: "Membre",
  lecteur: "Lecteur",
}

export const typesNotifications: Record<string, string> = {
  "tache.creee": "Création de tâche",
  "tache.modifiee": "Modification de tâche",
  "tache.deplacee": "Déplacement de tâche",
  "tache.supprimee": "Suppression de tâche",
  "tache.commentee": "Commentaire",
  "tache.image.ajoutee": "Ajout d'image",
  "groupe.membre.ajoute": "Ajout à un groupe",
}

export const libellesActivites: Record<string, string> = {
  "tache.creee": "a créé la tâche",
  "tache.modifiee": "a modifié la tâche",
  "tache.deplacee": "a déplacé la tâche",
  "tache.commentee": "a commenté",
  "tache.image.ajoutee": "a ajouté une image",
  "tache.image.supprimee": "a supprimé une image",
}
