<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from "vue"
import { useRoute } from "vue-router"
import { useLocalStorage } from "@vueuse/core"
import { VueDraggable } from "vue-draggable-plus"
import {
  utiliserDeplacementTache,
  utiliserDetailProjet,
  utiliserOrdreColonnes,
  utiliserSuppressionColonne,
  utiliserTaches,
} from "@/api/requetes"
import type { Colonne, FiltreTaches, Tache } from "@/api/types"
import { abonnerProjet, connectes, desabonnerProjet } from "@/tempsreel/prise"
import Avatar from "@/composants/ui/Avatar.vue"
import Bouton from "@/composants/ui/Bouton.vue"
import BoutonRetour from "@/composants/ui/BoutonRetour.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import MenuContextuel, { type ElementMenu } from "@/composants/ui/MenuContextuel.vue"
import ColonneKanban from "@/composants/kanban/ColonneKanban.vue"
import { type Densite } from "@/composants/kanban/CarteTache.vue"
import FiltresTaches from "@/composants/kanban/FiltresTaches.vue"
import DialogueTache from "@/composants/kanban/DialogueTache.vue"
import DialogueColonne from "@/composants/kanban/DialogueColonne.vue"
import DialogueReferentiels from "@/composants/kanban/DialogueReferentiels.vue"
import DialogueAccesIA from "@/composants/projet/DialogueAccesIA.vue"
import DialogueCorbeille from "@/composants/projet/DialogueCorbeille.vue"
import DialogueParametresProjet from "@/composants/projet/DialogueParametresProjet.vue"

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
const dialogueAccesIA = ref(false)
const dialogueCorbeille = ref(false)
const dialogueParametres = ref(false)

const densite = useLocalStorage<Densite>("historykanban.densite", "defaut")
const densites: { valeur: Densite; libelle: string }[] = [
  { valeur: "compacte", libelle: "Compacte" },
  { valeur: "defaut", libelle: "Défaut" },
  { valeur: "detaillee", libelle: "Détaillée" },
]

const suppressionColonne = utiliserSuppressionColonne()
const colonneASupprimer = ref<Colonne | null>(null)

const menu = ref<{ ouvert: boolean; x: number; y: number; colonne: Colonne | null }>({
  ouvert: false,
  x: 0,
  y: 0,
  colonne: null,
})

function ouvrirMenu(colonne: Colonne, evenement: MouseEvent) {
  if (!edition.value) return
  menu.value = { ouvert: true, x: evenement.clientX, y: evenement.clientY, colonne }
}

const elementsMenu = computed<ElementMenu[]>(() => {
  if (!menu.value.colonne) return []
  const indice = colonnesLocales.value.findIndex((colonne) => colonne.id === menu.value.colonne?.id)
  const elements: ElementMenu[] = [{ id: "tache", libelle: "Nouvelle tâche" }]
  if (gestion.value) {
    elements.push({ id: "modifier", libelle: "Modifier la colonne" })
    if (indice > 0) elements.push({ id: "gauche", libelle: "Déplacer à gauche" })
    if (indice < colonnesLocales.value.length - 1) elements.push({ id: "droite", libelle: "Déplacer à droite" })
    elements.push({ id: "supprimer", libelle: "Supprimer la colonne", danger: true })
  }
  return elements
})

function decalerColonne(direction: number) {
  const colonne = menu.value.colonne
  if (!colonne) return
  const indice = colonnesLocales.value.findIndex((element) => element.id === colonne.id)
  const cible = indice + direction
  if (indice < 0 || cible < 0 || cible >= colonnesLocales.value.length) return
  const copie = [...colonnesLocales.value]
  ;[copie[indice], copie[cible]] = [copie[cible], copie[indice]]
  colonnesLocales.value = copie
  surOrdreColonnes()
}

function choisirMenu(action: string) {
  const colonne = menu.value.colonne
  if (!colonne) return
  switch (action) {
    case "tache":
      ouvrirCreation(colonne.id)
      break
    case "modifier":
      ouvrirColonne(colonne)
      break
    case "gauche":
      decalerColonne(-1)
      break
    case "droite":
      decalerColonne(1)
      break
    case "supprimer":
      colonneASupprimer.value = colonne
      break
  }
}

