import { z } from "zod"

const nomCourt = (message: string, maximum = 100) =>
  z.string().trim().min(1, message).max(maximum, `${maximum} caractères maximum`)

const couleurHexadecimale = z.string().regex(/^#[0-9a-fA-F]{6}$/, "Couleur invalide")

const entierFacultatif = (message: string) =>
  z.union([z.string(), z.number()]).refine((valeur) => {
    if (valeur === "") return true
    const nombre = Number(valeur)
    return Number.isInteger(nombre) && nombre >= 0
  }, message)

export const formulaireGroupe = z.object({
  nom: nomCourt("Le nom du groupe est requis"),
  description: z.string().trim().max(500, "500 caractères maximum"),
})

export const formulaireProjet = z.object({
  nom: nomCourt("Le nom du projet est requis"),
  description: z.string().trim().max(500, "500 caractères maximum"),
  couleur: couleurHexadecimale,
})

export const formulaireMembre = z.object({
  courriel: z.string().trim().email("Courriel invalide"),
  role: z.enum(["administrateur", "membre", "lecteur"], { message: "Rôle invalide" }),
})

export const formulaireColonne = z.object({
  nom: nomCourt("Le nom de la colonne est requis", 50),
  couleur: couleurHexadecimale,
  limite: entierFacultatif("La limite doit être un entier positif"),
})

export const formulaireTache = z.object({
  titre: nomCourt("Le titre est requis", 200),
  description: z.string().trim().max(5000, "5000 caractères maximum"),
  colonne: z.string().min(1, "La colonne est requise"),
  points: entierFacultatif("Les points doivent être un entier positif"),
})

export const formulaireEtiquette = z.object({
  nom: nomCourt("Le nom de l'étiquette est requis", 50),
})

export const formulaireLot = z.object({
  nom: nomCourt("Le nom du lot est requis"),
})

export function valider<Schema extends z.ZodTypeAny>(
  schema: Schema,
  donnees: unknown,
): { donnees: z.infer<Schema> | null; erreurs: Record<string, string> } {
  const resultat = schema.safeParse(donnees)
  if (resultat.success) {
    return { donnees: resultat.data, erreurs: {} }
  }
  const erreurs: Record<string, string> = {}
  for (const probleme of resultat.error.issues) {
    const cle = String(probleme.path[0] ?? "formulaire")
    if (!erreurs[cle]) {
      erreurs[cle] = probleme.message
    }
  }
  return { donnees: null, erreurs }
}
