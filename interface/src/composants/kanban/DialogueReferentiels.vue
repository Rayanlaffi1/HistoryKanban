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
  "h-9 rounded-lg border border-neutral-300 bg-white px-2.5 text-sm focus:border-neutral-500 focus:outline-none dark:border-neutral-700 dark:bg-neutral-900"
</script>

<template>
  <Dialogue :ouvert="ouvert" titre="Étiquettes et lots" large @fermer="emissions('fermer')">
    <div class="space-y-6">
      <section>
        <h3 class="text-sm font-semibold">Étiquettes</h3>
        <div class="mt-3 flex flex-wrap gap-2">
          <span v-for="etiquette in etiquettes" :key="etiquette.id" class="inline-flex items-center gap-1">
            <Badge :couleur="etiquette.couleur">{{ etiquette.nom }}</Badge>
            <button
              class="text-xs text-neutral-400 hover:text-red-700 dark:hover:text-red-500"
              title="Supprimer"
              @click="suppressionEtiquette.mutate({ id: etiquette.id, projet })"
            >
              ✕
            </button>
          </span>
          <span v-if="!etiquettes.length" class="text-sm text-neutral-500">Aucune étiquette.</span>
        </div>
        <form class="mt-3 flex items-center gap-2" @submit.prevent="creerEtiquette(projet)">
          <input v-model="nouvelleEtiquette.nom" placeholder="Nouvelle étiquette" :class="classeChamp" class="flex-1" />
          <input v-model="nouvelleEtiquette.couleur" type="color" class="h-9 w-12 cursor-pointer rounded-lg border border-neutral-300 dark:border-neutral-700" />
          <Bouton taille="petite" type="submit" :desactive="mutationEtiquette.isPending.value">Ajouter</Bouton>
        </form>
        <p v-if="erreurEtiquette" class="mt-1 text-xs text-red-700 dark:text-red-500">{{ erreurEtiquette }}</p>
      </section>

      <section class="border-t border-neutral-200 pt-5 dark:border-neutral-800">
        <h3 class="text-sm font-semibold">Lots de tâches</h3>
        <p class="mt-1 text-xs text-neutral-500">
          Regroupez des tâches sous une même échéance commune.
        </p>
        <ul class="mt-3 space-y-2">
          <li
            v-for="lot in lots"
            :key="lot.id"
            class="flex items-center justify-between rounded-lg border border-neutral-200 px-3 py-2 dark:border-neutral-800"
          >
            <div class="flex items-center gap-2">
              <span class="h-2.5 w-2.5 rounded-full" :style="{ backgroundColor: lot.couleur }"></span>
              <span class="text-sm font-medium">{{ lot.nom }}</span>
              <span v-if="lot.echeance" class="text-xs text-neutral-500">
                échéance {{ formaterDateHeure(lot.echeance) }}
              </span>
            </div>
            <button
              class="text-xs text-neutral-400 hover:text-red-700 dark:hover:text-red-500"
              @click="suppressionLot.mutate({ id: lot.id, projet })"
            >
              Supprimer
            </button>
          </li>
          <li v-if="!lots.length" class="text-sm text-neutral-500">Aucun lot.</li>
        </ul>
        <form class="mt-3 flex flex-wrap items-center gap-2" @submit.prevent="creerLot(projet)">
          <input v-model="nouveauLot.nom" placeholder="Nouveau lot" :class="classeChamp" class="flex-1" />
          <input v-model="nouveauLot.echeance" type="datetime-local" :class="classeChamp" />
          <input v-model="nouveauLot.couleur" type="color" class="h-9 w-12 cursor-pointer rounded-lg border border-neutral-300 dark:border-neutral-700" />
          <Bouton taille="petite" type="submit" :desactive="mutationLot.isPending.value">Ajouter</Bouton>
        </form>
        <p v-if="erreurLot" class="mt-1 text-xs text-red-700 dark:text-red-500">{{ erreurLot }}</p>
      </section>
    </div>
    <template #pied>
      <Bouton variante="secondaire" @click="emissions('fermer')">Fermer</Bouton>
    </template>
  </Dialogue>
</template>
