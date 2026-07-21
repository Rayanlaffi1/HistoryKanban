<script setup lang="ts">
import { computed, ref } from "vue"
import { useRoute, useRouter } from "vue-router"
import {
  utiliserGroupes,
  utiliserMembres,
  utiliserMutationGroupe,
  utiliserMutationMembre,
  utiliserMutationProjet,
  utiliserMutationRole,
  utiliserProfil,
  utiliserProjets,
  utiliserRetraitMembre,
  utiliserSuppressionGroupe,
} from "@/api/requetes"
import { libellesRoles, roles } from "@/api/types"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import Zone from "@/composants/ui/Zone.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import Selection from "@/composants/ui/Selection.vue"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"

const route = useRoute()
const routeur = useRouter()
const identifiant = computed(() => String(route.params.id))

const { data: profil } = utiliserProfil()
const { data: groupes } = utiliserGroupes()
const { data: membres } = utiliserMembres(identifiant)
const { data: projets } = utiliserProjets(identifiant)

const groupe = computed(() => groupes.value?.find((groupe) => groupe.id === identifiant.value))
const monRole = computed(() => groupe.value?.role ?? "")
const gestionnaire = computed(() => ["proprietaire", "administrateur"].includes(monRole.value))
const proprietaire = computed(() => monRole.value === "proprietaire")

const mutationGroupe = utiliserMutationGroupe()
const suppressionGroupe = utiliserSuppressionGroupe()
const mutationMembre = utiliserMutationMembre()
const mutationRole = utiliserMutationRole()
const retraitMembre = utiliserRetraitMembre()
const mutationProjet = utiliserMutationProjet()

const dialogueMembre = ref(false)
const courrielMembre = ref("")
const roleMembre = ref("membre")
const erreurMembre = ref("")

function ajouterMembre() {
  erreurMembre.value = ""
  mutationMembre.mutate(
    { groupe: identifiant.value, courriel: courrielMembre.value.trim(), role: roleMembre.value },
    {
      onSuccess: () => {
        dialogueMembre.value = false
        courrielMembre.value = ""
      },
      onError: (erreur: any) => {
        erreurMembre.value = erreur.response?.data?.erreur ?? "Ajout impossible"
      },
    },
  )
}

const dialogueProjet = ref(false)
const nomProjet = ref("")
const descriptionProjet = ref("")
const couleurProjet = ref("#737373")

function creerProjet() {
  if (!nomProjet.value.trim()) return
  mutationProjet.mutate(
    {
      groupe: identifiant.value,
      nom: nomProjet.value.trim(),
      description: descriptionProjet.value.trim(),
      couleur: couleurProjet.value,
    },
    {
      onSuccess: () => {
        dialogueProjet.value = false
        nomProjet.value = ""
        descriptionProjet.value = ""
      },
    },
  )
}

const dialogueEdition = ref(false)
const nomEdition = ref("")
const descriptionEdition = ref("")

function ouvrirEdition() {
  nomEdition.value = groupe.value?.nom ?? ""
  descriptionEdition.value = groupe.value?.description ?? ""
  dialogueEdition.value = true
}

function modifierGroupe() {
  mutationGroupe.mutate(
    { id: identifiant.value, nom: nomEdition.value.trim(), description: descriptionEdition.value.trim() },
    { onSuccess: () => (dialogueEdition.value = false) },
  )
}

const confirmationSuppression = ref(false)

function supprimerGroupe() {
  suppressionGroupe.mutate(identifiant.value, { onSuccess: () => routeur.push("/") })
}

