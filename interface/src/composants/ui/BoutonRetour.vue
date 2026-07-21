<script setup lang="ts">
defineProps<{
  vers: string
  etiquette: string
}>()
</script>

<template>
  <RouterLink
    :to="vers"
    class="group inline-flex h-10 items-center gap-2.5 rounded-lg border border-neutral-300 bg-white pl-3 pr-4 text-sm font-semibold text-neutral-800 shadow-sm transition-colors hover:border-neutral-500 hover:bg-neutral-100 dark:border-neutral-600 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:border-neutral-400 dark:hover:bg-neutral-700"
  >
    <span
      class="flex h-6 w-6 items-center justify-center rounded-md bg-neutral-900 text-white transition-transform group-hover:-translate-x-0.5 dark:bg-neutral-100 dark:text-neutral-900"
    >
      <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M19 12H5" />
        <path d="m12 19-7-7 7-7" />
      </svg>
    </span>
    {{ etiquette }}
  </RouterLink>
</template>
