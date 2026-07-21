<script setup lang="ts">
import { computed, ref } from "vue"
import { useRouter } from "vue-router"
import { onClickOutside, refDebounced } from "@vueuse/core"
import { utiliserRecherche } from "@/api/requetes"
import { texteBrut } from "@/utilitaires/html"
import type { Resultat } from "@/api/types"

const routeur = useRouter()

const saisie = ref("")
const differee = refDebounced(saisie, 250)
const ouvert = ref(false)
const zone = ref<HTMLElement | null>(null)

onClickOutside(zone, () => (ouvert.value = false))

const { data: resultats, isFetching: chargement } = utiliserRecherche(differee)

const liste = computed(() => resultats.value ?? [])
const assezLong = computed(() => saisie.value.trim().length >= 2)

const libellesOrigines: Record<string, string> = {
  titre: "titre",
  description: "description",
  commentaire: "commentaire",
}

function extrait(resultat: Resultat) {
  const brut = texteBrut(resultat.extrait).trim()
  if (resultat.origine === "titre" || brut.length <= 120) return brut
  const position = brut.toLowerCase().indexOf(differee.value.trim().toLowerCase())
  if (position < 0) return `${brut.slice(0, 120)}…`
  const debut = Math.max(0, position - 40)
  return `${debut > 0 ? "…" : ""}${brut.slice(debut, debut + 120)}…`
}

function ouvrir(resultat: Resultat) {
  ouvert.value = false
  saisie.value = ""
  routeur.push({ name: "projet", params: { id: resultat.projet }, query: { tache: resultat.tache } })
}
</script>

<template>
  <div ref="zone" class="relative">
    <input
      v-model="saisie"
      type="search"
      placeholder="Rechercher partout…"
      class="h-9 w-36 rounded-lg border border-neutral-300 bg-white px-2.5 text-sm text-neutral-900 shadow-sm transition placeholder:text-neutral-400 focus:w-56 focus:border-neutral-500 focus:outline-none focus:ring-2 focus:ring-neutral-900/10 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:placeholder:text-neutral-500 dark:focus:border-neutral-400 dark:focus:ring-neutral-100/10 sm:w-48"
      @focus="ouvert = true"
      @keydown.escape="ouvert = false"
    />

    <div
      v-if="ouvert && assezLong"
      class="absolute right-0 top-11 z-50 max-h-[70vh] w-[22rem] overflow-y-auto rounded-xl border border-neutral-200 bg-white shadow-lg dark:border-neutral-700 dark:bg-neutral-900 sm:w-[26rem]"
    >
      <p v-if="chargement && !liste.length" class="px-4 py-3 text-sm text-neutral-500">Recherche…</p>
      <p v-else-if="!liste.length" class="px-4 py-3 text-sm text-neutral-500">Aucun résultat.</p>
      <ul v-else class="divide-y divide-neutral-100 dark:divide-neutral-800">
        <li v-for="resultat in liste" :key="resultat.tache">
          <button
            class="w-full px-4 py-2.5 text-left transition-colors hover:bg-neutral-50 dark:hover:bg-neutral-800"
            @click="ouvrir(resultat)"
          >
            <div class="flex items-center gap-2">
              <span class="h-2 w-2 shrink-0 rounded-full" :style="{ backgroundColor: resultat.projetcouleur }"></span>
              <span class="truncate text-xs text-neutral-500">
                {{ resultat.groupenom }} · {{ resultat.projetnom }} · {{ resultat.colonne }}
              </span>
            </div>
            <p class="mt-0.5 break-words text-sm font-medium">{{ resultat.titre }}</p>
            <p v-if="resultat.origine !== 'titre'" class="mt-0.5 line-clamp-2 text-xs text-neutral-500">
              <span class="text-neutral-400">{{ libellesOrigines[resultat.origine] }} :</span>
              {{ extrait(resultat) }}
            </p>
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>
