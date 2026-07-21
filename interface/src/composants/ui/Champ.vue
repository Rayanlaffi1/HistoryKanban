<script setup lang="ts">
defineProps<{
  etiquette?: string
  type?: string
  obligatoire?: boolean
  indication?: string
  erreur?: string
}>()

const modele = defineModel<string>({ default: "" })
</script>

<template>
  <label class="block space-y-1.5">
    <span v-if="etiquette" class="flex items-baseline gap-1 text-sm font-medium text-neutral-700 dark:text-neutral-300">
      {{ etiquette }}
      <span v-if="obligatoire" class="text-red-600 dark:text-red-500">*</span>
    </span>
    <input
      v-model="modele"
      :type="type ?? 'text'"
      :placeholder="indication"
      class="h-10 w-full rounded-lg border bg-white px-3 text-sm text-neutral-900 shadow-sm transition placeholder:text-neutral-400 focus:outline-none focus:ring-2 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder:text-neutral-500"
      :class="
        erreur
          ? 'border-red-500 focus:border-red-600 focus:ring-red-600/15 dark:border-red-700'
          : 'border-neutral-300 focus:border-neutral-500 focus:ring-neutral-900/10 dark:border-neutral-700 dark:focus:border-neutral-400 dark:focus:ring-neutral-100/10'
      "
    />
    <span v-if="erreur" class="block text-xs font-medium text-red-700 dark:text-red-500">{{ erreur }}</span>
  </label>
</template>
