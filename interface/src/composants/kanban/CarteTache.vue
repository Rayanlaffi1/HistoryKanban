<script setup lang="ts">
import { computed } from "vue"
import type { Etiquette, Lot, Membre, Tache } from "@/api/types"
import { estDepassee, formaterDate } from "@/utilitaires/dates"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"

const proprietes = defineProps<{
  tache: Tache
  etiquettes: Record<string, Etiquette>
  membres: Record<string, Membre>
  lots: Record<string, Lot>
}>()

defineEmits<{ ouvrir: [] }>()

const lot = computed(() => (proprietes.tache.lot ? proprietes.lots[proprietes.tache.lot] : undefined))
const affiches = computed(() => proprietes.tache.affectations.slice(0, 4))
const apercu = computed(() => proprietes.tache.images[0])
</script>

<template>
  <article
    class="cursor-pointer rounded-lg border border-neutral-200 bg-white p-3 shadow-sm transition-shadow hover:shadow-md dark:border-neutral-700 dark:bg-neutral-800"
    @click="$emit('ouvrir')"
  >
    <img
      v-if="apercu"
      :src="apercu.url"
      :alt="apercu.nom"
      class="mb-2 h-28 w-full rounded-md object-cover"
      loading="lazy"
    />
    <div class="flex items-start justify-between gap-2">
      <h4 class="text-sm font-medium leading-snug">{{ tache.titre }}</h4>
      <Badge v-if="tache.points > 0">{{ tache.points }} pt{{ tache.points > 1 ? "s" : "" }}</Badge>
    </div>
    <div v-if="tache.etiquettes.length || lot" class="mt-2 flex flex-wrap gap-1">
      <Badge v-if="lot" :couleur="lot.couleur">{{ lot.nom }}</Badge>
      <Badge
        v-for="identifiant in tache.etiquettes"
        :key="identifiant"
        :couleur="etiquettes[identifiant]?.couleur"
      >
        {{ etiquettes[identifiant]?.nom ?? "?" }}
      </Badge>
    </div>
    <div class="mt-2 flex items-center justify-between">
      <span
        v-if="tache.echeance"
        class="text-xs"
        :class="estDepassee(tache.echeance) ? 'font-semibold text-red-700 dark:text-red-500' : 'text-neutral-500'"
      >
        {{ formaterDate(tache.echeance) }}
      </span>
      <span v-else></span>
      <div class="flex -space-x-2">
        <Avatar
          v-for="identifiant in affiches"
          :key="identifiant"
          :nom="membres[identifiant]?.nom ?? '?'"
          :prenom="membres[identifiant]?.prenom ?? ''"
          petite
        />
      </div>
    </div>
  </article>
</template>
