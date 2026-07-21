<script setup lang="ts">
import { ref, watch } from "vue"
import { VueDraggable } from "vue-draggable-plus"
import type { Colonne, Etiquette, Lot, Membre, Tache } from "@/api/types"
import CarteTache, { type Densite } from "@/composants/kanban/CarteTache.vue"

const proprietes = defineProps<{
  colonne: Colonne
  taches: Tache[]
  etiquettes: Record<string, Etiquette>
  membres: Record<string, Membre>
  lots: Record<string, Lot>
  edition: boolean
  gestion: boolean
  densite: Densite
}>()

const emissions = defineEmits<{
  deplacer: [{ id: string; colonne: string; position: number }]
  ouvrir: [Tache]
  supprimer: [Tache]
  creer: []
  modifier: []
  menu: [MouseEvent]
}>()

const locales = ref<Tache[]>([])

watch(
  () => proprietes.taches,
  (valeur) => {
    locales.value = [...valeur]
  },
  { immediate: true },
)

function surChangement(evenement: { newIndex?: number }) {
  const position = evenement.newIndex ?? 0
  const tache = locales.value[position]
  if (!tache) return
  emissions("deplacer", { id: tache.id, colonne: proprietes.colonne.id, position })
}
</script>

<template>
  <section
    class="flex w-[17rem] shrink-0 flex-col rounded-xl border border-neutral-200 bg-neutral-50 dark:border-neutral-800 dark:bg-neutral-900 sm:w-72"
    @contextmenu.prevent="emissions('menu', $event)"
  >
    <header class="flex items-center justify-between gap-2 px-3 py-2.5">
      <div class="poignee flex min-w-0 items-center gap-2" :class="gestion && 'cursor-grab'">
        <span class="h-2.5 w-2.5 shrink-0 rounded-full" :style="{ backgroundColor: colonne.couleur }"></span>
        <h3 class="truncate text-sm font-semibold">{{ colonne.nom }}</h3>
        <span
          class="shrink-0 whitespace-nowrap rounded-full bg-neutral-200 px-1.5 text-xs text-neutral-600 dark:bg-neutral-800 dark:text-neutral-400"
          :class="colonne.limite !== null && taches.length > colonne.limite && '!bg-red-100 !text-red-700 dark:!bg-red-950 dark:!text-red-400'"
        >
          {{ taches.length }}<template v-if="colonne.limite !== null">/{{ colonne.limite }}</template>
        </span>
      </div>
      <button
        v-if="gestion"
        class="shrink-0 text-xs text-neutral-500 hover:text-neutral-900 hover:underline dark:hover:text-neutral-100"
        @click="emissions('modifier')"
      >
        Modifier
      </button>
    </header>

    <VueDraggable
      v-model="locales"
      group="taches"
      :animation="150"
      ghost-class="fantome"
      :disabled="!edition"
      class="flex min-h-24 flex-1 flex-col gap-2 overflow-y-auto px-3 pb-3"
      @add="surChangement"
      @update="surChangement"
    >
      <CarteTache
        v-for="tache in locales"
        :key="tache.id"
        :tache="tache"
        :etiquettes="etiquettes"
        :membres="membres"
        :lots="lots"
        :densite="densite"
        :edition="edition"
        @ouvrir="emissions('ouvrir', tache)"
        @supprimer="emissions('supprimer', tache)"
      />
    </VueDraggable>

    <footer v-if="edition" class="px-3 pb-3">
      <button
        class="w-full rounded-lg border border-dashed border-neutral-300 py-2 text-sm text-neutral-500 hover:border-neutral-400 hover:text-neutral-900 dark:border-neutral-700 dark:hover:border-neutral-500 dark:hover:text-neutral-100"
        @click="emissions('creer')"
      >
        + Ajouter une tâche
      </button>
    </footer>
  </section>
</template>
