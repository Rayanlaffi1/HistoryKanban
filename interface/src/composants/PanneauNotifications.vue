<script setup lang="ts">
import { utiliserNotifications, utiliserLectureNotification, utiliserLectureTotale } from "@/api/requetes"
import { typesNotifications } from "@/api/types"
import { depuis } from "@/utilitaires/dates"

defineEmits<{ fermer: [] }>()

const { data: notifications } = utiliserNotifications()
const lecture = utiliserLectureNotification()
const lectureTotale = utiliserLectureTotale()
</script>

<template>
  <div
    class="absolute right-0 top-full mt-2 w-80 overflow-hidden rounded-xl border border-neutral-200 bg-white shadow-lg dark:border-neutral-800 dark:bg-neutral-900"
  >
    <div class="flex items-center justify-between border-b border-neutral-200 px-4 py-3 dark:border-neutral-800">
      <span class="text-sm font-semibold">Notifications</span>
      <button
        class="text-xs text-neutral-500 hover:text-neutral-900 hover:underline dark:hover:text-neutral-100"
        @click="lectureTotale.mutate(undefined)"
      >
        Tout marquer lu
      </button>
    </div>
    <div class="max-h-96 overflow-y-auto">
      <p v-if="!notifications?.length" class="px-4 py-6 text-center text-sm text-neutral-500">
        Aucune notification
      </p>
      <button
        v-for="notification in notifications"
        :key="notification.id"
        class="block w-full border-b border-neutral-100 px-4 py-3 text-left last:border-0 hover:bg-neutral-50 dark:border-neutral-800/60 dark:hover:bg-neutral-800/50"
        :class="!notification.lue && 'bg-neutral-100/70 dark:bg-neutral-800/70'"
        @click="lecture.mutate(notification.id)"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="text-xs font-semibold">
            {{ typesNotifications[notification.type] ?? notification.type }}
          </span>
          <span class="shrink-0 text-[11px] text-neutral-400">{{ depuis(notification.creation) }}</span>
        </div>
        <p class="mt-0.5 truncate text-sm text-neutral-600 dark:text-neutral-400">
          {{ notification.contenu.titre }}
          <span v-if="notification.contenu.acteurnom" class="text-neutral-400 dark:text-neutral-500">
            — {{ notification.contenu.acteurnom }}
          </span>
        </p>
      </button>
    </div>
  </div>
</template>
