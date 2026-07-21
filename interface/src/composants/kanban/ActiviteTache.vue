<script setup lang="ts">
import { computed } from "vue"
import { utiliserActivites } from "@/api/requetes"
import { libellesActivites } from "@/api/types"
import { depuis, formaterDateHeure } from "@/utilitaires/dates"
import Avatar from "@/composants/ui/Avatar.vue"

const proprietes = defineProps<{
  tache: string
}>()

const identifiant = computed(() => proprietes.tache)
const { data: activites, isLoading: chargement } = utiliserActivites(identifiant)
</script>

<template>
  <div>
    <p v-if="chargement" class="py-6 text-center text-sm text-neutral-500">Chargement…</p>
    <ol v-else-if="activites?.length" class="relative space-y-4 border-l border-neutral-200 pl-5 dark:border-neutral-800">
      <li v-for="activite in activites" :key="activite.id" class="relative">
        <span class="absolute -left-[27px] top-1 h-2.5 w-2.5 rounded-full bg-neutral-400 ring-4 ring-white dark:bg-neutral-500 dark:ring-neutral-900"></span>
        <div class="flex items-start gap-2.5">
          <Avatar :nom="activite.nom || '?'" :prenom="activite.prenom" petite />
          <div class="min-w-0">
            <p class="text-sm">
              <span class="font-medium">{{ activite.prenom }} {{ activite.nom }}</span>
              {{ libellesActivites[activite.type] ?? activite.type }}
              <span v-if="activite.detail" class="text-neutral-500">{{ activite.detail }}</span>
            </p>
            <p class="text-xs text-neutral-400" :title="formaterDateHeure(activite.creation)">
              {{ depuis(activite.creation) }}
            </p>
          </div>
        </div>
      </li>
    </ol>
    <p v-else class="rounded-lg border border-dashed border-neutral-300 p-6 text-center text-sm text-neutral-500 dark:border-neutral-700">
      Aucune activité enregistrée sur cette tâche.
    </p>
  </div>
</template>
