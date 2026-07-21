<script setup lang="ts">
defineProps<{
  etiquette?: string
}>()

const modele = defineModel<boolean>({ default: false })
</script>

<template>
  <label class="flex cursor-pointer items-center justify-between gap-4">
    <span v-if="etiquette" class="text-sm text-neutral-700 dark:text-neutral-300">{{ etiquette }}</span>
    <button
      type="button"
      role="switch"
      :aria-checked="modele"
      class="relative h-6 w-11 rounded-full transition-colors"
      :class="modele ? 'bg-neutral-900 dark:bg-neutral-100' : 'bg-neutral-300 dark:bg-neutral-700'"
      @click="modele = !modele"
    >
      <span
        class="absolute top-0.5 h-5 w-5 rounded-full bg-white transition-all dark:bg-neutral-900"
        :class="modele ? 'left-[22px]' : 'left-0.5'"
      ></span>
    </button>
  </label>
</template>
