<script setup lang="ts">
import { reactive, ref, watch } from "vue"
import type { Projet } from "@/api/types"
import { utiliserMutationProjet } from "@/api/requetes"
import { formulaireProjet, valider } from "@/utilitaires/validation"
import { extraireErreur } from "@/utilitaires/erreurs"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import ChampCouleur from "@/composants/ui/ChampCouleur.vue"
import Zone from "@/composants/ui/Zone.vue"
import Bascule from "@/composants/ui/Bascule.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"

const proprietes = defineProps<{ ouvert: boolean; projet: Projet | null }>()
const emissions = defineEmits<{ fermer: [] }>()

const formulaire = reactive({
  nom: "",
  description: "",
  couleur: "#737373",
  archive: false,
})

const erreurs = ref<Record<string, string>>({})
const erreurApi = ref("")

watch(
  () => proprietes.ouvert,
  (ouvert) => {
    if (!ouvert || !proprietes.projet) return
    formulaire.nom = proprietes.projet.nom
    formulaire.description = proprietes.projet.description
    formulaire.couleur = proprietes.projet.couleur
    formulaire.archive = proprietes.projet.archive
    erreurs.value = {}
    erreurApi.value = ""
  },
)

const mutation = utiliserMutationProjet()

function enregistrer() {
  if (!proprietes.projet) return
  erreurApi.value = ""
  const resultat = valider(formulaireProjet, {
    nom: formulaire.nom,
    description: formulaire.description,
    couleur: formulaire.couleur,
  })
  erreurs.value = resultat.erreurs
  if (!resultat.donnees) return
  mutation.mutate(
    {
      id: proprietes.projet.id,
      groupe: proprietes.projet.groupe,
      ...resultat.donnees,
      archive: formulaire.archive,
    },
    {
      onSuccess: () => emissions("fermer"),
      onError: (erreur) => {
        erreurApi.value = extraireErreur(erreur)
      },
    },
  )
}
</script>

<template>
  <Dialogue :ouvert="ouvert" titre="Paramètres du projet" @fermer="emissions('fermer')">
    <form class="space-y-4" @submit.prevent="enregistrer">
      <Champ v-model="formulaire.nom" etiquette="Nom" obligatoire :erreur="erreurs.nom" />
      <Zone v-model="formulaire.description" etiquette="Description" :erreur="erreurs.description" />
      <ChampCouleur v-model="formulaire.couleur" etiquette="Couleur du projet" />
      <Bascule v-model="formulaire.archive" etiquette="Projet archivé" />
      <p v-if="erreurApi" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">
        {{ erreurApi }}
      </p>
    </form>
    <template #pied>
      <Bouton variante="secondaire" @click="emissions('fermer')">Annuler</Bouton>
      <Bouton :desactive="mutation.isPending.value" @click="enregistrer">Enregistrer</Bouton>
    </template>
  </Dialogue>
</template>
