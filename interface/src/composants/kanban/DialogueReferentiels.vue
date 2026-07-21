<script setup lang="ts">
import { reactive, ref } from "vue"
import { formulaireEtiquette, formulaireLot, valider } from "@/utilitaires/validation"
import type { Etiquette, Lot } from "@/api/types"
import {
  utiliserMutationEtiquette,
  utiliserMutationLot,
  utiliserSuppressionEtiquette,
  utiliserSuppressionLot,
} from "@/api/requetes"
import { depuisChampDate, formaterDateHeure } from "@/utilitaires/dates"
import Bouton from "@/composants/ui/Bouton.vue"
import Badge from "@/composants/ui/Badge.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import SectionFormulaire from "@/composants/ui/SectionFormulaire.vue"

defineProps<{
  ouvert: boolean
  projet: string
  etiquettes: Etiquette[]
  lots: Lot[]
}>()

const emissions = defineEmits<{ fermer: [] }>()

const nouvelleEtiquette = reactive({ nom: "", couleur: "#737373" })
const nouveauLot = reactive({ nom: "", couleur: "#525252", echeance: "" })

const mutationEtiquette = utiliserMutationEtiquette()
const suppressionEtiquette = utiliserSuppressionEtiquette()
const mutationLot = utiliserMutationLot()
const suppressionLot = utiliserSuppressionLot()

const erreurEtiquette = ref("")
const erreurLot = ref("")

function creerEtiquette(projet: string) {
  const resultat = valider(formulaireEtiquette, { nom: nouvelleEtiquette.nom })
  erreurEtiquette.value = resultat.erreurs.nom ?? ""
  if (!resultat.donnees) return
  mutationEtiquette.mutate(
    { projet, nom: resultat.donnees.nom, couleur: nouvelleEtiquette.couleur },
    { onSuccess: () => (nouvelleEtiquette.nom = "") },
  )
}

function creerLot(projet: string) {
  const resultat = valider(formulaireLot, { nom: nouveauLot.nom })
  erreurLot.value = resultat.erreurs.nom ?? ""
  if (!resultat.donnees) return
  mutationLot.mutate(
    {
      projet,
      nom: resultat.donnees.nom,
      couleur: nouveauLot.couleur,
      echeance: depuisChampDate(nouveauLot.echeance),
    },
    {
      onSuccess: () => {
        nouveauLot.nom = ""
        nouveauLot.echeance = ""
      },
    },
  )
}

const classeChamp =
  "h-10 rounded-lg border border-neutral-300 bg-white px-3 text-sm shadow-sm transition placeholder:text-neutral-400 focus:border-neutral-500 focus:outline-none focus:ring-2 focus:ring-neutral-900/10 dark:border-neutral-700 dark:bg-neutral-900 dark:placeholder:text-neutral-500 dark:focus:border-neutral-400 dark:focus:ring-neutral-100/10"

const classeCouleur =
  "relative h-10 w-12 shrink-0 cursor-pointer overflow-hidden rounded-lg border border-neutral-300 shadow-sm dark:border-neutral-700"
</script>

<template>
  <Dialogue :ouvert="ouvert" titre="Étiquettes et lots" large @fermer="emissions('fermer')">
    <div class="space-y-5">
      <SectionFormulaire titre="Étiquettes">
        <div class="flex flex-wrap gap-2">
          <span
            v-for="etiquette in etiquettes"
            :key="etiquette.id"
            class="inline-flex items-center gap-1.5 rounded-full border border-neutral-200 py-1 pl-1.5 pr-2 dark:border-neutral-700"
          >
            <Badge :couleur="etiquette.couleur">{{ etiquette.nom }}</Badge>
            <button
              class="text-xs text-neutral-400 transition-colors hover:text-red-700 dark:hover:text-red-500"
              title="Supprimer"
              @click="suppressionEtiquette.mutate({ id: etiquette.id, projet })"
            >
              ✕
            </button>
          </span>
          <span v-if="!etiquettes.length" class="text-sm text-neutral-500">Aucune étiquette pour le moment.</span>
        </div>
        <form class="flex items-center gap-2 border-t border-neutral-100 pt-3 dark:border-neutral-800" @submit.prevent="creerEtiquette(projet)">
          <input v-model="nouvelleEtiquette.nom" placeholder="Nouvelle étiquette" :class="classeChamp" class="flex-1" />
          <label :class="classeCouleur" title="Couleur">
            <span class="absolute inset-0" :style="{ backgroundColor: nouvelleEtiquette.couleur }"></span>
            <input v-model="nouvelleEtiquette.couleur" type="color" class="absolute inset-0 cursor-pointer opacity-0" />
          </label>
          <Bouton type="submit" :desactive="mutationEtiquette.isPending.value">Ajouter</Bouton>
        </form>
        <p v-if="erreurEtiquette" class="text-xs font-medium text-red-700 dark:text-red-500">{{ erreurEtiquette }}</p>
      </SectionFormulaire>

      <SectionFormulaire titre="Lots de tâches">
        <p class="text-xs text-neutral-500">Regroupez des tâches sous une même échéance commune.</p>
        <ul class="space-y-2">
          <li
            v-for="lot in lots"
            :key="lot.id"
            class="flex items-center justify-between rounded-lg border border-neutral-200 bg-neutral-50 px-3 py-2.5 dark:border-neutral-700 dark:bg-neutral-800/50"
          >
            <div class="flex min-w-0 items-center gap-2.5">
              <span class="h-3 w-3 shrink-0 rounded-full" :style="{ backgroundColor: lot.couleur }"></span>
              <span class="truncate text-sm font-medium">{{ lot.nom }}</span>
              <span v-if="lot.echeance" class="shrink-0 text-xs text-neutral-500">
                échéance {{ formaterDateHeure(lot.echeance) }}
              </span>
            </div>
            <button
              class="shrink-0 text-xs text-neutral-400 transition-colors hover:text-red-700 dark:hover:text-red-500"
              @click="suppressionLot.mutate({ id: lot.id, projet })"
            >
              Supprimer
            </button>
          </li>
          <li v-if="!lots.length" class="text-sm text-neutral-500">Aucun lot pour le moment.</li>
        </ul>
        <form class="space-y-2 border-t border-neutral-100 pt-3 dark:border-neutral-800" @submit.prevent="creerLot(projet)">
          <div class="flex items-center gap-2">
            <input v-model="nouveauLot.nom" placeholder="Nouveau lot" :class="classeChamp" class="flex-1" />
            <label :class="classeCouleur" title="Couleur">
              <span class="absolute inset-0" :style="{ backgroundColor: nouveauLot.couleur }"></span>
              <input v-model="nouveauLot.couleur" type="color" class="absolute inset-0 cursor-pointer opacity-0" />
            </label>
          </div>
          <div class="flex items-center gap-2">
            <input v-model="nouveauLot.echeance" type="datetime-local" :class="classeChamp" class="flex-1" />
            <Bouton type="submit" :desactive="mutationLot.isPending.value">Ajouter</Bouton>
          </div>
        </form>
        <p v-if="erreurLot" class="text-xs font-medium text-red-700 dark:text-red-500">{{ erreurLot }}</p>
      </SectionFormulaire>
    </div>
    <template #pied>
      <Bouton variante="secondaire" @click="emissions('fermer')">Fermer</Bouton>
    </template>
  </Dialogue>
</template>
