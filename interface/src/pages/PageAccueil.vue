<script setup lang="ts">
import { ref } from "vue"
import { utiliserGroupes, utiliserMutationGroupe } from "@/api/requetes"
import { libellesRoles } from "@/api/types"
import { formulaireGroupe, valider } from "@/utilitaires/validation"
import { extraireErreur } from "@/utilitaires/erreurs"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import Zone from "@/composants/ui/Zone.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import Badge from "@/composants/ui/Badge.vue"

const { data: groupes, isLoading: chargement } = utiliserGroupes()
const mutation = utiliserMutationGroupe()

const dialogueOuvert = ref(false)
const nom = ref("")
const description = ref("")
const erreurs = ref<Record<string, string>>({})
const erreurApi = ref("")

function creer() {
  erreurApi.value = ""
  const resultat = valider(formulaireGroupe, { nom: nom.value, description: description.value })
  erreurs.value = resultat.erreurs
  if (!resultat.donnees) return
  mutation.mutate(resultat.donnees, {
    onSuccess: () => {
      dialogueOuvert.value = false
      nom.value = ""
      description.value = ""
    },
    onError: (erreur) => {
      erreurApi.value = extraireErreur(erreur)
    },
  })
}
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-8">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-xl font-bold sm:text-2xl">Mes groupes</h1>
        <p class="mt-1 text-sm text-neutral-500">Retrouvez vos équipes et leurs projets kanban.</p>
      </div>
      <Bouton @click="dialogueOuvert = true">Nouveau groupe</Bouton>
    </div>

    <p v-if="chargement" class="mt-12 text-center text-sm text-neutral-500">Chargement…</p>

    <div v-else-if="groupes?.length" class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <RouterLink
        v-for="groupe in groupes"
        :key="groupe.id"
        :to="`/groupes/${groupe.id}`"
        class="group rounded-xl border border-neutral-200 bg-white p-5 transition-shadow hover:shadow-md dark:border-neutral-800 dark:bg-neutral-900"
      >
        <div class="flex items-start justify-between gap-2">
          <h2 class="font-semibold group-hover:underline">{{ groupe.nom }}</h2>
          <Badge>{{ libellesRoles[groupe.role] ?? groupe.role }}</Badge>
        </div>
        <p class="mt-1 line-clamp-2 min-h-10 text-sm text-neutral-500">{{ groupe.description }}</p>
        <div class="mt-4 flex gap-4 text-xs text-neutral-500">
          <span>{{ groupe.nbmembres }} membre{{ groupe.nbmembres > 1 ? "s" : "" }}</span>
          <span>{{ groupe.nbprojets }} projet{{ groupe.nbprojets > 1 ? "s" : "" }}</span>
        </div>
      </RouterLink>
    </div>

    <div
      v-else
      class="mt-12 rounded-xl border border-dashed border-neutral-300 p-12 text-center dark:border-neutral-700"
    >
      <p class="font-medium">Aucun groupe pour le moment</p>
      <p class="mt-1 text-sm text-neutral-500">Créez votre premier groupe pour démarrer un tableau kanban.</p>
      <Bouton class="mt-4" @click="dialogueOuvert = true">Créer un groupe</Bouton>
    </div>

    <Dialogue :ouvert="dialogueOuvert" titre="Nouveau groupe" @fermer="dialogueOuvert = false">
      <form class="space-y-4" @submit.prevent="creer">
        <Champ v-model="nom" etiquette="Nom" obligatoire indication="Équipe produit" :erreur="erreurs.nom" />
        <Zone v-model="description" etiquette="Description" indication="À quoi sert ce groupe ?" :erreur="erreurs.description" />
        <p v-if="erreurApi" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">
          {{ erreurApi }}
        </p>
      </form>
      <template #pied>
        <Bouton variante="secondaire" @click="dialogueOuvert = false">Annuler</Bouton>
        <Bouton :desactive="mutation.isPending.value" @click="creer">Créer</Bouton>
      </template>
    </Dialogue>
  </div>
</template>
