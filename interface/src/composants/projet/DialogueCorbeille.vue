<script setup lang="ts">
import { computed, ref, toRef } from "vue"
import { utiliserCorbeille, utiliserPurgeTache, utiliserRestaurationTache } from "@/api/requetes"
import { formaterDateHeure } from "@/utilitaires/dates"
import Bouton from "@/composants/ui/Bouton.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"

const proprietes = defineProps<{ ouvert: boolean; projet: string }>()
const emissions = defineEmits<{ fermer: [] }>()

const ouvert = toRef(proprietes, "ouvert")
const projet = toRef(proprietes, "projet")

const { data: taches, isLoading: chargement } = utiliserCorbeille(projet, ouvert)
const restauration = utiliserRestaurationTache()
const purge = utiliserPurgeTache()

const confirmation = ref("")

function joursRestants(suppression: string | null) {
  if (!suppression) return 0
  const ecoules = (Date.now() - new Date(suppression).getTime()) / 86400000
  return Math.max(0, Math.ceil(30 - ecoules))
}

const liste = computed(() => taches.value ?? [])
</script>

<template>
  <Dialogue :ouvert="ouvert" titre="Corbeille" geant @fermer="emissions('fermer')">
    <p class="text-sm text-neutral-600 dark:text-neutral-400">
      Les tâches supprimées restent ici trente jours, puis sont effacées définitivement avec leurs pièces
      jointes. Elles ne comptent plus dans le tableau ni dans les statistiques.
    </p>

    <p v-if="chargement" class="mt-4 text-sm text-neutral-500">Chargement…</p>
    <p v-else-if="!liste.length" class="mt-4 text-sm italic text-neutral-500">La corbeille est vide.</p>

    <ul v-else class="mt-4 space-y-2">
      <li
        v-for="tache in liste"
        :key="tache.id"
        class="rounded-lg border border-neutral-200 p-3 dark:border-neutral-700"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <p class="break-words text-sm font-medium">{{ tache.titre }}</p>
            <p class="mt-0.5 text-xs text-neutral-500">
              Supprimée le {{ formaterDateHeure(tache.suppression ?? "") }} · effacement définitif dans
              {{ joursRestants(tache.suppression) }} jour{{ joursRestants(tache.suppression) > 1 ? "s" : "" }}
            </p>
          </div>
          <div class="flex shrink-0 gap-2">
            <Bouton
              taille="petite"
              variante="secondaire"
              :desactive="restauration.isPending.value"
              @click="restauration.mutate({ id: tache.id, projet })"
            >
              Restaurer
            </Bouton>
            <Bouton taille="petite" variante="danger" @click="confirmation = tache.id">Effacer</Bouton>
          </div>
        </div>
        <div
          v-if="confirmation === tache.id"
          class="mt-2 rounded-lg border border-red-300 bg-red-50 p-2 dark:border-red-900 dark:bg-red-950/40"
        >
          <p class="text-xs text-red-800 dark:text-red-300">
            Effacer définitivement cette tâche et ses pièces jointes ? Cette action est irréversible.
          </p>
          <div class="mt-2 flex gap-2">
            <Bouton
              taille="petite"
              variante="danger"
              :desactive="purge.isPending.value"
              @click="purge.mutate({ id: tache.id, projet }), (confirmation = '')"
            >
              Oui, effacer
            </Bouton>
            <Bouton taille="petite" variante="secondaire" @click="confirmation = ''">Annuler</Bouton>
          </div>
        </div>
      </li>
    </ul>

    <template #pied>
      <Bouton variante="secondaire" @click="emissions('fermer')">Fermer</Bouton>
    </template>
  </Dialogue>
</template>
