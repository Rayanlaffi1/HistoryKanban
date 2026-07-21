<script setup lang="ts">
import { computed } from "vue"
import { estImage, formaterTaille, libelleType } from "@/utilitaires/fichiers"

const proprietes = defineProps<{
  nom: string
  taille: number
  typecontenu: string
  url: string
}>()

const image = computed(() => estImage(proprietes.typecontenu))
const etiquette = computed(() => libelleType(proprietes.typecontenu, proprietes.nom))
const poids = computed(() => formaterTaille(proprietes.taille))
</script>

<template>
  <a v-if="image" :href="url" target="_blank" rel="noopener" class="block">
    <img :src="url" :alt="nom" class="h-24 w-full rounded-lg object-cover" loading="lazy" />
  </a>
  <a
    v-else
    :href="url"
    target="_blank"
    rel="noopener"
    class="flex h-24 flex-col justify-between rounded-lg border border-neutral-200 bg-neutral-50 p-2 transition-colors hover:border-neutral-400 dark:border-neutral-700 dark:bg-neutral-800/60 dark:hover:border-neutral-500"
  >
    <span
      class="w-fit rounded bg-neutral-200 px-1.5 py-0.5 font-mono text-[10px] font-semibold text-neutral-600 dark:bg-neutral-700 dark:text-neutral-300"
    >
      {{ etiquette }}
    </span>
    <span class="line-clamp-2 break-all text-xs text-neutral-700 dark:text-neutral-300">{{ nom }}</span>
    <span class="text-[11px] text-neutral-500">{{ poids }}</span>
  </a>
</template>
