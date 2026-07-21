<script setup lang="ts">
defineProps<{
  etiquette?: string
}>()

const modele = defineModel<string>({ default: "" })
</script>

<template>
  <label class="block space-y-1.5">
    <span v-if="etiquette" class="text-sm font-medium text-neutral-700 dark:text-neutral-300">{{ etiquette }}</span>
    <select
      v-model="modele"
      class="h-10 w-full rounded-lg border border-neutral-300 bg-white px-3 text-sm text-neutral-900 focus:border-neutral-500 focus:outline-none dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100"
    >
      <slot />
    </select>
  </label>
</template>
