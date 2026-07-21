<script setup lang="ts">
import { computed } from "vue"
import type { Etiquette, FiltreTaches, Lot, Membre } from "@/api/types"

defineProps<{
  membres: Membre[]
  etiquettes: Etiquette[]
  lots: Lot[]
}>()

const filtre = defineModel<FiltreTaches>({ required: true })

const actif = computed(() =>
  Boolean(
    filtre.value.texte ||
      filtre.value.membre ||
      filtre.value.etiquette ||
      filtre.value.lot ||
      filtre.value.echeance ||
      filtre.value.pointsmin !== undefined ||
      filtre.value.pointsmax !== undefined,
  ),
)

function changer(cle: keyof FiltreTaches, valeur: string) {
  const copie = { ...filtre.value }
  if (cle === "pointsmin" || cle === "pointsmax") {
    copie[cle] = valeur === "" ? undefined : Number(valeur)
  } else {
    copie[cle] = valeur === "" ? undefined : valeur
  }
  filtre.value = copie
}

function reinitialiser() {
  filtre.value = {}
}

const classeChamp =
  "h-9 rounded-lg border border-neutral-300 bg-white px-2.5 text-sm text-neutral-900 focus:border-neutral-500 focus:outline-none dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100"
</script>

<template>
  <div class="flex flex-wrap items-center gap-2">
    <input
      :value="filtre.texte ?? ''"
      type="search"
      placeholder="Rechercher…"
      :class="classeChamp"
      class="w-44"
      @input="changer('texte', ($event.target as HTMLInputElement).value)"
    />
    <select :value="filtre.membre ?? ''" :class="classeChamp" @change="changer('membre', ($event.target as HTMLSelectElement).value)">
      <option value="">Tous les membres</option>
      <option v-for="membre in membres" :key="membre.utilisateur" :value="membre.utilisateur">
        {{ membre.prenom }} {{ membre.nom }}
      </option>
    </select>
    <select :value="filtre.etiquette ?? ''" :class="classeChamp" @change="changer('etiquette', ($event.target as HTMLSelectElement).value)">
      <option value="">Toutes les étiquettes</option>
      <option v-for="etiquette in etiquettes" :key="etiquette.id" :value="etiquette.id">{{ etiquette.nom }}</option>
    </select>
    <select :value="filtre.lot ?? ''" :class="classeChamp" @change="changer('lot', ($event.target as HTMLSelectElement).value)">
      <option value="">Tous les lots</option>
      <option v-for="lot in lots" :key="lot.id" :value="lot.id">{{ lot.nom }}</option>
    </select>
    <select :value="filtre.echeance ?? ''" :class="classeChamp" @change="changer('echeance', ($event.target as HTMLSelectElement).value)">
      <option value="">Toutes les échéances</option>
      <option value="depassee">En retard</option>
      <option value="semaine">Sous 7 jours</option>
      <option value="sans">Sans échéance</option>
    </select>
    <input
      :value="filtre.pointsmin ?? ''"
      type="number"
      min="0"
      placeholder="Pts min"
      :class="classeChamp"
      class="w-20"
      @input="changer('pointsmin', ($event.target as HTMLInputElement).value)"
    />
    <input
      :value="filtre.pointsmax ?? ''"
      type="number"
      min="0"
      placeholder="Pts max"
      :class="classeChamp"
      class="w-20"
      @input="changer('pointsmax', ($event.target as HTMLInputElement).value)"
    />
    <button
      v-if="actif"
      class="text-sm text-neutral-500 hover:text-neutral-900 hover:underline dark:hover:text-neutral-100"
      @click="reinitialiser"
    >
      Réinitialiser
    </button>
  </div>
</template>
