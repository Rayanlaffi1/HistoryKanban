import { computed, type Ref } from "vue"
import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { z } from "zod"
import { client } from "@/api/client"
import {
  schemaActivite,
  schemaCommentaire,
  schemaDetailProjet,
  schemaGroupe,
  schemaMembre,
  schemaNotification,
  schemaPreferences,
  schemaProfil,
  schemaProjet,
  schemaReponseCle,
  schemaStatistiques,
  schemaTache,
  type FiltreTaches,
  type ParametresStatistiques,
  type Preferences,
} from "@/api/types"

function invalidation(cles: (variables: any) => unknown[][]) {
  const clientRequetes = useQueryClient()
  return (_donnees: unknown, variables: unknown) => {
    for (const cle of cles(variables)) {
      clientRequetes.invalidateQueries({ queryKey: cle })
    }
  }
}

export function utiliserProfil() {
  return useQuery({
    queryKey: ["profil"],
    queryFn: async () => schemaProfil.parse((await client.get("/moi")).data),
  })
}

export function utiliserGroupes() {
  return useQuery({
    queryKey: ["groupes"],
    queryFn: async () => z.array(schemaGroupe).parse((await client.get("/groupes")).data),
  })
}

export function utiliserMembres(groupe: Ref<string>) {
  return useQuery({
    queryKey: computed(() => ["membres", groupe.value]),
    queryFn: async () => z.array(schemaMembre).parse((await client.get(`/groupes/${groupe.value}/membres`)).data),
    enabled: computed(() => groupe.value !== ""),
  })
}

export function utiliserProjets(groupe: Ref<string>) {
  return useQuery({
    queryKey: computed(() => ["projets", groupe.value]),
    queryFn: async () => z.array(schemaProjet).parse((await client.get(`/groupes/${groupe.value}/projets`)).data),
    enabled: computed(() => groupe.value !== ""),
  })
}

export function utiliserDetailProjet(projet: Ref<string>) {
  return useQuery({
    queryKey: computed(() => ["projet", projet.value]),
    queryFn: async () => schemaDetailProjet.parse((await client.get(`/projets/${projet.value}`)).data),
    enabled: computed(() => projet.value !== ""),
  })
}

export function utiliserTaches(projet: Ref<string>, filtre: Ref<FiltreTaches>) {
  return useQuery({
    queryKey: computed(() => ["taches", projet.value, filtre.value]),
    queryFn: async () =>
      z.array(schemaTache).parse((await client.get(`/projets/${projet.value}/taches`, { params: filtre.value })).data),
    enabled: computed(() => projet.value !== ""),
  })
}

export function utiliserCommentaires(tache: Ref<string>) {
  return useQuery({
    queryKey: computed(() => ["commentaires", tache.value]),
    queryFn: async () => z.array(schemaCommentaire).parse((await client.get(`/taches/${tache.value}/commentaires`)).data),
    enabled: computed(() => tache.value !== ""),
  })
}

export function utiliserStatistiques(groupe: Ref<string>, parametres: Ref<ParametresStatistiques>) {
  return useQuery({
    queryKey: computed(() => ["statistiques", groupe.value, parametres.value]),
    queryFn: async () =>
      schemaStatistiques.parse(
        (
          await client.get(`/groupes/${groupe.value}/statistiques`, {
            params: {
              debut: parametres.value.debut,
              fin: parametres.value.fin,
              granularite: parametres.value.granularite,
              projet: parametres.value.projet || undefined,
            },
          })
        ).data,
      ),
    enabled: computed(() => groupe.value !== ""),
  })
}

export function utiliserCle(projet: Ref<string>, active: Ref<boolean>) {
  return useQuery({
    queryKey: computed(() => ["cle", projet.value]),
    queryFn: async () => schemaReponseCle.parse((await client.get(`/projets/${projet.value}/cle`)).data).cle,
    enabled: computed(() => projet.value !== "" && active.value),
  })
}

export function utiliserGenerationCle() {
  return useMutation({
    mutationFn: (projet: string) => client.post(`/projets/${projet}/cle`),
    onSuccess: invalidation((projet: string) => [["cle", projet]]),
  })
}

export function utiliserSuppressionCle() {
  return useMutation({
    mutationFn: (projet: string) => client.delete(`/projets/${projet}/cle`),
    onSuccess: invalidation((projet: string) => [["cle", projet]]),
  })
}

