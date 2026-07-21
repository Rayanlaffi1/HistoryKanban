<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue"
import type { Colonne, Etiquette, Lot, Membre, Tache } from "@/api/types"
import {
  televerserFichier,
  utiliserMutationTache,
  utiliserSuppressionImage,
  utiliserSuppressionTache,
  utiliserTeleversementImage,
} from "@/api/requetes"
import { depuisChampDate, formaterDateHeure, versChampDate } from "@/utilitaires/dates"
import { formulaireTache, valider } from "@/utilitaires/validation"
import { extraireErreur } from "@/utilitaires/erreurs"
import { assainir } from "@/utilitaires/html"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import ZoneRiche from "@/composants/ui/ZoneRiche.vue"
import Selection from "@/composants/ui/Selection.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"
import SectionFormulaire from "@/composants/ui/SectionFormulaire.vue"
import CommentairesTache from "@/composants/kanban/CommentairesTache.vue"
import ActiviteTache from "@/composants/kanban/ActiviteTache.vue"

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

type Onglet = "detail" | "commentaires" | "activite"
const onglet = ref<Onglet>("detail")
const onglets: { valeur: Onglet; libelle: string }[] = [
  { valeur: "detail", libelle: "Détail" },
  { valeur: "commentaires", libelle: "Commentaires" },
  { valeur: "activite", libelle: "Activité" },
]

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

const erreurs = ref<Record<string, string>>({})
const erreurApi = ref("")
const confirmationSuppression = ref(false)
const champFichier = ref<HTMLInputElement | null>(null)

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

watch(
  () => [proprietes.ouvert, proprietes.tache] as const,
  ([ouvert], anciens) => {
    if (!ouvert) return
    const ancienOuvert = anciens?.[0] ?? false
    if (!ancienOuvert) {
      onglet.value = "detail"
      erreurs.value = {}
      erreurApi.value = ""
      confirmationSuppression.value = false
      viderEnAttente()
    }
    formulaire.titre = proprietes.tache?.titre ?? ""
    formulaire.description = proprietes.tache?.description ?? ""
    formulaire.colonne = proprietes.tache?.colonne ?? proprietes.colonneInitiale
    formulaire.lot = proprietes.tache?.lot ?? ""
    formulaire.points = String(proprietes.tache?.points ?? 0)
    formulaire.echeance = versChampDate(proprietes.tache?.echeance)
    formulaire.affectations = [...(proprietes.tache?.affectations ?? [])]
    formulaire.etiquettes = [...(proprietes.tache?.etiquettes ?? [])]
  },
  { immediate: true },
)

const nomColonne = computed(
  () => proprietes.colonnes.find((colonne) => colonne.id === proprietes.tache?.colonne)?.nom ?? "",
)
const lotCourant = computed(() => proprietes.lots.find((lot) => lot.id === proprietes.tache?.lot))
const etiquettesCourantes = computed(() =>
  proprietes.etiquettes.filter((etiquette) => proprietes.tache?.etiquettes.includes(etiquette.id)),
)
const membresAffectes = computed(() =>
  proprietes.membres.filter((membre) => proprietes.tache?.affectations.includes(membre.utilisateur)),
)

const mutationTache = utiliserMutationTache()
const suppressionTache = utiliserSuppressionTache()
const televersement = utiliserTeleversementImage()
const suppressionImage = utiliserSuppressionImage()

function basculer(liste: string[], valeur: string) {
  const indice = liste.indexOf(valeur)
  if (indice >= 0) liste.splice(indice, 1)
  else liste.push(valeur)
}

