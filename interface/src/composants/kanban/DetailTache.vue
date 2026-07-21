<script setup lang="ts">
import { computed } from "vue"
import type { Colonne, Etiquette, Lot, Membre, Tache } from "@/api/types"
import { couleursUrgences, libellesUrgences } from "@/api/types"
import { formaterDateHeure } from "@/utilitaires/dates"
import { assainir } from "@/utilitaires/html"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"
import PieceJointe from "@/composants/kanban/PieceJointe.vue"

const proprietes = defineProps<{
  tache: Tache
  colonnes: Colonne[]
  membres: Membre[]
  etiquettes: Etiquette[]
  lots: Lot[]
  depot?: string
}>()

const nomColonne = computed(
  () => proprietes.colonnes.find((colonne) => colonne.id === proprietes.tache.colonne)?.nom ?? "",
)
const lotCourant = computed(() => proprietes.lots.find((lot) => lot.id === proprietes.tache.lot))
const etiquettesCourantes = computed(() =>
  proprietes.etiquettes.filter((etiquette) => proprietes.tache.etiquettes.includes(etiquette.id)),
)
const membresAffectes = computed(() =>
  proprietes.membres.filter((membre) => proprietes.tache.affectations.includes(membre.utilisateur)),
)

const lienCommit = computed(() => {
  if (!proprietes.depot || !proprietes.tache.commit) return ""
  return `${proprietes.depot.replace(/\/+$/, "")}/commit/${proprietes.tache.commit}`
})
</script>

<template>
  <div class="space-y-5">
    <div class="flex items-start justify-between gap-3">
      <h3 class="min-w-0 break-words text-lg font-semibold leading-snug">{{ tache.titre }}</h3>
      <div class="flex shrink-0 items-center gap-1.5">
        <span
          v-if="tache.urgence !== 'normale'"
          class="rounded-full px-2 py-0.5 text-xs font-medium"
          :class="couleursUrgences[tache.urgence]"
        >
          {{ libellesUrgences[tache.urgence] }}
        </span>
        <Badge v-if="tache.points > 0">{{ tache.points }} pt{{ tache.points > 1 ? "s" : "" }}</Badge>
      </div>
    </div>

    <div v-if="tache.description" class="texteriche text-sm" v-html="assainir(tache.description)"></div>
    <p v-else class="text-sm italic text-neutral-500">Aucune description.</p>

    <dl class="grid grid-cols-2 gap-x-4 gap-y-2 rounded-lg bg-neutral-100 p-4 text-sm dark:bg-neutral-800/60">
      <dt class="text-neutral-500">Colonne</dt>
      <dd class="font-medium">{{ nomColonne || "—" }}</dd>
      <dt class="text-neutral-500">Lot</dt>
      <dd class="font-medium">
        <template v-if="lotCourant">
          {{ lotCourant.nom }}
          <span v-if="lotCourant.echeance" class="text-neutral-500">
            — {{ formaterDateHeure(lotCourant.echeance) }}
          </span>
        </template>
        <template v-else>—</template>
      </dd>
      <dt class="text-neutral-500">Urgence</dt>
      <dd class="font-medium">{{ libellesUrgences[tache.urgence] }}</dd>
      <dt class="text-neutral-500">Échéance</dt>
      <dd class="font-medium">{{ tache.echeance ? formaterDateHeure(tache.echeance) : "—" }}</dd>
      <dt class="text-neutral-500">Créée le</dt>
      <dd class="font-medium">{{ formaterDateHeure(tache.creation) }}</dd>
      <template v-if="tache.commit">
        <dt class="text-neutral-500">Commit</dt>
        <dd class="font-mono text-xs font-medium">
          <a
            v-if="lienCommit"
            :href="lienCommit"
            target="_blank"
            rel="noopener"
            class="underline decoration-neutral-400 underline-offset-2 hover:decoration-neutral-900 dark:hover:decoration-neutral-100"
          >
            {{ tache.commit }}
          </a>
          <template v-else>{{ tache.commit }}</template>
        </dd>
      </template>
    </dl>

    <div v-if="etiquettesCourantes.length">
      <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Étiquettes</span>
      <div class="mt-2 flex flex-wrap gap-1.5">
        <Badge v-for="etiquette in etiquettesCourantes" :key="etiquette.id" :couleur="etiquette.couleur">
          {{ etiquette.nom }}
        </Badge>
      </div>
    </div>

    <div v-if="membresAffectes.length">
      <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Affectations</span>
      <div class="mt-2 flex flex-wrap gap-2">
        <span
          v-for="membre in membresAffectes"
          :key="membre.utilisateur"
          class="flex items-center gap-2 rounded-full border border-neutral-300 px-2 py-1 text-xs text-neutral-700 dark:border-neutral-700 dark:text-neutral-300"
        >
          <Avatar :nom="membre.nom" :prenom="membre.prenom" petite />
          {{ membre.prenom }} {{ membre.nom }}
        </span>
      </div>
    </div>

    <div v-if="tache.images.length">
      <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Pièces jointes</span>
      <div class="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-3">
        <PieceJointe
          v-for="image in tache.images"
          :key="image.id"
          :nom="image.nom"
          :taille="image.taille"
          :typecontenu="image.typecontenu"
          :url="image.url"
        />
      </div>
    </div>
  </div>
</template>
