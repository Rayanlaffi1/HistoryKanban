<script setup lang="ts">
import { computed } from "vue"

const proprietes = defineProps<{
  nom: string
  prenom: string
  petite?: boolean
}>()

const initiales = computed(() =>
  ((proprietes.prenom[0] ?? "") + (proprietes.nom[0] ?? "")).toUpperCase() || "?",
)
</script>

<template>
  <span
    :title="`${prenom} ${nom}`.trim()"
    class="inline-flex items-center justify-center rounded-full bg-neutral-300 font-semibold text-neutral-700 ring-2 ring-white dark:bg-neutral-700 dark:text-neutral-200 dark:ring-neutral-900"
    :class="petite ? 'h-6 w-6 text-[10px]' : 'h-8 w-8 text-xs'"
  >
    {{ initiales }}
  </span>
</template>