function envoyerImageContenu(fichier: File) {
  return televerserFichier(proprietes.projet, fichier)
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

function televerser(evenement: Event) {
  const fichier = (evenement.target as HTMLInputElement).files?.[0]
  if (!fichier || !proprietes.tache) return
  televersement.mutate({ tache: proprietes.tache.id, projet: proprietes.projet, fichier })
  if (champFichier.value) champFichier.value.value = ""
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
      description: resultat.donnees.description === "" ? "" : assainir(resultat.donnees.description),
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
</script>

<template>
  <Dialogue
    :ouvert="ouvert"
    :titre="tache ? 'Détail de la tâche' : 'Nouvelle tâche'"
    geant
    @fermer="emissions('fermer')"
  >
    <div v-if="tache" class="mb-5 flex gap-1 rounded-lg bg-neutral-100 p-1 dark:bg-neutral-800">
      <button
        v-for="option in onglets"
        :key="option.valeur"
        type="button"
        class="flex-1 rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
        :class="
          onglet === option.valeur
            ? 'bg-white text-neutral-900 shadow-sm dark:bg-neutral-900 dark:text-neutral-100'
            : 'text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100'
        "
        @click="onglet = option.valeur"
      >
        {{ option.libelle }}
      </button>
    </div>

    <div v-show="!tache || onglet === 'detail'">
      <div v-if="!edition && tache" class="space-y-5">
        <div class="flex items-start justify-between gap-3">
          <h3 class="text-lg font-semibold leading-snug">{{ tache.titre }}</h3>
          <Badge v-if="tache.points > 0">{{ tache.points }} pt{{ tache.points > 1 ? "s" : "" }}</Badge>
        </div>

        <div v-if="tache.description" class="texteriche text-sm" v-html="assainir(tache.description)"></div>
        <p v-else class="text-sm italic text-neutral-500">Aucune description.</p>

        <dl class="grid grid-cols-2 gap-x-4 gap-y-2 rounded-lg bg-neutral-100 p-4 text-sm dark:bg-neutral-800/60">
          <dt class="text-neutral-500">Colonne</dt>
          <dd class="font-medium">{{ nomColonne || "—" }}</dd>
          <dt class="text-neutral-500">Lot</dt>
          <dd class="font-medium">
            <template v-if="lotCourant">
              {{ lotCourant.nom }}
              <span v-if="lotCourant.echeance" class="text-neutral-500">
                — {{ formaterDateHeure(lotCourant.echeance) }}
              </span>
            </template>
            <template v-else>—</template>
          </dd>
          <dt class="text-neutral-500">Échéance</dt>
          <dd class="font-medium">{{ tache.echeance ? formaterDateHeure(tache.echeance) : "—" }}</dd>
          <dt class="text-neutral-500">Créée le</dt>
          <dd class="font-medium">{{ formaterDateHeure(tache.creation) }}</dd>
        </dl>

        <div v-if="etiquettesCourantes.length">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Étiquettes</span>
          <div class="mt-2 flex flex-wrap gap-1.5">
            <Badge v-for="etiquette in etiquettesCourantes" :key="etiquette.id" :couleur="etiquette.couleur">
              {{ etiquette.nom }}
            </Badge>
          </div>
        </div>

        <div v-if="membresAffectes.length">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Affectations</span>
          <div class="mt-2 flex flex-wrap gap-2">
            <span
              v-for="membre in membresAffectes"
              :key="membre.utilisateur"
              class="flex items-center gap-2 rounded-full border border-neutral-300 px-2 py-1 text-xs text-neutral-700 dark:border-neutral-700 dark:text-neutral-300"
            >
              <Avatar :nom="membre.nom" :prenom="membre.prenom" petite />
              {{ membre.prenom }} {{ membre.nom }}
            </span>
          </div>
        </div>

        <div v-if="tache.images.length">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Images</span>
          <div class="mt-2 grid grid-cols-3 gap-2">
            <a v-for="image in tache.images" :key="image.id" :href="image.url" target="_blank" rel="noopener">
              <img :src="image.url" :alt="image.nom" class="h-24 w-full rounded-lg object-cover" />
            </a>
          </div>
        </div>
      </div>

      <div v-else class="space-y-4">
        <Champ v-model="formulaire.titre" etiquette="Titre" obligatoire indication="Que faut-il faire ?" :erreur="erreurs.titre" />
        <ZoneRiche
          v-model="formulaire.description"
          etiquette="Description"
          indication="Décrivez la tâche, insérez des images…"
          :erreur="erreurs.description"
          :televerser="envoyerImageContenu"
        />

        <SectionFormulaire titre="Planification">
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
            <Champ v-model="formulaire.echeance" etiquette="Échéance" type="datetime-local" />
            <Champ v-model="formulaire.points" etiquette="Points" type="number" :erreur="erreurs.points" />
          </div>
        </SectionFormulaire>

        <SectionFormulaire v-if="etiquettes.length || membres.length" titre="Organisation">
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
          <div v-if="membres.length">
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
        </SectionFormulaire>

        <SectionFormulaire titre="Images">
          <template v-if="tache">
            <div v-if="tache.images.length" class="grid grid-cols-3 gap-2">
              <div v-for="image in tache.images" :key="image.id" class="group relative">
                <a :href="image.url" target="_blank" rel="noopener">
                  <img :src="image.url" :alt="image.nom" class="h-24 w-full rounded-lg object-cover" />
                </a>
                <button
                  class="absolute right-1 top-1 hidden rounded-md bg-neutral-950/70 px-1.5 py-0.5 text-xs text-white group-hover:block"
                  @click="suppressionImage.mutate({ id: image.id, projet })"
                >
                  Retirer
                </button>
              </div>
            </div>
            <p v-else class="text-xs text-neutral-500">Aucune image jointe.</p>
            <label
              class="inline-flex h-9 cursor-pointer items-center rounded-lg border border-dashed border-neutral-300 px-3 text-xs font-medium text-neutral-600 transition-colors hover:border-neutral-500 hover:text-neutral-900 dark:border-neutral-600 dark:text-neutral-400 dark:hover:border-neutral-400 dark:hover:text-neutral-100"
            >
              + Ajouter une image
              <input ref="champFichier" type="file" accept="image/*" class="hidden" @change="televerser" />
            </label>
          </template>
          <template v-else>
            <div v-if="enAttente.length" class="grid grid-cols-3 gap-2">
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
            <p v-else class="text-xs text-neutral-500">Elles seront téléversées à la création de la tâche.</p>
            <label
              class="inline-flex h-9 cursor-pointer items-center rounded-lg border border-dashed border-neutral-300 px-3 text-xs font-medium text-neutral-600 transition-colors hover:border-neutral-500 hover:text-neutral-900 dark:border-neutral-600 dark:text-neutral-400 dark:hover:border-neutral-400 dark:hover:text-neutral-100"
            >
              + Ajouter des images
              <input type="file" accept="image/*" multiple class="hidden" @change="selectionner" />
            </label>
          </template>
        </SectionFormulaire>

        <div v-if="tache" class="flex justify-end">
          <button
            class="text-xs text-red-700 hover:underline dark:text-red-500"
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
    </div>

    <div v-if="tache" v-show="onglet === 'commentaires'">
      <CommentairesTache :tache="tache" :projet="projet" :edition="edition" />
    </div>

    <div v-if="tache" v-show="onglet === 'activite'">
      <ActiviteTache :tache="tache.id" />
    </div>

    <template #pied>
      <p v-if="erreurApi" class="mr-auto self-center rounded-lg bg-red-50 px-3 py-1.5 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">
        {{ erreurApi }}
      </p>
      <Bouton variante="secondaire" @click="emissions('fermer')">Fermer</Bouton>
      <Bouton
        v-if="edition && (!tache || onglet === 'detail')"
        :desactive="mutationTache.isPending.value"
        @click="enregistrer"
      >
        {{ tache ? "Enregistrer" : "Créer la tâche" }}
      </Bouton>
    </template>
  </Dialogue>
</template>