function quitterGroupe() {
  if (!profil.value) return
  retraitMembre.mutate(
    { groupe: identifiant.value, utilisateur: profil.value.id },
    { onSuccess: () => routeur.push("/") },
  )
}
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-8">
    <RouterLink to="/" class="text-sm text-neutral-500 hover:underline">← Retour aux groupes</RouterLink>

    <div class="mt-4 flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold">{{ groupe?.nom }}</h1>
        <p class="mt-1 text-sm text-neutral-500">{{ groupe?.description }}</p>
      </div>
      <div class="flex gap-2">
        <Bouton v-if="gestionnaire" variante="secondaire" @click="ouvrirEdition">Modifier</Bouton>
        <Bouton v-if="proprietaire" variante="danger" @click="confirmationSuppression = true">Supprimer</Bouton>
        <Bouton v-else variante="secondaire" @click="quitterGroupe">Quitter le groupe</Bouton>
      </div>
    </div>

    <div class="mt-8 grid gap-8 lg:grid-cols-3">
      <section class="lg:col-span-2">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Projets</h2>
          <Bouton v-if="monRole !== 'lecteur'" taille="petite" @click="dialogueProjet = true">Nouveau projet</Bouton>
        </div>
        <div v-if="projets?.length" class="mt-4 grid gap-4 sm:grid-cols-2">
          <RouterLink
            v-for="projet in projets"
            :key="projet.id"
            :to="`/projets/${projet.id}`"
            class="group rounded-xl border border-neutral-200 bg-white p-5 transition-shadow hover:shadow-md dark:border-neutral-800 dark:bg-neutral-900"
          >
            <div class="flex items-center gap-2">
              <span class="h-3 w-3 rounded-full" :style="{ backgroundColor: projet.couleur }"></span>
              <h3 class="font-semibold group-hover:underline">{{ projet.nom }}</h3>
              <Badge v-if="projet.archive">Archivé</Badge>
            </div>
            <p class="mt-1 line-clamp-2 min-h-10 text-sm text-neutral-500">{{ projet.description }}</p>
            <p class="mt-3 text-xs text-neutral-500">{{ projet.nbtaches }} tâche{{ projet.nbtaches > 1 ? "s" : "" }}</p>
          </RouterLink>
        </div>
        <p v-else class="mt-4 rounded-xl border border-dashed border-neutral-300 p-8 text-center text-sm text-neutral-500 dark:border-neutral-700">
          Aucun projet dans ce groupe.
        </p>
      </section>

      <section>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Membres</h2>
          <Bouton v-if="gestionnaire" taille="petite" variante="secondaire" @click="dialogueMembre = true">
            Ajouter
          </Bouton>
        </div>
        <ul class="mt-4 space-y-2">
          <li
            v-for="membre in membres"
            :key="membre.utilisateur"
            class="flex items-center justify-between gap-3 rounded-lg border border-neutral-200 bg-white px-3 py-2.5 dark:border-neutral-800 dark:bg-neutral-900"
          >
            <div class="flex min-w-0 items-center gap-3">
              <Avatar :nom="membre.nom" :prenom="membre.prenom" />
              <div class="min-w-0">
                <p class="truncate text-sm font-medium">{{ membre.prenom }} {{ membre.nom }}</p>
                <p class="truncate text-xs text-neutral-500">{{ membre.courriel }}</p>
              </div>
            </div>
            <div class="flex shrink-0 items-center gap-2">
              <template v-if="proprietaire && membre.role !== 'proprietaire'">
                <select
                  :value="membre.role"
                  class="h-8 rounded-md border border-neutral-300 bg-white px-2 text-xs dark:border-neutral-700 dark:bg-neutral-900"
                  @change="
                    mutationRole.mutate({
                      groupe: identifiant,
                      utilisateur: membre.utilisateur,
                      role: ($event.target as HTMLSelectElement).value,
                    })
                  "
                >
                  <option v-for="role in roles" :key="role" :value="role">{{ libellesRoles[role] }}</option>
                </select>
                <button
                  class="text-xs text-red-700 hover:underline dark:text-red-500"
                  @click="retraitMembre.mutate({ groupe: identifiant, utilisateur: membre.utilisateur })"
                >
                  Retirer
                </button>
              </template>
              <Badge v-else>{{ libellesRoles[membre.role] ?? membre.role }}</Badge>
            </div>
          </li>
        </ul>
      </section>
    </div>

    <Dialogue :ouvert="dialogueMembre" titre="Ajouter un membre" @fermer="dialogueMembre = false">
      <form class="space-y-4" @submit.prevent="ajouterMembre">
        <Champ v-model="courrielMembre" etiquette="Courriel" type="email" obligatoire indication="collegue@exemple.fr" />
        <Selection v-model="roleMembre" etiquette="Rôle">
          <option v-for="role in roles" :key="role" :value="role">{{ libellesRoles[role] }}</option>
        </Selection>
        <p v-if="erreurMembre" class="text-sm text-red-700 dark:text-red-500">{{ erreurMembre }}</p>
        <p class="text-xs text-neutral-500">
          La personne doit s'être connectée au moins une fois à HistoryKanban.
        </p>
      </form>
      <template #pied>
        <Bouton variante="secondaire" @click="dialogueMembre = false">Annuler</Bouton>
        <Bouton :desactive="mutationMembre.isPending.value" @click="ajouterMembre">Ajouter</Bouton>
      </template>
    </Dialogue>

    <Dialogue :ouvert="dialogueProjet" titre="Nouveau projet" @fermer="dialogueProjet = false">
      <form class="space-y-4" @submit.prevent="creerProjet">
        <Champ v-model="nomProjet" etiquette="Nom" obligatoire indication="Refonte du site" />
        <Zone v-model="descriptionProjet" etiquette="Description" />
        <label class="block space-y-1.5">
          <span class="text-sm font-medium text-neutral-700 dark:text-neutral-300">Couleur</span>
          <input v-model="couleurProjet" type="color" class="h-10 w-20 cursor-pointer rounded-lg border border-neutral-300 bg-white dark:border-neutral-700 dark:bg-neutral-900" />
        </label>
      </form>
      <template #pied>
        <Bouton variante="secondaire" @click="dialogueProjet = false">Annuler</Bouton>
        <Bouton :desactive="mutationProjet.isPending.value" @click="creerProjet">Créer</Bouton>
      </template>
    </Dialogue>

    <Dialogue :ouvert="dialogueEdition" titre="Modifier le groupe" @fermer="dialogueEdition = false">
      <form class="space-y-4" @submit.prevent="modifierGroupe">
        <Champ v-model="nomEdition" etiquette="Nom" obligatoire />
        <Zone v-model="descriptionEdition" etiquette="Description" />
      </form>
      <template #pied>
        <Bouton variante="secondaire" @click="dialogueEdition = false">Annuler</Bouton>
        <Bouton :desactive="mutationGroupe.isPending.value" @click="modifierGroupe">Enregistrer</Bouton>
      </template>
    </Dialogue>

    <Dialogue :ouvert="confirmationSuppression" titre="Supprimer le groupe" @fermer="confirmationSuppression = false">
      <p class="text-sm text-neutral-600 dark:text-neutral-400">
        Cette action supprimera définitivement le groupe, ses projets et toutes leurs tâches.
      </p>
      <template #pied>
        <Bouton variante="secondaire" @click="confirmationSuppression = false">Annuler</Bouton>
        <Bouton variante="danger" :desactive="suppressionGroupe.isPending.value" @click="supprimerGroupe">
          Supprimer définitivement
        </Bouton>
      </template>
    </Dialogue>
  </div>
</template>
