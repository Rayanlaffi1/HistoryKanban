<script setup lang="ts">
defineProps<{
  etiquette?: string
  lignes?: number
  indication?: string
}>()

const modele = defineModel<string>({ default: "" })
</script>

<template>
  <label class="block space-y-1.5">
    <span v-if="etiquette" class="text-sm font-medium text-neutral-700 dark:text-neutral-300">{{ etiquette }}</span>
    <textarea
      v-model="modele"
      :rows="lignes ?? 3"
      :placeholder="indication"
      class="w-full rounded-lg border border-neutral-300 bg-white px-3 py-2 text-sm text-neutral-900 placeholder:text-neutral-400 focus:border-neutral-500 focus:outline-none dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder:text-neutral-500"
    ></textarea>
  </label>
</template>
