<script setup lang="ts">
import { computed, ref } from "vue"
import type { Tache } from "@/api/types"
import {
  televerserFichier,
  utiliserCommentaires,
  utiliserMutationCommentaire,
  utiliserProfil,
  utiliserSuppressionCommentaire,
} from "@/api/requetes"
import { depuis } from "@/utilitaires/dates"
import { assainir, contenuVide } from "@/utilitaires/html"
import Avatar from "@/composants/ui/Avatar.vue"
import Bouton from "@/composants/ui/Bouton.vue"
import ZoneRiche from "@/composants/ui/ZoneRiche.vue"

const proprietes = defineProps<{
  tache: Tache
  projet: string
  edition: boolean
}>()

const identifiant = computed(() => proprietes.tache.id)
const { data: profil } = utiliserProfil()
const { data: commentaires } = utiliserCommentaires(identifiant)
const mutation = utiliserMutationCommentaire()
const suppression = utiliserSuppressionCommentaire()

const nouveau = ref("")

function envoyerImage(fichier: File) {
  return televerserFichier(proprietes.projet, fichier)
}

function commenter() {
  if (contenuVide(nouveau.value)) return
  mutation.mutate(
    { tache: proprietes.tache.id, contenu: assainir(nouveau.value) },
    { onSuccess: () => (nouveau.value = "") },
  )
}
</script>

<template>
  <div class="space-y-4">
    <ul v-if="commentaires?.length" class="space-y-4">
      <li v-for="commentaire in commentaires" :key="commentaire.id" class="flex gap-3">
        <Avatar :nom="commentaire.nom" :prenom="commentaire.prenom" />
        <div class="min-w-0 flex-1 rounded-lg bg-neutral-100 px-3 py-2 dark:bg-neutral-800/60">
          <p class="text-xs text-neutral-500">
            <span class="font-medium text-neutral-700 dark:text-neutral-300">
              {{ commentaire.prenom }} {{ commentaire.nom }}
            </span>
            · {{ depuis(commentaire.creation) }}
            <button
              v-if="commentaire.auteur === profil?.id"
              class="ml-1 text-red-700 hover:underline dark:text-red-500"
              @click="suppression.mutate({ id: commentaire.id, tache: tache.id })"
            >
              Supprimer
            </button>
          </p>
          <div class="texteriche mt-1 text-sm" v-html="assainir(commentaire.contenu)"></div>
        </div>
      </li>
    </ul>
    <p v-else class="rounded-lg border border-dashed border-neutral-300 p-6 text-center text-sm text-neutral-500 dark:border-neutral-700">
      Aucun commentaire pour le moment.
    </p>

    <form v-if="edition" class="space-y-2 border-t border-neutral-200 pt-4 dark:border-neutral-800" @submit.prevent="commenter">
      <ZoneRiche v-model="nouveau" compact indication="Écrire un commentaire…" :televerser="envoyerImage" />
      <div class="flex justify-end">
        <Bouton taille="petite" type="submit" :desactive="mutation.isPending.value">Envoyer</Bouton>
      </div>
    </form>
  </div>
</template>
