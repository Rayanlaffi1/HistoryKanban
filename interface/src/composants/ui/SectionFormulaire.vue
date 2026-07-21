<script setup lang="ts">
defineProps<{
  titre?: string
}>()
</script>

<template>
  <fieldset class="rounded-xl border border-neutral-200 px-4 pb-4 dark:border-neutral-800" :class="titre ? 'pt-1' : 'pt-4'">
    <legend v-if="titre" class="px-2 text-xs font-semibold uppercase tracking-wider text-neutral-500 dark:text-neutral-400">
      {{ titre }}
    </legend>
    <div class="space-y-4" :class="titre && 'mt-2'">
      <slot />
    </div>
  </fieldset>
</template>
