<script setup lang="ts">
defineProps<{
  etiquette?: string
  erreur?: string
}>()

const modele = defineModel<string>({ default: "" })
</script>

<template>
  <label class="block space-y-1.5">
    <span v-if="etiquette" class="block text-sm font-medium text-neutral-700 dark:text-neutral-300">{{ etiquette }}</span>
    <span class="relative block">
      <select
        v-model="modele"
        class="h-10 w-full appearance-none rounded-lg border bg-white pl-3 pr-9 text-sm text-neutral-900 shadow-sm transition focus:outline-none focus:ring-2 dark:bg-neutral-900 dark:text-neutral-100"
        :class="
          erreur
            ? 'border-red-500 focus:border-red-600 focus:ring-red-600/15 dark:border-red-700'
            : 'border-neutral-300 focus:border-neutral-500 focus:ring-neutral-900/10 dark:border-neutral-700 dark:focus:border-neutral-400 dark:focus:ring-neutral-100/10'
        "
      >
        <slot />
      </select>
      <svg
        class="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-neutral-400"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="m6 9 6 6 6-6" />
      </svg>
    </span>
    <span v-if="erreur" class="block text-xs font-medium text-red-700 dark:text-red-500">{{ erreur }}</span>
  </label>
</template>