export function utiliserActivites(tache: Ref<string>) {
  return useQuery({
    queryKey: computed(() => ["activites", tache.value]),
    queryFn: async () => z.array(schemaActivite).parse((await client.get(`/taches/${tache.value}/activites`)).data),
    enabled: computed(() => tache.value !== ""),
  })
}

export function utiliserNotifications() {
  return useQuery({
    queryKey: ["notifications"],
    queryFn: async () => z.array(schemaNotification).parse((await client.get("/notifications")).data),
  })
}

export function utiliserPreferences() {
  return useQuery({
    queryKey: ["preferences"],
    queryFn: async () => schemaPreferences.parse((await client.get("/moi/preferences")).data),
  })
}

export function utiliserMutationGroupe() {
  return useMutation({
    mutationFn: (corps: { id?: string; nom: string; description: string }) =>
      corps.id ? client.put(`/groupes/${corps.id}`, corps) : client.post("/groupes", corps),
    onSuccess: invalidation(() => [["groupes"]]),
  })
}

export function utiliserSuppressionGroupe() {
  return useMutation({
    mutationFn: (id: string) => client.delete(`/groupes/${id}`),
    onSuccess: invalidation(() => [["groupes"]]),
  })
}

export function utiliserMutationMembre() {
  return useMutation({
    mutationFn: (corps: { groupe: string; courriel: string; role: string }) =>
      client.post(`/groupes/${corps.groupe}/membres`, corps),
    onSuccess: invalidation((variables: { groupe: string }) => [["membres", variables.groupe], ["groupes"]]),
  })
}

export function utiliserMutationRole() {
  return useMutation({
    mutationFn: (corps: { groupe: string; utilisateur: string; role: string }) =>
      client.put(`/groupes/${corps.groupe}/membres/${corps.utilisateur}`, { role: corps.role }),
    onSuccess: invalidation((variables: { groupe: string }) => [["membres", variables.groupe]]),
  })
}

export function utiliserMutationFonction() {
  return useMutation({
    mutationFn: (corps: { groupe: string; utilisateur: string; fonction: string }) =>
      client.put(`/groupes/${corps.groupe}/membres/${corps.utilisateur}/fonction`, { fonction: corps.fonction }),
    onSuccess: invalidation((variables: { groupe: string }) => [["membres", variables.groupe]]),
  })
}

export function utiliserRetraitMembre() {
  return useMutation({
    mutationFn: (corps: { groupe: string; utilisateur: string }) =>
      client.delete(`/groupes/${corps.groupe}/membres/${corps.utilisateur}`),
    onSuccess: invalidation((variables: { groupe: string }) => [["membres", variables.groupe], ["groupes"]]),
  })
}

export function utiliserMutationProjet() {
  return useMutation({
    mutationFn: (corps: { id?: string; groupe: string; nom: string; description: string; couleur: string; depot?: string; archive?: boolean }) =>
      corps.id ? client.put(`/projets/${corps.id}`, corps) : client.post(`/groupes/${corps.groupe}/projets`, corps),
    onSuccess: invalidation((variables: { id?: string; groupe: string }) => [
      ["projets", variables.groupe],
      ["groupes"],
      ...(variables.id ? [["projet", variables.id]] : []),
    ]),
  })
}

export function utiliserSuppressionProjet() {
  return useMutation({
    mutationFn: (corps: { id: string; groupe: string }) => client.delete(`/projets/${corps.id}`),
    onSuccess: invalidation((variables: { groupe: string }) => [["projets", variables.groupe], ["groupes"]]),
  })
}

export function utiliserMutationColonne() {
  return useMutation({
    mutationFn: (corps: { id?: string; projet: string; nom: string; couleur: string; limite: number | null }) =>
      corps.id ? client.put(`/colonnes/${corps.id}`, corps) : client.post(`/projets/${corps.projet}/colonnes`, corps),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet]]),
  })
}

export function utiliserSuppressionColonne() {
  return useMutation({
    mutationFn: (corps: { id: string; projet: string }) => client.delete(`/colonnes/${corps.id}`),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet], ["taches", variables.projet]]),
  })
}

export function utiliserOrdreColonnes() {
  return useMutation({
    mutationFn: (corps: { projet: string; ordre: string[] }) =>
      client.put(`/projets/${corps.projet}/colonnes/ordre`, corps),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet]]),
  })
}

