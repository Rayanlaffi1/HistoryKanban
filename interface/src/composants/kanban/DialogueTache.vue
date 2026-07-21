<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue"
import type { Colonne, Etiquette, Lot, Membre, Tache } from "@/api/types"
import {
  utiliserCommentaires,
  utiliserMutationCommentaire,
  utiliserMutationTache,
  utiliserSuppressionCommentaire,
  utiliserSuppressionImage,
  utiliserSuppressionTache,
  utiliserTeleversementImage,
  utiliserProfil,
} from "@/api/requetes"
import { depuis, depuisChampDate, formaterDateHeure, versChampDate } from "@/utilitaires/dates"
import { formulaireTache, valider } from "@/utilitaires/validation"
import { extraireErreur } from "@/utilitaires/erreurs"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import Zone from "@/composants/ui/Zone.vue"
import Selection from "@/composants/ui/Selection.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"

const proprietes = defineProps<{
  ouvert: boolean
  projet: string
  colonnes: Colonne[]
  membres: Membre[]
  etiquettes: Etiquette[]
  lots: Lot[]
  tache: Tache | null
  colonneInitiale: string
  edition: boolean
}>()

const emissions = defineEmits<{ fermer: [] }>()

const formulaire = reactive({
  titre: "",
  description: "",
  colonne: "",
  lot: "",
  points: "0",
  echeance: "",
  affectations: [] as string[],
  etiquettes: [] as string[],
})

watch(
  () => [proprietes.ouvert, proprietes.tache] as const,
  ([ouvert]) => {
    if (!ouvert) return
    formulaire.titre = proprietes.tache?.titre ?? ""
    formulaire.description = proprietes.tache?.description ?? ""
    formulaire.colonne = proprietes.tache?.colonne ?? proprietes.colonneInitiale
    formulaire.lot = proprietes.tache?.lot ?? ""
    formulaire.points = String(proprietes.tache?.points ?? 0)
    formulaire.echeance = versChampDate(proprietes.tache?.echeance)
    formulaire.affectations = [...(proprietes.tache?.affectations ?? [])]
    formulaire.etiquettes = [...(proprietes.tache?.etiquettes ?? [])]
    viderEnAttente()
  },
  { immediate: true },
)

const identifiantTache = computed(() => proprietes.tache?.id ?? "")
const { data: profil } = utiliserProfil()
const { data: commentaires } = utiliserCommentaires(identifiantTache)

const mutationTache = utiliserMutationTache()
const suppressionTache = utiliserSuppressionTache()
const televersement = utiliserTeleversementImage()
const suppressionImage = utiliserSuppressionImage()
const mutationCommentaire = utiliserMutationCommentaire()
const suppressionCommentaire = utiliserSuppressionCommentaire()

const nouveauCommentaire = ref("")
const confirmationSuppression = ref(false)
const champFichier = ref<HTMLInputElement | null>(null)
const erreurs = ref<Record<string, string>>({})
const erreurApi = ref("")

interface FichierEnAttente {
  fichier: File
  url: string
}

const enAttente = ref<FichierEnAttente[]>([])

function viderEnAttente() {
  for (const element of enAttente.value) {
    URL.revokeObjectURL(element.url)
  }
  enAttente.value = []
}

function selectionner(evenement: Event) {
  const fichiers = (evenement.target as HTMLInputElement).files
  if (!fichiers) return
  for (const fichier of Array.from(fichiers)) {
    enAttente.value.push({ fichier, url: URL.createObjectURL(fichier) })
  }
  ;(evenement.target as HTMLInputElement).value = ""
}

function retirerEnAttente(indice: number) {
  URL.revokeObjectURL(enAttente.value[indice].url)
  enAttente.value.splice(indice, 1)
}

function basculer(liste: string[], valeur: string) {
  const indice = liste.indexOf(valeur)
  if (indice >= 0) liste.splice(indice, 1)
  else liste.push(valeur)
}

function enregistrer() {
  erreurApi.value = ""
  const resultat = valider(formulaireTache, {
    titre: formulaire.titre,
    description: formulaire.description,
    colonne: formulaire.colonne,
    points: formulaire.points,
  })
  erreurs.value = resultat.erreurs
  if (!resultat.donnees) return
  mutationTache.mutate(
    {
      id: proprietes.tache?.id,
      projet: proprietes.projet,
      colonne: resultat.donnees.colonne,
      titre: resultat.donnees.titre,
      description: resultat.donnees.description,
      lot: formulaire.lot || null,
      points: Number(resultat.donnees.points) || 0,
      echeance: depuisChampDate(formulaire.echeance),
      affectations: formulaire.affectations,
      etiquettes: formulaire.etiquettes,
    },
    {
      onSuccess: async (reponse) => {
        if (!proprietes.tache && enAttente.value.length > 0) {
          const identifiantCree = reponse.data?.id
          if (identifiantCree) {
            for (const element of enAttente.value) {
              try {
                await televersement.mutateAsync({
                  tache: identifiantCree,
                  projet: proprietes.projet,
                  fichier: element.fichier,
                })
              } catch {
                break
              }
            }
          }
          viderEnAttente()
        }
        emissions("fermer")
      },
      onError: (erreur) => {
        erreurApi.value = extraireErreur(erreur)
      },
    },
  )
}

