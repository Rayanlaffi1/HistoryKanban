<script setup lang="ts">
defineProps<{
  etiquette?: string
}>()

const modele = defineModel<string>({ default: "#737373" })

const predefinies = [
  "#737373",
  "#525252",
  "#404040",
  "#dc2626",
  "#ea580c",
  "#ca8a04",
  "#16a34a",
  "#0891b2",
  "#2563eb",
  "#7c3aed",
  "#db2777",
]
</script>

<template>
  <div class="space-y-1.5">
    <span v-if="etiquette" class="block text-sm font-medium text-neutral-700 dark:text-neutral-300">{{ etiquette }}</span>
    <div class="flex items-center gap-3 rounded-lg border border-neutral-300 bg-white px-3 py-2.5 shadow-sm dark:border-neutral-700 dark:bg-neutral-900">
      <label
        class="relative h-8 w-8 shrink-0 cursor-pointer overflow-hidden rounded-lg border border-neutral-300 shadow-inner dark:border-neutral-600"
        title="Couleur personnalisée"
      >
        <span class="absolute inset-0" :style="{ backgroundColor: modele }"></span>
        <input v-model="modele" type="color" class="absolute inset-0 cursor-pointer opacity-0" />
      </label>
      <div class="flex flex-wrap gap-1.5">
        <button
          v-for="couleur in predefinies"
          :key="couleur"
          type="button"
          class="h-5 w-5 rounded-full transition-transform hover:scale-110"
          :class="modele === couleur && 'ring-2 ring-neutral-900 ring-offset-1 ring-offset-white dark:ring-neutral-100 dark:ring-offset-neutral-900'"
          :style="{ backgroundColor: couleur }"
          :title="couleur"
          @click="modele = couleur"
        ></button>
      </div>
      <span class="ml-auto shrink-0 font-mono text-xs uppercase text-neutral-400">{{ modele }}</span>
    </div>
  </div>
</template>