function supprimerColonne() {
  if (!colonneASupprimer.value) return
  suppressionColonne.mutate(
    { id: colonneASupprimer.value.id, projet: identifiant.value },
    { onSuccess: () => (colonneASupprimer.value = null) },
  )
}
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)] flex-col">
    <div class="border-b border-neutral-200 bg-white px-3 py-3 dark:border-neutral-800 dark:bg-neutral-950 sm:px-4 sm:py-4">
      <div class="mx-auto max-w-full">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex min-w-0 flex-wrap items-center gap-2 sm:gap-3">
            <BoutonRetour v-if="detail" :vers="`/groupes/${detail.projet.groupe}`" etiquette="Groupe" />
            <span
              v-if="detail"
              class="h-3 w-3 shrink-0 rounded-full"
              :style="{ backgroundColor: detail.projet.couleur }"
            ></span>
            <h1 class="truncate text-base font-bold sm:text-lg">{{ detail?.projet.nom }}</h1>
            <span class="text-xs text-neutral-500 sm:text-sm">{{ taches?.length ?? 0 }} tâches · {{ totalPoints }} points</span>
          </div>
          <div class="flex flex-wrap items-center gap-2 sm:gap-3">
            <div class="hidden -space-x-1.5 sm:flex">
              <span
                v-for="membre in (detail?.membres ?? []).slice(0, 8)"
                :key="membre.utilisateur"
                class="relative inline-flex"
              >
                <Avatar :nom="membre.nom" :prenom="membre.prenom" />
                <span
                  class="absolute -bottom-0.5 -right-0.5 h-2.5 w-2.5 rounded-full border-2 border-white dark:border-neutral-950"
                  :class="connectes.includes(membre.utilisateur) ? 'bg-green-500' : 'bg-neutral-400'"
                  :title="connectes.includes(membre.utilisateur) ? 'En ligne' : 'Hors ligne'"
                ></span>
              </span>
            </div>
            <div class="flex rounded-lg border border-neutral-300 p-0.5 dark:border-neutral-700">
              <button
                v-for="option in densites"
                :key="option.valeur"
                class="rounded-md px-2.5 py-1 text-xs font-medium transition-colors"
                :class="
                  densite === option.valeur
                    ? 'bg-neutral-900 text-white dark:bg-neutral-100 dark:text-neutral-900'
                    : 'text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100'
                "
                @click="densite = option.valeur"
              >
                {{ option.libelle }}
              </button>
            </div>
            <Bouton v-if="edition" taille="petite" variante="secondaire" @click="dialogueAccesIA = true">
              Accès IA
            </Bouton>
            <Bouton v-if="edition" taille="petite" variante="secondaire" @click="dialogueCorbeille = true">
              Corbeille
            </Bouton>
            <Bouton v-if="edition" taille="petite" variante="secondaire" @click="dialogueReferentiels = true">
              Étiquettes et lots
            </Bouton>
            <Bouton v-if="gestion" taille="petite" variante="secondaire" @click="ouvrirColonne(null)">
              Nouvelle colonne
            </Bouton>
            <Bouton v-if="gestion" taille="petite" variante="secondaire" @click="dialogueParametres = true">
              Paramètres
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

    <div class="flex-1 overflow-x-auto overflow-y-hidden p-3 sm:p-4">
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
          :densite="densite"
          @deplacer="deplacer"
          @ouvrir="ouvrirTache"
          @creer="ouvrirCreation(colonne.id)"
          @modifier="ouvrirColonne(colonne)"
          @menu="ouvrirMenu(colonne, $event)"
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
      :depot="detail?.projet.depot ?? ''"
      @fermer="dialogueTache = false"
    />

    <DialogueParametresProjet
      :ouvert="dialogueParametres"
      :projet="detail?.projet ?? null"
      @fermer="dialogueParametres = false"
    />

    <DialogueColonne
      :ouvert="dialogueColonne"
      :projet="identifiant"
      :colonne="colonneOuverte"
      @fermer="dialogueColonne = false"
    />

    <DialogueAccesIA :ouvert="dialogueAccesIA" :projet="identifiant" @fermer="dialogueAccesIA = false" />
    <DialogueCorbeille :ouvert="dialogueCorbeille" :projet="identifiant" @fermer="dialogueCorbeille = false" />

    <DialogueReferentiels
      :ouvert="dialogueReferentiels"
      :projet="identifiant"
      :etiquettes="detail?.etiquettes ?? []"
      :lots="detail?.lots ?? []"
      @fermer="dialogueReferentiels = false"
    />

    <MenuContextuel
      :ouvert="menu.ouvert"
      :x="menu.x"
      :y="menu.y"
      :elements="elementsMenu"
      @choisir="choisirMenu"
      @fermer="menu.ouvert = false"
    />

    <Dialogue
      :ouvert="colonneASupprimer !== null"
      titre="Supprimer la colonne"
      @fermer="colonneASupprimer = null"
    >
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        La colonne « {{ colonneASupprimer?.nom }} » et toutes ses tâches seront supprimées définitivement.
      </p>
      <template #pied>
        <Bouton variante="secondaire" @click="colonneASupprimer = null">Annuler</Bouton>
        <Bouton variante="danger" :desactive="suppressionColonne.isPending.value" @click="supprimerColonne">
          Supprimer
        </Bouton>
      </template>
    </Dialogue>
  </div>
</template>
