<script setup lang="ts">
import { computed } from "vue"
import type { Etiquette, Lot, Membre, Tache } from "@/api/types"
import { estDepassee, formaterDate, formaterDateHeure } from "@/utilitaires/dates"
import { texteBrut } from "@/utilitaires/html"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"

export type Densite = "compacte" | "defaut" | "detaillee"

const proprietes = defineProps<{
  tache: Tache
  etiquettes: Record<string, Etiquette>
  membres: Record<string, Membre>
  lots: Record<string, Lot>
  densite: Densite
}>()

defineEmits<{ ouvrir: [] }>()

const lot = computed(() => (proprietes.tache.lot ? proprietes.lots[proprietes.tache.lot] : undefined))
const affiches = computed(() => proprietes.tache.affectations.slice(0, 4))
const apercu = computed(() => proprietes.tache.images[0])
const galerie = computed(() => proprietes.tache.images.slice(0, 3))
const resume = computed(() => texteBrut(proprietes.tache.description))
</script>

<template>
  <article
    v-if="densite === 'compacte'"
    class="flex cursor-pointer items-center justify-between gap-2 rounded-lg border border-neutral-200 bg-white px-3 py-1.5 shadow-sm transition-shadow hover:shadow-md dark:border-neutral-700 dark:bg-neutral-800"
    @click="$emit('ouvrir')"
  >
    <h4 class="truncate text-sm font-medium">{{ tache.titre }}</h4>
    <div class="flex shrink-0 items-center gap-1.5">
      <span
        v-if="tache.echeance && estDepassee(tache.echeance)"
        class="h-2 w-2 rounded-full bg-red-600"
        title="Échéance dépassée"
      ></span>
      <Badge v-if="tache.points > 0">{{ tache.points }}</Badge>
    </div>
  </article>

  <article
    v-else
    class="cursor-pointer rounded-lg border border-neutral-200 bg-white p-3 shadow-sm transition-shadow hover:shadow-md dark:border-neutral-700 dark:bg-neutral-800"
    @click="$emit('ouvrir')"
  >
    <template v-if="densite === 'defaut'">
      <img
        v-if="apercu"
        :src="apercu.url"
        :alt="apercu.nom"
        class="mb-2 h-28 w-full rounded-md object-cover"
        loading="lazy"
      />
    </template>
    <template v-else>
      <div v-if="galerie.length" class="mb-2 grid gap-1" :class="galerie.length > 1 ? 'grid-cols-3' : ''">
        <img
          v-for="image in galerie"
          :key="image.id"
          :src="image.url"
          :alt="image.nom"
          class="w-full rounded-md object-cover"
          :class="galerie.length > 1 ? 'h-16' : 'h-28'"
          loading="lazy"
        />
      </div>
    </template>

    <div class="flex items-start justify-between gap-2">
      <h4 class="text-sm font-medium leading-snug">{{ tache.titre }}</h4>
      <Badge v-if="tache.points > 0">{{ tache.points }} pt{{ tache.points > 1 ? "s" : "" }}</Badge>
    </div>

    <p v-if="densite === 'detaillee' && resume" class="mt-1 line-clamp-3 text-xs text-neutral-500 dark:text-neutral-400">
      {{ resume }}
    </p>

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
        {{ densite === "detaillee" ? formaterDateHeure(tache.echeance) : formaterDate(tache.echeance) }}
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

    <p v-if="densite === 'detaillee'" class="mt-2 border-t border-neutral-100 pt-1.5 text-[11px] text-neutral-400 dark:border-neutral-700">
      Créée le {{ formaterDate(tache.creation) }} · {{ tache.images.length }} image{{ tache.images.length > 1 ? "s" : "" }}
    </p>
  </article>
</template>
