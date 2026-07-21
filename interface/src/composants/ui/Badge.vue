<script setup lang="ts">
defineProps<{
  couleur?: string
}>()
</script>

<template>
  <span
    class="inline-flex max-w-full shrink-0 items-center gap-1 overflow-hidden whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium"
    :style="
      couleur
        ? { backgroundColor: couleur + '22', color: couleur, border: `1px solid ${couleur}55` }
        : undefined
    "
    :class="!couleur && 'bg-neutral-200 text-neutral-700 dark:bg-neutral-800 dark:text-neutral-300'"
  >
    <slot />
  </span>
</template>
