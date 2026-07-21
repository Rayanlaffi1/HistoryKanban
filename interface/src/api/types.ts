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
  echeance: z.string().nullable(),
  position: z.number(),
  createur: z.string(),
  creation: z.string(),
  modification: z.string(),
  affectations: z.array(z.string()),
  etiquettes: z.array(z.string()),
  images: z.array(schemaImage),
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
export type Notification = z.infer<typeof schemaNotification>
export type Preferences = z.infer<typeof schemaPreferences>
export type DetailProjet = z.infer<typeof schemaDetailProjet>

export interface FiltreTaches {
  texte?: string
  membre?: string
  etiquette?: string
  lot?: string
  echeance?: string
  pointsmin?: number
  pointsmax?: number
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
