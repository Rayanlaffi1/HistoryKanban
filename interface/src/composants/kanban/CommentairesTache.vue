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
import { rendreMarkdown } from "@/utilitaires/markdown"
import Avatar from "@/composants/ui/Avatar.vue"
import Bouton from "@/composants/ui/Bouton.vue"

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
const zone = ref<HTMLTextAreaElement | null>(null)
const envoiImage = ref(false)

async function collerImage(evenement: ClipboardEvent) {
  const fichiers = Array.from(evenement.clipboardData?.files ?? []).filter((fichier) =>
    fichier.type.startsWith("image/"),
  )
  if (!fichiers.length) return
  evenement.preventDefault()
  envoiImage.value = true
  try {
    for (const fichier of fichiers) {
      const url = await televerserFichier(proprietes.projet, fichier)
      const champ = zone.value
      const position = champ?.selectionStart ?? nouveau.value.length
      const insertion = `![${fichier.name}](${url})`
      nouveau.value = nouveau.value.slice(0, position) + insertion + nouveau.value.slice(position)
    }
  } finally {
    envoiImage.value = false
  }
}

function commenter() {
  const contenu = nouveau.value.trim()
  if (!contenu) return
  mutation.mutate(
    { tache: proprietes.tache.id, contenu },
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
            <span
              v-if="commentaire.agent"
              class="ml-1 rounded border border-neutral-300 px-1 py-0.5 align-middle text-[10px] font-semibold uppercase tracking-wide text-neutral-500 dark:border-neutral-600 dark:text-neutral-400"
              title="Commentaire publié par un agent via une clé d'API"
            >
              agent
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
          <div class="texteriche mt-1 text-sm" v-html="rendreMarkdown(commentaire.contenu)"></div>
        </div>
      </li>
    </ul>
    <p v-else class="rounded-lg border border-dashed border-neutral-300 p-6 text-center text-sm text-neutral-500 dark:border-neutral-700">
      Aucun commentaire pour le moment.
    </p>

    <form v-if="edition" class="space-y-2 border-t border-neutral-200 pt-4 dark:border-neutral-800" @submit.prevent="commenter">
      <textarea
        ref="zone"
        v-model="nouveau"
        rows="3"
        placeholder="Écrire un commentaire…"
        class="w-full rounded-lg border border-neutral-300 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-neutral-400 dark:border-neutral-700 dark:bg-neutral-900 dark:focus:ring-neutral-500"
        @paste="collerImage"
      ></textarea>
      <div class="flex items-center justify-between gap-2">
        <p class="text-xs text-neutral-500">
          {{ envoiImage ? "Envoi de l'image…" : "Markdown pris en charge : **gras**, *italique*, listes, [liens](url), `code`. Coller une image l'insère." }}
        </p>
        <Bouton taille="petite" type="submit" :desactive="mutation.isPending.value">Envoyer</Bouton>
      </div>
    </form>
  </div>
</template>
