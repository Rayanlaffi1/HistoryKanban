<script setup lang="ts">
import { computed, ref } from "vue"
import type { Tache } from "@/api/types"
import {
  utiliserCreationSousTache,
  utiliserMutationSousTache,
  utiliserSuppressionSousTache,
} from "@/api/requetes"
import Bouton from "@/composants/ui/Bouton.vue"

const proprietes = defineProps<{ tache: Tache; projet: string; edition: boolean }>()

const creation = utiliserCreationSousTache()
const mutation = utiliserMutationSousTache()
const suppression = utiliserSuppressionSousTache()

const nouvelle = ref("")

const liste = computed(() => proprietes.tache.soustaches ?? [])
const faites = computed(() => liste.value.filter((sousTache) => sousTache.faite).length)
const proportion = computed(() => (liste.value.length ? (faites.value / liste.value.length) * 100 : 0))

function ajouter() {
  const libelle = nouvelle.value.trim()
  if (!libelle) return
  creation.mutate({ tache: proprietes.tache.id, projet: proprietes.projet, libelle })
  nouvelle.value = ""
}

function basculer(identifiant: string, faite: boolean) {
  mutation.mutate({ id: identifiant, projet: proprietes.projet, faite: !faite })
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between gap-3">
      <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Sous-tâches</span>
      <span v-if="liste.length" class="text-xs text-neutral-500">{{ faites }} / {{ liste.length }}</span>
    </div>

    <div v-if="liste.length" class="mt-2 h-1.5 w-full overflow-hidden rounded-full bg-neutral-200 dark:bg-neutral-700">
      <div
        class="h-full rounded-full bg-neutral-900 transition-all dark:bg-neutral-100"
        :style="{ width: `${proportion}%` }"
      ></div>
    </div>

    <ul v-if="liste.length" class="mt-2 space-y-1">
      <li v-for="sousTache in liste" :key="sousTache.id" class="group flex items-center gap-2">
        <input
          type="checkbox"
          class="h-4 w-4 shrink-0 rounded border-neutral-300 text-neutral-900 focus:ring-neutral-500 dark:border-neutral-600 dark:bg-neutral-800"
          :checked="sousTache.faite"
          :disabled="!edition"
          @change="basculer(sousTache.id, sousTache.faite)"
        />
        <span
          class="min-w-0 flex-1 break-words text-sm"
          :class="sousTache.faite ? 'text-neutral-400 line-through' : ''"
        >
          {{ sousTache.libelle }}
        </span>
        <button
          v-if="edition"
          class="hidden shrink-0 text-xs text-neutral-400 hover:text-red-700 group-hover:block dark:hover:text-red-500"
          @click="suppression.mutate({ id: sousTache.id, projet })"
        >
          Retirer
        </button>
      </li>
    </ul>
    <p v-else class="mt-2 text-xs text-neutral-500">Aucune sous-tâche.</p>

    <div v-if="edition" class="mt-3 flex gap-2">
      <input
        v-model="nouvelle"
        type="text"
        maxlength="200"
        placeholder="Ajouter une sous-tâche…"
        class="h-9 min-w-0 flex-1 rounded-lg border border-neutral-300 bg-white px-2.5 text-sm text-neutral-900 shadow-sm transition placeholder:text-neutral-400 focus:border-neutral-500 focus:outline-none focus:ring-2 focus:ring-neutral-900/10 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder:text-neutral-500"
        @keydown.enter.prevent="ajouter"
      />
      <Bouton taille="petite" variante="secondaire" :desactive="creation.isPending.value" @click="ajouter">
        Ajouter
      </Bouton>
    </div>
  </div>
</template>
