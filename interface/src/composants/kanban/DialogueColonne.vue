<script setup lang="ts">
import { reactive, ref, watch } from "vue"
import type { Colonne } from "@/api/types"
import { utiliserMutationColonne, utiliserSuppressionColonne } from "@/api/requetes"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"

const proprietes = defineProps<{
  ouvert: boolean
  projet: string
  colonne: Colonne | null
}>()

const emissions = defineEmits<{ fermer: [] }>()

const formulaire = reactive({ nom: "", couleur: "#a3a3a3", limite: "" })
const confirmation = ref(false)

watch(
  () => proprietes.ouvert,
  (ouvert) => {
    if (!ouvert) return
    confirmation.value = false
    formulaire.nom = proprietes.colonne?.nom ?? ""
    formulaire.couleur = proprietes.colonne?.couleur ?? "#a3a3a3"
    formulaire.limite = proprietes.colonne?.limite?.toString() ?? ""
  },
  { immediate: true },
)

const mutation = utiliserMutationColonne()
const suppression = utiliserSuppressionColonne()

function enregistrer() {
  if (!formulaire.nom.trim()) return
  mutation.mutate(
    {
      id: proprietes.colonne?.id,
      projet: proprietes.projet,
      nom: formulaire.nom.trim(),
      couleur: formulaire.couleur,
      limite: formulaire.limite === "" ? null : Number(formulaire.limite),
    },
    { onSuccess: () => emissions("fermer") },
  )
}

function supprimer() {
  if (!proprietes.colonne) return
  suppression.mutate(
    { id: proprietes.colonne.id, projet: proprietes.projet },
    { onSuccess: () => emissions("fermer") },
  )
}
</script>

<template>
  <Dialogue :ouvert="ouvert" :titre="colonne ? 'Modifier la colonne' : 'Nouvelle colonne'" @fermer="emissions('fermer')">
    <form class="space-y-4" @submit.prevent="enregistrer">
      <Champ v-model="formulaire.nom" etiquette="Nom" obligatoire indication="En revue" />
      <div class="grid grid-cols-2 gap-4">
        <label class="block space-y-1.5">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Couleur</span>
          <input v-model="formulaire.couleur" type="color" class="h-10 w-20 cursor-pointer rounded-lg border border-neutral-300 bg-white dark:border-neutral-700 dark:bg-neutral-900" />
        </label>
        <Champ v-model="formulaire.limite" etiquette="Limite de tâches" type="number" indication="Aucune" />
      </div>
      <div v-if="colonne" class="border-t border-neutral-200 pt-3 dark:border-neutral-800">
        <button v-if="!confirmation" type="button" class="text-sm text-red-700 hover:underline dark:text-red-500" @click="confirmation = true">
          Supprimer la colonne et ses tâches
        </button>
        <div v-else class="flex items-center gap-2">
          <span class="text-sm text-red-700 dark:text-red-500">Confirmer la suppression ?</span>
          <Bouton taille="petite" variante="danger" :desactive="suppression.isPending.value" @click="supprimer">Oui</Bouton>
          <Bouton taille="petite" variante="secondaire" @click="confirmation = false">Non</Bouton>
        </div>
      </div>
    </form>
    <template #pied>
      <Bouton variante="secondaire" @click="emissions('fermer')">Annuler</Bouton>
      <Bouton :desactive="mutation.isPending.value" @click="enregistrer">{{ colonne ? "Enregistrer" : "Créer" }}</Bouton>
    </template>
  </Dialogue>
</template>
