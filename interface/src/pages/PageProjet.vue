<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from "vue"
import { useRoute } from "vue-router"
import { VueDraggable } from "vue-draggable-plus"
import {
  utiliserDeplacementTache,
  utiliserDetailProjet,
  utiliserOrdreColonnes,
  utiliserTaches,
} from "@/api/requetes"
import type { Colonne, FiltreTaches, Tache } from "@/api/types"
import { abonnerProjet, desabonnerProjet } from "@/tempsreel/prise"
import Bouton from "@/composants/ui/Bouton.vue"
import BoutonRetour from "@/composants/ui/BoutonRetour.vue"
import ColonneKanban from "@/composants/kanban/ColonneKanban.vue"
import FiltresTaches from "@/composants/kanban/FiltresTaches.vue"
import DialogueTache from "@/composants/kanban/DialogueTache.vue"
import DialogueColonne from "@/composants/kanban/DialogueColonne.vue"
import DialogueReferentiels from "@/composants/kanban/DialogueReferentiels.vue"

const route = useRoute()
const identifiant = computed(() => String(route.params.id))

const { data: detail } = utiliserDetailProjet(identifiant)
const filtre = ref<FiltreTaches>({})
const { data: taches } = utiliserTaches(identifiant, filtre)

const deplacement = utiliserDeplacementTache()
const ordreColonnes = utiliserOrdreColonnes()

const edition = computed(() => ["proprietaire", "administrateur", "membre"].includes(detail.value?.role ?? ""))
const gestion = computed(() => ["proprietaire", "administrateur"].includes(detail.value?.role ?? ""))

const colonnesLocales = ref<Colonne[]>([])
watch(
  () => detail.value?.colonnes,
  (valeur) => {
    colonnesLocales.value = valeur ? [...valeur] : []
  },
  { immediate: true },
)

const parColonne = computed(() => {
  const groupes: Record<string, Tache[]> = {}
  for (const colonne of colonnesLocales.value) {
    groupes[colonne.id] = []
  }
  for (const tache of taches.value ?? []) {
    groupes[tache.colonne]?.push(tache)
  }
  return groupes
})

const etiquettesParId = computed(() =>
  Object.fromEntries((detail.value?.etiquettes ?? []).map((etiquette) => [etiquette.id, etiquette])),
)
const membresParId = computed(() =>
  Object.fromEntries((detail.value?.membres ?? []).map((membre) => [membre.utilisateur, membre])),
)
const lotsParId = computed(() => Object.fromEntries((detail.value?.lots ?? []).map((lot) => [lot.id, lot])))

const totalPoints = computed(() => (taches.value ?? []).reduce((somme, tache) => somme + tache.points, 0))

watch(
  identifiant,
  (nouveau, ancien) => {
    if (ancien) desabonnerProjet(ancien)
    if (nouveau) abonnerProjet(nouveau)
  },
  { immediate: true },
)
onBeforeUnmount(() => desabonnerProjet(identifiant.value))

function surOrdreColonnes() {
  ordreColonnes.mutate({
    projet: identifiant.value,
    ordre: colonnesLocales.value.map((colonne) => colonne.id),
  })
}

function deplacer(mouvement: { id: string; colonne: string; position: number }) {
  deplacement.mutate({ ...mouvement, projet: identifiant.value })
}

const dialogueTache = ref(false)
const tacheOuverte = ref<Tache | null>(null)
const colonneInitiale = ref("")

function ouvrirCreation(colonne: string) {
  tacheOuverte.value = null
  colonneInitiale.value = colonne
  dialogueTache.value = true
}

function ouvrirTache(tache: Tache) {
  tacheOuverte.value = tache
  colonneInitiale.value = tache.colonne
  dialogueTache.value = true
}

watch(taches, (valeur) => {
  if (!tacheOuverte.value) return
  const rafraichie = valeur?.find((tache) => tache.id === tacheOuverte.value?.id)
  if (rafraichie) tacheOuverte.value = rafraichie
})

const dialogueColonne = ref(false)
const colonneOuverte = ref<Colonne | null>(null)

function ouvrirColonne(colonne: Colonne | null) {
  colonneOuverte.value = colonne
  dialogueColonne.value = true
}

const dialogueReferentiels = ref(false)
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)] flex-col">
    <div class="border-b border-neutral-200 bg-white px-4 py-4 dark:border-neutral-800 dark:bg-neutral-950">
      <div class="mx-auto max-w-full">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-3">
            <BoutonRetour v-if="detail" :vers="`/groupes/${detail.projet.groupe}`" etiquette="Groupe" />
            <span
              v-if="detail"
              class="h-3 w-3 rounded-full"
              :style="{ backgroundColor: detail.projet.couleur }"
            ></span>
            <h1 class="text-lg font-bold">{{ detail?.projet.nom }}</h1>
            <span class="text-sm text-neutral-500">{{ taches?.length ?? 0 }} tâches · {{ totalPoints }} points</span>
          </div>
          <div class="flex gap-2">
            <Bouton v-if="edition" taille="petite" variante="secondaire" @click="dialogueReferentiels = true">
              Étiquettes et lots
            </Bouton>
            <Bouton v-if="gestion" taille="petite" variante="secondaire" @click="ouvrirColonne(null)">
              Nouvelle colonne
            </Bouton>
          </div>
        </div>
        <div class="mt-3">
          <FiltresTaches
            v-model="filtre"
            :membres="detail?.membres ?? []"
            :etiquettes="detail?.etiquettes ?? []"
            :lots="detail?.lots ?? []"
          />
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-x-auto overflow-y-hidden p-4">
      <VueDraggable
        v-model="colonnesLocales"
        :animation="150"
        handle=".poignee"
        :disabled="!gestion"
        class="flex h-full items-stretch gap-4"
        @update="surOrdreColonnes"
      >
        <ColonneKanban
          v-for="colonne in colonnesLocales"
          :key="colonne.id"
          :colonne="colonne"
          :taches="parColonne[colonne.id] ?? []"
          :etiquettes="etiquettesParId"
          :membres="membresParId"
          :lots="lotsParId"
          :edition="edition"
          :gestion="gestion"
          @deplacer="deplacer"
          @ouvrir="ouvrirTache"
          @creer="ouvrirCreation(colonne.id)"
          @modifier="ouvrirColonne(colonne)"
        />
      </VueDraggable>
    </div>

    <DialogueTache
      :ouvert="dialogueTache"
      :projet="identifiant"
      :colonnes="detail?.colonnes ?? []"
      :membres="detail?.membres ?? []"
      :etiquettes="detail?.etiquettes ?? []"
      :lots="detail?.lots ?? []"
      :tache="tacheOuverte"
      :colonne-initiale="colonneInitiale"
      :edition="edition"
      @fermer="dialogueTache = false"
    />

    <DialogueColonne
      :ouvert="dialogueColonne"
      :projet="identifiant"
      :colonne="colonneOuverte"
      @fermer="dialogueColonne = false"
    />

    <DialogueReferentiels
      :ouvert="dialogueReferentiels"
      :projet="identifiant"
      :etiquettes="detail?.etiquettes ?? []"
      :lots="detail?.lots ?? []"
      @fermer="dialogueReferentiels = false"
    />
  </div>
</template>