export function utiliserMutationEtiquette() {
  return useMutation({
    mutationFn: (corps: { id?: string; projet: string; nom: string; couleur: string }) =>
      corps.id ? client.put(`/etiquettes/${corps.id}`, corps) : client.post(`/projets/${corps.projet}/etiquettes`, corps),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet]]),
  })
}

export function utiliserSuppressionEtiquette() {
  return useMutation({
    mutationFn: (corps: { id: string; projet: string }) => client.delete(`/etiquettes/${corps.id}`),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet]]),
  })
}

export function utiliserMutationLot() {
  return useMutation({
    mutationFn: (corps: { id?: string; projet: string; nom: string; couleur: string; echeance: string | null }) =>
      corps.id ? client.put(`/lots/${corps.id}`, corps) : client.post(`/projets/${corps.projet}/lots`, corps),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet]]),
  })
}

export function utiliserSuppressionLot() {
  return useMutation({
    mutationFn: (corps: { id: string; projet: string }) => client.delete(`/lots/${corps.id}`),
    onSuccess: invalidation((variables: { projet: string }) => [["projet", variables.projet], ["taches", variables.projet]]),
  })
}

export interface CorpsTache {
  id?: string
  projet: string
  colonne?: string
  titre: string
  description: string
  lot: string | null
  points: number
  urgence: string
  echeance: string | null
  commit?: string
  affectations: string[]
  etiquettes: string[]
}

export function utiliserMutationTache() {
  return useMutation({
    mutationFn: (corps: CorpsTache) =>
      corps.id ? client.put(`/taches/${corps.id}`, corps) : client.post(`/projets/${corps.projet}/taches`, corps),
    onSuccess: invalidation((variables: { projet: string }) => [["taches", variables.projet]]),
  })
}

export function utiliserSuppressionTache() {
  return useMutation({
    mutationFn: (corps: { id: string; projet: string }) => client.delete(`/taches/${corps.id}`),
    onSuccess: invalidation((variables: { projet: string }) => [["taches", variables.projet]]),
  })
}

export function utiliserDeplacementTache() {
  return useMutation({
    mutationFn: (corps: { id: string; projet: string; colonne: string; position: number }) =>
      client.put(`/taches/${corps.id}/deplacer`, corps),
    onSuccess: invalidation((variables: { projet: string }) => [["taches", variables.projet]]),
  })
}

export function utiliserTeleversementImage() {
  return useMutation({
    mutationFn: (corps: { tache: string; projet: string; fichier: File }) => {
      const donnees = new FormData()
      donnees.append("fichier", corps.fichier)
      return client.post(`/taches/${corps.tache}/images`, donnees)
    },
    onSuccess: invalidation((variables: { projet: string }) => [["taches", variables.projet]]),
  })
}

export function utiliserSuppressionImage() {
  return useMutation({
    mutationFn: (corps: { id: string; projet: string }) => client.delete(`/images/${corps.id}`),
    onSuccess: invalidation((variables: { projet: string }) => [["taches", variables.projet]]),
  })
}

export function utiliserMutationCommentaire() {
  return useMutation({
    mutationFn: (corps: { tache: string; contenu: string }) =>
      client.post(`/taches/${corps.tache}/commentaires`, corps),
    onSuccess: invalidation((variables: { tache: string }) => [["commentaires", variables.tache]]),
  })
}

export function utiliserSuppressionCommentaire() {
  return useMutation({
    mutationFn: (corps: { id: string; tache: string }) => client.delete(`/commentaires/${corps.id}`),
    onSuccess: invalidation((variables: { tache: string }) => [["commentaires", variables.tache]]),
  })
}

export function utiliserLectureNotification() {
  return useMutation({
    mutationFn: (id: string) => client.put(`/notifications/${id}/lue`),
    onSuccess: invalidation(() => [["notifications"]]),
  })
}

export function utiliserLectureTotale() {
  return useMutation({
    mutationFn: () => client.put("/notifications/tout"),
    onSuccess: invalidation(() => [["notifications"]]),
  })
}

export async function televerserFichier(projet: string, fichier: File): Promise<string> {
  const donnees = new FormData()
  donnees.append("fichier", fichier)
  const reponse = await client.post(`/projets/${projet}/fichiers`, donnees)
  return reponse.data.url
}

export function utiliserMutationPreferences() {
  return useMutation({
    mutationFn: (corps: Omit<Preferences, "utilisateur">) => client.put("/moi/preferences", corps),
    onSuccess: invalidation(() => [["preferences"]]),
  })
}
