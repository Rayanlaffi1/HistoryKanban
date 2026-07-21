<script setup lang="ts">
import { utiliserMagasinNotifications } from "@/magasins/notifications"

const magasin = utiliserMagasinNotifications()
</script>

<template>
  <div class="pointer-events-none fixed bottom-4 right-4 z-50 flex w-80 flex-col gap-2">
    <TransitionGroup
      enter-active-class="transition duration-200"
      enter-from-class="translate-y-2 opacity-0"
      leave-active-class="transition duration-200"
      leave-to-class="opacity-0"
    >
      <div
        v-for="message in magasin.messages"
        :key="message.id"
        class="pointer-events-auto rounded-xl border border-neutral-200 bg-white p-4 shadow-lg dark:border-neutral-700 dark:bg-neutral-900"
      >
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <p class="text-sm font-semibold">{{ message.titre }}</p>
            <p v-if="message.detail" class="mt-0.5 text-sm text-neutral-600 dark:text-neutral-400">
              {{ message.detail }}
            </p>
          </div>
          <button
            class="shrink-0 text-neutral-400 hover:text-neutral-900 dark:hover:text-neutral-100"
            @click="magasin.fermer(message.id)"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6 6 18M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>