function supprimer() {
  if (!proprietes.tache) return
  suppressionTache.mutate(
    { id: proprietes.tache.id, projet: proprietes.projet },
    {
      onSuccess: () => {
        confirmationSuppression.value = false
        emissions("fermer")
      },
    },
  )
}

function televerser(evenement: Event) {
  const fichier = (evenement.target as HTMLInputElement).files?.[0]
  if (!fichier || !proprietes.tache) return
  televersement.mutate({ tache: proprietes.tache.id, projet: proprietes.projet, fichier })
  if (champFichier.value) champFichier.value.value = ""
}

function commenter() {
  if (!nouveauCommentaire.value.trim() || !proprietes.tache) return
  mutationCommentaire.mutate(
    { tache: proprietes.tache.id, contenu: nouveauCommentaire.value.trim() },
    { onSuccess: () => (nouveauCommentaire.value = "") },
  )
}
</script>

<template>
  <Dialogue
    :ouvert="ouvert"
    :titre="tache ? 'Détail de la tâche' : 'Nouvelle tâche'"
    large
    @fermer="emissions('fermer')"
  >
    <div class="space-y-4">
      <Champ v-model="formulaire.titre" etiquette="Titre" indication="Que faut-il faire ?" :erreur="erreurs.titre" />
      <Zone v-model="formulaire.description" etiquette="Description" :lignes="4" :erreur="erreurs.description" />

      <div class="grid gap-4 sm:grid-cols-2">
        <Selection v-model="formulaire.colonne" etiquette="Colonne">
          <option v-for="colonne in colonnes" :key="colonne.id" :value="colonne.id">{{ colonne.nom }}</option>
        </Selection>
        <Selection v-model="formulaire.lot" etiquette="Lot de tâches">
          <option value="">Aucun lot</option>
          <option v-for="lot in lots" :key="lot.id" :value="lot.id">
            {{ lot.nom }}<template v-if="lot.echeance"> — {{ formaterDateHeure(lot.echeance) }}</template>
          </option>
        </Selection>
        <Champ v-model="formulaire.points" etiquette="Points" type="number" :erreur="erreurs.points" />
        <Champ v-model="formulaire.echeance" etiquette="Échéance" type="datetime-local" />
      </div>

      <div v-if="etiquettes.length">
        <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Étiquettes</span>
        <div class="mt-2 flex flex-wrap gap-2">
          <button
            v-for="etiquette in etiquettes"
            :key="etiquette.id"
            type="button"
            class="rounded-full border px-3 py-1 text-xs font-medium transition-colors"
            :style="
              formulaire.etiquettes.includes(etiquette.id)
                ? { backgroundColor: etiquette.couleur, borderColor: etiquette.couleur, color: '#ffffff' }
                : { borderColor: etiquette.couleur, color: etiquette.couleur }
            "
            @click="basculer(formulaire.etiquettes, etiquette.id)"
          >
            {{ etiquette.nom }}
          </button>
        </div>
      </div>

      <div>
        <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Affectations</span>
        <div class="mt-2 flex flex-wrap gap-2">
          <button
            v-for="membre in membres"
            :key="membre.utilisateur"
            type="button"
            class="flex items-center gap-2 rounded-full border px-2 py-1 text-xs transition-colors"
            :class="
              formulaire.affectations.includes(membre.utilisateur)
                ? 'border-neutral-900 bg-neutral-900 text-white dark:border-neutral-100 dark:bg-neutral-100 dark:text-neutral-900'
                : 'border-neutral-300 text-neutral-600 hover:border-neutral-500 dark:border-neutral-700 dark:text-neutral-400'
            "
            @click="basculer(formulaire.affectations, membre.utilisateur)"
          >
            <Avatar :nom="membre.nom" :prenom="membre.prenom" petite />
            {{ membre.prenom }} {{ membre.nom }}
          </button>
        </div>
      </div>

      <div v-if="tache">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Images</span>
          <label
            v-if="edition"
            class="cursor-pointer text-xs text-neutral-500 hover:text-neutral-900 hover:underline dark:hover:text-neutral-100"
          >
            Ajouter une image
            <input ref="champFichier" type="file" accept="image/*" class="hidden" @change="televerser" />
          </label>
        </div>
        <div v-if="tache.images.length" class="mt-2 grid grid-cols-3 gap-2">
          <div v-for="image in tache.images" :key="image.id" class="group relative">
            <a :href="image.url" target="_blank" rel="noopener">
              <img :src="image.url" :alt="image.nom" class="h-24 w-full rounded-lg object-cover" />
            </a>
            <button
              v-if="edition"
              class="absolute right-1 top-1 hidden rounded-md bg-neutral-950/70 px-1.5 py-0.5 text-xs text-white group-hover:block"
              @click="suppressionImage.mutate({ id: image.id, projet })"
            >
              Retirer
            </button>
          </div>
        </div>
        <p v-else class="mt-2 text-xs text-neutral-500">Aucune image jointe.</p>
      </div>
      <div v-else>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Images</span>
          <label class="cursor-pointer text-xs text-neutral-500 hover:text-neutral-900 hover:underline dark:hover:text-neutral-100">
            Ajouter des images
            <input type="file" accept="image/*" multiple class="hidden" @change="selectionner" />
          </label>
        </div>
        <div v-if="enAttente.length" class="mt-2 grid grid-cols-3 gap-2">
          <div v-for="(element, indice) in enAttente" :key="element.url" class="group relative">
            <img :src="element.url" :alt="element.fichier.name" class="h-24 w-full rounded-lg object-cover" />
            <button
              class="absolute right-1 top-1 hidden rounded-md bg-neutral-950/70 px-1.5 py-0.5 text-xs text-white group-hover:block"
              @click="retirerEnAttente(indice)"
            >
              Retirer
            </button>
          </div>
        </div>
        <p v-else class="mt-2 text-xs text-neutral-500">Elles seront téléversées à la création de la tâche.</p>
      </div>

      <div v-if="tache" class="border-t border-neutral-200 pt-4 dark:border-neutral-800">
        <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Commentaires</span>
        <ul class="mt-3 space-y-3">
          <li v-for="commentaire in commentaires" :key="commentaire.id" class="flex gap-3">
            <Avatar :nom="commentaire.nom" :prenom="commentaire.prenom" petite />
            <div class="min-w-0 flex-1">
              <p class="text-xs text-neutral-500">
                <span class="font-medium text-neutral-700 dark:text-neutral-300">
                  {{ commentaire.prenom }} {{ commentaire.nom }}
                </span>
                · {{ depuis(commentaire.creation) }}
                <button
                  v-if="commentaire.auteur === profil?.id"
                  class="ml-1 text-red-700 hover:underline dark:text-red-500"
                  @click="suppressionCommentaire.mutate({ id: commentaire.id, tache: tache.id })"
                >
                  Supprimer
                </button>
              </p>
              <p class="text-sm">{{ commentaire.contenu }}</p>
            </div>
          </li>
        </ul>
        <form v-if="edition" class="mt-3 flex gap-2" @submit.prevent="commenter">
          <input
            v-model="nouveauCommentaire"
            placeholder="Écrire un commentaire…"
            class="h-9 flex-1 rounded-lg border border-neutral-300 bg-white px-3 text-sm focus:border-neutral-500 focus:outline-none dark:border-neutral-700 dark:bg-neutral-900"
          />
          <Bouton taille="petite" type="submit" :desactive="mutationCommentaire.isPending.value">Envoyer</Bouton>
        </form>
      </div>

      <div v-if="tache" class="flex items-center justify-between text-xs text-neutral-400">
        <Badge>Créée le {{ formaterDateHeure(tache.creation) }}</Badge>
        <button
          v-if="edition"
          class="text-red-700 hover:underline dark:text-red-500"
          @click="confirmationSuppression = true"
        >
          Supprimer la tâche
        </button>
      </div>

      <div v-if="confirmationSuppression" class="rounded-lg border border-red-300 bg-red-50 p-3 dark:border-red-900 dark:bg-red-950/40">
        <p class="text-sm text-red-800 dark:text-red-300">Supprimer définitivement cette tâche ?</p>
        <div class="mt-2 flex gap-2">
          <Bouton taille="petite" variante="danger" :desactive="suppressionTache.isPending.value" @click="supprimer">
            Oui, supprimer
          </Bouton>
          <Bouton taille="petite" variante="secondaire" @click="confirmationSuppression = false">Annuler</Bouton>
        </div>
      </div>
    </div>

    <template #pied>
      <p v-if="erreurApi" class="mr-auto self-center rounded-lg bg-red-50 px-3 py-1.5 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">
        {{ erreurApi }}
      </p>
      <Bouton variante="secondaire" @click="emissions('fermer')">Fermer</Bouton>
      <Bouton v-if="edition" :desactive="mutationTache.isPending.value" @click="enregistrer">
        {{ tache ? "Enregistrer" : "Créer la tâche" }}
      </Bouton>
    </template>
  </Dialogue>
</template>
