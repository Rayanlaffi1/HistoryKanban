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
import { formulaireGroupe, formulaireMembre, formulaireProjet, valider } from "@/utilitaires/validation"
import { extraireErreur } from "@/utilitaires/erreurs"
import { connectes } from "@/tempsreel/prise"
import Bouton from "@/composants/ui/Bouton.vue"
import Champ from "@/composants/ui/Champ.vue"
import ChampCouleur from "@/composants/ui/ChampCouleur.vue"
import Zone from "@/composants/ui/Zone.vue"
import Dialogue from "@/composants/ui/Dialogue.vue"
import Selection from "@/composants/ui/Selection.vue"
import Avatar from "@/composants/ui/Avatar.vue"
import Badge from "@/composants/ui/Badge.vue"
import BoutonRetour from "@/composants/ui/BoutonRetour.vue"

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
const erreursMembre = ref<Record<string, string>>({})
const erreurMembre = ref("")

function ajouterMembre() {
  erreurMembre.value = ""
  const resultat = valider(formulaireMembre, { courriel: courrielMembre.value, role: roleMembre.value })
  erreursMembre.value = resultat.erreurs
  if (!resultat.donnees) return
  mutationMembre.mutate(
    { groupe: identifiant.value, ...resultat.donnees },
    {
      onSuccess: () => {
        dialogueMembre.value = false
        courrielMembre.value = ""
      },
      onError: (erreur) => {
        erreurMembre.value = extraireErreur(erreur)
      },
    },
  )
}

const dialogueProjet = ref(false)
const nomProjet = ref("")
const descriptionProjet = ref("")
const couleurProjet = ref("#737373")
const erreursProjet = ref<Record<string, string>>({})
const erreurProjet = ref("")

function creerProjet() {
  erreurProjet.value = ""
  const resultat = valider(formulaireProjet, {
    nom: nomProjet.value,
    description: descriptionProjet.value,
    couleur: couleurProjet.value,
  })
  erreursProjet.value = resultat.erreurs
  if (!resultat.donnees) return
  mutationProjet.mutate(
    { groupe: identifiant.value, ...resultat.donnees },
    {
      onSuccess: () => {
        dialogueProjet.value = false
        nomProjet.value = ""
        descriptionProjet.value = ""
      },
      onError: (erreur) => {
        erreurProjet.value = extraireErreur(erreur)
      },
    },
  )
}

const dialogueEdition = ref(false)
const nomEdition = ref("")
const descriptionEdition = ref("")
const erreursEdition = ref<Record<string, string>>({})
const erreurEdition = ref("")

function ouvrirEdition() {
  nomEdition.value = groupe.value?.nom ?? ""
  descriptionEdition.value = groupe.value?.description ?? ""
  erreursEdition.value = {}
  erreurEdition.value = ""
  dialogueEdition.value = true
}

function modifierGroupe() {
  erreurEdition.value = ""
  const resultat = valider(formulaireGroupe, { nom: nomEdition.value, description: descriptionEdition.value })
  erreursEdition.value = resultat.erreurs
  if (!resultat.donnees) return
  mutationGroupe.mutate(
    { id: identifiant.value, ...resultat.donnees },
    {
      onSuccess: () => (dialogueEdition.value = false),
      onError: (erreur) => {
        erreurEdition.value = extraireErreur(erreur)
      },
    },
  )
}

const confirmationSuppression = ref(false)

const nbEnLigne = computed(
  () => (membres.value ?? []).filter((membre) => connectes.value.includes(membre.utilisateur)).length,
)

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
    <BoutonRetour vers="/" etiquette="Retour aux groupes" />

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
          <h2 class="flex items-baseline gap-2 text-lg font-semibold">
            Membres
            <span class="text-xs font-normal text-neutral-500">
              {{ nbEnLigne }} en ligne
            </span>
          </h2>
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
              <span class="relative inline-flex shrink-0">
                <Avatar :nom="membre.nom" :prenom="membre.prenom" />
                <span
                  class="absolute -bottom-0.5 -right-0.5 h-2.5 w-2.5 rounded-full border-2 border-white dark:border-neutral-900"
                  :class="connectes.includes(membre.utilisateur) ? 'bg-green-500' : 'bg-neutral-400'"
                  :title="connectes.includes(membre.utilisateur) ? 'En ligne' : 'Hors ligne'"
                ></span>
              </span>
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
        <div class="grid gap-4 sm:grid-cols-[1fr,10rem]">
          <Champ v-model="courrielMembre" etiquette="Courriel" type="email" obligatoire indication="collegue@exemple.fr" :erreur="erreursMembre.courriel" />
          <Selection v-model="roleMembre" etiquette="Rôle">
            <option v-for="role in roles" :key="role" :value="role">{{ libellesRoles[role] }}</option>
          </Selection>
        </div>
        <p v-if="erreurMembre" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">{{ erreurMembre }}</p>
        <p class="rounded-lg bg-neutral-100 px-3 py-2 text-xs text-neutral-500 dark:bg-neutral-800/60">
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
        <Champ v-model="nomProjet" etiquette="Nom" obligatoire indication="Refonte du site" :erreur="erreursProjet.nom" />
        <Zone v-model="descriptionProjet" etiquette="Description" indication="Objectif du projet…" :erreur="erreursProjet.description" />
        <ChampCouleur v-model="couleurProjet" etiquette="Couleur du projet" />
        <p v-if="erreurProjet" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">{{ erreurProjet }}</p>
      </form>
      <template #pied>
        <Bouton variante="secondaire" @click="dialogueProjet = false">Annuler</Bouton>
        <Bouton :desactive="mutationProjet.isPending.value" @click="creerProjet">Créer</Bouton>
      </template>
    </Dialogue>

    <Dialogue :ouvert="dialogueEdition" titre="Modifier le groupe" @fermer="dialogueEdition = false">
      <form class="space-y-4" @submit.prevent="modifierGroupe">
        <Champ v-model="nomEdition" etiquette="Nom" :erreur="erreursEdition.nom" />
        <Zone v-model="descriptionEdition" etiquette="Description" :erreur="erreursEdition.description" />
        <p v-if="erreurEdition" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-800 dark:bg-red-950/40 dark:text-red-300">{{ erreurEdition }}</p>
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
