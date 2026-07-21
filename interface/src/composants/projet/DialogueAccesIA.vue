<script setup lang="ts">
import { computed, ref, toRef } from "vue"
import { useClipboard } from "@vueuse/core"
import { utiliserCle, utiliserGenerationCle, utiliserSuppressionCle } from "@/api/requetes"
import { formaterDateHeure } from "@/utilitaires/dates"
import Bouton from "@/composants/ui/Bouton.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"

const proprietes = defineProps<{ ouvert: boolean; projet: string }>()
const emissions = defineEmits<{ fermer: [] }>()

const ouvert = toRef(proprietes, "ouvert")
const projet = toRef(proprietes, "projet")

const { data: cle, isLoading: chargement } = utiliserCle(projet, ouvert)
const generation = utiliserGenerationCle()
const suppression = utiliserSuppressionCle()

const { copy: copier, copied: copiee } = useClipboard()

const base = computed(() => `${window.location.origin}/api/agent`)
const confirmationRegeneration = ref(false)

function generer() {
  confirmationRegeneration.value = false
  generation.mutate(projet.value)
}

const pointsEntree = [
  { methode: "GET", chemin: "/projet", detail: "Colonnes, membres, étiquettes et lots du projet" },
  { methode: "GET", chemin: "/taches", detail: "Toutes les tâches du tableau" },
  { methode: "POST", chemin: "/taches", detail: "Créer une tâche : { titre, colonne, description?, points?, echeance?, affectations? }" },
  { methode: "GET", chemin: "/taches/{id}", detail: "Détail d'une tâche" },
  { methode: "PUT", chemin: "/taches/{id}", detail: "Modifier une tâche : { titre, description?, points?, echeance?, commit? }" },
  { methode: "PUT", chemin: "/taches/{id}/deplacer", detail: "Déplacer vers une colonne : { colonne, position }" },
  { methode: "PUT", chemin: "/taches/{id}/commit", detail: "Renseigner le commit une fois la tâche terminée : { commit }" },
  { methode: "GET", chemin: "/taches/{id}/commentaires", detail: "Lire les commentaires" },
  { methode: "POST", chemin: "/taches/{id}/commentaires", detail: "Commenter : { contenu }" },
  { methode: "GET", chemin: "/taches/{id}/activites", detail: "Journal d'activité de la tâche" },
]

const exemple = computed(
  () =>
    `curl -k ${base.value}/taches \\\n  -H "X-Cle-API: ${cle.value?.cle ?? "<votre cle>"}" \\\n  -H "Content-Type: application/json" \\\n  -d '{"titre":"Ma tache","colonne":"<id colonne>"}' -X POST`,
)
</script>

<template>
  <Dialogue :ouvert="ouvert" titre="Accès IA" geant @fermer="emissions('fermer')">
    <div class="space-y-5">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        Cette clé permet à un agent (IA, script, intégration continue…) d'agir sur ce projet en votre nom :
        créer, déplacer et commenter des tâches, ou renseigner le commit d'une tâche terminée.
        Elle est personnelle et limitée à ce projet.
      </p>

      <div class="rounded-xl border border-neutral-200 bg-neutral-50 p-4 dark:border-neutral-700 dark:bg-neutral-800/50">
        <div class="flex items-center justify-between gap-3">
          <span class="text-sm font-semibold">Votre clé d'API</span>
          <span v-if="cle" class="text-xs text-neutral-500">générée le {{ formaterDateHeure(cle.creation) }}</span>
        </div>
        <template v-if="cle">
          <div class="mt-3 flex items-center gap-2">
            <code
              class="min-w-0 flex-1 overflow-x-auto whitespace-nowrap rounded-lg border border-neutral-300 bg-white px-3 py-2 font-mono text-xs dark:border-neutral-700 dark:bg-neutral-900"
            >
              {{ cle.cle }}
            </code>
            <Bouton taille="petite" variante="secondaire" @click="copier(cle.cle)">
              {{ copiee ? "Copiée !" : "Copier" }}
            </Bouton>
          </div>
          <div class="mt-3 flex flex-wrap items-center gap-2">
            <template v-if="confirmationRegeneration">
              <span class="text-xs text-red-700 dark:text-red-400">L'ancienne clé cessera de fonctionner.</span>
              <Bouton taille="petite" variante="danger" :desactive="generation.isPending.value" @click="generer">
                Confirmer la régénération
              </Bouton>
              <Bouton taille="petite" variante="secondaire" @click="confirmationRegeneration = false">Annuler</Bouton>
            </template>
            <template v-else>
              <Bouton taille="petite" variante="secondaire" @click="confirmationRegeneration = true">Régénérer</Bouton>
              <Bouton
                taille="petite"
                variante="secondaire"
                :desactive="suppression.isPending.value"
                @click="suppression.mutate(projet)"
              >
                Révoquer
              </Bouton>
            </template>
          </div>
        </template>
        <template v-else>
          <p class="mt-2 text-sm text-neutral-500">
            {{ chargement ? "Chargement…" : "Aucune clé pour ce projet." }}
          </p>
          <Bouton
            v-if="!chargement"
            class="mt-3"
            taille="petite"
            :desactive="generation.isPending.value"
            @click="generer"
          >
            Générer une clé
          </Bouton>
        </template>
      </div>

      <div>
        <h3 class="text-sm font-semibold">Utilisation</h3>
        <p class="mt-1 text-sm text-neutral-600 dark:text-neutral-400">
          Adresse de base : <code class="rounded bg-neutral-100 px-1.5 py-0.5 font-mono text-xs dark:bg-neutral-800">{{ base }}</code>
          — authentification par l'en-tête
          <code class="rounded bg-neutral-100 px-1.5 py-0.5 font-mono text-xs dark:bg-neutral-800">X-Cle-API</code>
          ou
          <code class="rounded bg-neutral-100 px-1.5 py-0.5 font-mono text-xs dark:bg-neutral-800">Authorization: Bearer &lt;cle&gt;</code>.
        </p>
        <div class="mt-3 overflow-x-auto rounded-xl border border-neutral-200 dark:border-neutral-700">
          <table class="w-full min-w-[36rem] text-left text-xs">
            <tbody>
              <tr
                v-for="point in pointsEntree"
                :key="point.methode + point.chemin"
                class="border-b border-neutral-100 last:border-0 dark:border-neutral-800"
              >
                <td class="w-14 px-3 py-2 font-mono font-semibold text-neutral-500">{{ point.methode }}</td>
                <td class="whitespace-nowrap px-3 py-2 font-mono">{{ point.chemin }}</td>
                <td class="px-3 py-2 text-neutral-600 dark:text-neutral-400">{{ point.detail }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <pre
          class="mt-3 overflow-x-auto rounded-xl bg-neutral-900 p-4 font-mono text-xs leading-relaxed text-neutral-100 dark:bg-neutral-800"
        >{{ exemple }}</pre>
      </div>
    </div>
    <template #pied>
      <Bouton variante="secondaire" @click="emissions('fermer')">Fermer</Bouton>
    </template>
  </Dialogue>
</template>
