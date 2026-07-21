<script setup lang="ts">
defineProps<{
  ouvert: boolean
  titre: string
  large?: boolean
  geant?: boolean
}>()

const emissions = defineEmits<{ fermer: [] }>()
</script>

<template>
  <Teleport to="body">
    <div
      v-if="ouvert"
      class="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-neutral-950/50 p-3 pt-8 sm:p-4 sm:pt-16"
      @click.self="emissions('fermer')"
    >
      <div
        class="w-full rounded-xl border border-neutral-200 bg-white shadow-xl dark:border-neutral-800 dark:bg-neutral-900"
        :class="geant ? 'max-w-4xl' : large ? 'max-w-2xl' : 'max-w-md'"
      >
        <div class="flex items-center justify-between border-b border-neutral-200 px-4 py-3.5 dark:border-neutral-800 sm:px-6 sm:py-4">
          <h2 class="text-base font-semibold tracking-tight">{{ titre }}</h2>
          <button
            class="rounded-md p-1.5 text-neutral-500 transition-colors hover:bg-neutral-100 hover:text-neutral-900 dark:hover:bg-neutral-800 dark:hover:text-neutral-100"
            @click="emissions('fermer')"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="max-h-[75vh] overflow-y-auto px-4 py-4 sm:max-h-[70vh] sm:px-6 sm:py-5">
          <slot />
        </div>
        <div
          v-if="$slots.pied"
          class="flex flex-wrap justify-end gap-2 rounded-b-xl border-t border-neutral-200 bg-neutral-50 px-4 py-3.5 dark:border-neutral-800 dark:bg-neutral-950/40 sm:px-6 sm:py-4"
        >
          <slot name="pied" />
        </div>
      </div>
    </div>
  </Teleport>
</template>
