<script setup lang="ts">
import { reactive, ref, watch } from "vue"
import type { Colonne } from "@/api/types"
import { utiliserMutationColonne, utiliserSuppressionColonne } from "@/api/requetes"
import { formulaireColonne, valider } from "@/utilitaires/validation"
import { extraireErreur } from "@/utilitaires/erreurs"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import ChampCouleur from "@/composants/ui/ChampCouleur.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import SectionFormulaire from "@/composants/ui/SectionFormulaire.vue"

const proprietes = defineProps<{
  ouvert: boolean
  projet: string
  colonne: Colonne | null
}>()

const emissions = defineEmits<{ fermer: [] }>()

const formulaire = reactive({ nom: "", couleur: "#a3a3a3", limite: "" })
const confirmation = ref(false)
const erreurs = ref<Record<string, string>>({})
const erreurApi = ref("")

watch(
  () => proprietes.ouvert,
  (ouvert) => {
    if (!ouvert) return
    confirmation.value = false
    erreurs.value = {}
    erreurApi.value = ""
    formulaire.nom = proprietes.colonne?.nom ?? ""
    formulaire.couleur = proprietes.colonne?.couleur ?? "#a3a3a3"
    formulaire.limite = proprietes.colonne?.limite?.toString() ?? ""
  },
  { immediate: true },
)

const mutation = utiliserMutationColonne()
const suppression = utiliserSuppressionColonne()

function enregistrer() {
  erreurApi.value = ""
  const resultat = valider(formulaireColonne, { ...formulaire })
  erreurs.value = resultat.erreurs
  if (!resultat.donnees) return
  mutation.mutate(
    {
      id: proprietes.colonne?.id,
      projet: proprietes.projet,
      nom: resultat.donnees.nom,
      couleur: resultat.donnees.couleur,
      limite: resultat.donnees.limite === "" ? null : Number(resultat.donnees.limite),
    },
    {
      onSuccess: () => emissions("fermer"),
      onError: (erreur) => {
        erreurApi.value = extraireErreur(erreur)
      },
    },
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
      <Champ v-model="formulaire.nom" etiquette="Nom" obligatoire indication="En revue" :erreur="erreurs.nom" />
      <SectionFormulaire titre="Apparence et règles">
        <ChampCouleur v-model="formulaire.couleur" etiquette="Couleur" />
        <Champ v-model="formulaire.limite" etiquette="Limite de tâches" type="number" indication="Aucune limite" :erreur="erreurs.limite" />
      </SectionFormulaire>
      <p v-if="erreurApi" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">{{ erreurApi }}</p>
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
