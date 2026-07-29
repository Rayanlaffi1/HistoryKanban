<script setup lang="ts">
import { computed, reactive, watch } from "vue"
import { useDark } from "@vueuse/core"
import {
  utiliserEtatMaj,
  utiliserMutationMaj,
  utiliserMutationPreferences,
  utiliserPreferences,
  utiliserProfil,
} from "@/api/requetes"
import { typesNotifications } from "@/api/types"
import Bascule from "@/composants/ui/Bascule.vue"
import Bouton from "@/composants/ui/Bouton.vue"
import { utiliserMagasinNotifications } from "@/magasins/notifications"
import { notifierDeplacement } from "@/tempsreel/prise"

const sombre = useDark()
const { data: profil } = utiliserProfil()
const { data: preferences } = utiliserPreferences()
const mutation = utiliserMutationPreferences()
const magasin = utiliserMagasinNotifications()

const formulaire = reactive({
  courriels: true,
  types: {} as Record<string, boolean>,
})

watch(
  preferences,
  (valeur) => {
    if (!valeur) return
    formulaire.courriels = valeur.courriels
    for (const type of Object.keys(typesNotifications)) {
      formulaire.types[type] = valeur.types[type] ?? true
    }
  },
  { immediate: true },
)

function enregistrer() {
  mutation.mutate(
    { courriels: formulaire.courriels, types: { ...formulaire.types } },
    { onSuccess: () => magasin.annoncer("Préférences enregistrées", "") },
  )
}

const estAdministrateur = computed(() => profil.value?.roles?.includes("administrateur") ?? false)
const { data: etatMaj, isLoading: etatMajEnCours } = utiliserEtatMaj(estAdministrateur)
const mutationMaj = utiliserMutationMaj()

function mettreAJour() {
  mutationMaj.mutate(undefined, {
    onSuccess: () => magasin.annoncer("Mise à jour lancée", "Le pipeline Jenkins reconstruit et redéploie l'application."),
  })
}
</script>

<template>
  <div class="mx-auto max-w-2xl px-4 py-8">
    <h1 class="text-2xl font-bold">Paramètres</h1>

    <section class="mt-8 rounded-xl border border-neutral-200 bg-white p-6 dark:border-neutral-800 dark:bg-neutral-900">
      <h2 class="font-semibold">Profil</h2>
      <div v-if="profil" class="mt-4 space-y-1 text-sm">
        <p>{{ profil.prenom }} {{ profil.nom }}</p>
        <p class="text-neutral-500">{{ profil.courriel }}</p>
      </div>
      <p class="mt-3 text-xs text-neutral-500">
        Le profil et le mot de passe se gèrent dans Keycloak.
      </p>
    </section>

    <section class="mt-6 rounded-xl border border-neutral-200 bg-white p-6 dark:border-neutral-800 dark:bg-neutral-900">
      <h2 class="font-semibold">Apparence</h2>
      <div class="mt-4">
        <Bascule v-model="sombre" etiquette="Mode sombre" />
      </div>
    </section>

    <section class="mt-6 rounded-xl border border-neutral-200 bg-white p-6 dark:border-neutral-800 dark:bg-neutral-900">
      <h2 class="font-semibold">Notifications de l'application</h2>
      <p class="mt-1 text-sm text-neutral-500">Alertes instantanées affichées dans l'interface.</p>
      <div class="mt-4">
        <Bascule v-model="notifierDeplacement" etiquette="Tâche déplacée d'une colonne à une autre" />
      </div>
    </section>

    <section class="mt-6 rounded-xl border border-neutral-200 bg-white p-6 dark:border-neutral-800 dark:bg-neutral-900">
      <h2 class="font-semibold">Courriels de notification</h2>
      <p class="mt-1 text-sm text-neutral-500">
        Choisissez les changements qui déclenchent l'envoi d'un courriel.
      </p>
      <div class="mt-4 space-y-3">
        <Bascule v-model="formulaire.courriels" etiquette="Activer les courriels" />
        <div
          class="space-y-3 border-t border-neutral-200 pt-3 dark:border-neutral-800"
          :class="!formulaire.courriels && 'pointer-events-none opacity-50'"
        >
          <Bascule
            v-for="(libelle, type) in typesNotifications"
            :key="type"
            v-model="formulaire.types[type]"
            :etiquette="libelle"
          />
        </div>
      </div>
      <div class="mt-6 flex justify-end">
        <Bouton :desactive="mutation.isPending.value" @click="enregistrer">Enregistrer</Bouton>
      </div>
    </section>

    <section
      v-if="estAdministrateur"
      class="mt-6 rounded-xl border border-neutral-200 bg-white p-6 dark:border-neutral-800 dark:bg-neutral-900"
    >
      <h2 class="font-semibold">Mise à jour de l'application</h2>
      <p class="mt-1 text-sm text-neutral-500">
        Déclenche le pipeline Jenkins qui reconstruit et redéploie l'application depuis la dernière version publiée.
      </p>
      <div class="mt-4 space-y-1 text-sm">
        <p v-if="etatMajEnCours" class="text-neutral-500">Vérification de la version…</p>
        <template v-else-if="etatMaj">
          <p>Version installée : {{ etatMaj.versionActuelle }}</p>
          <p v-if="etatMaj.derniereVersion">Dernière release GitHub : {{ etatMaj.derniereVersion }}</p>
          <p v-if="etatMaj.majDisponible" class="font-medium text-amber-700 dark:text-amber-500">
            Une nouvelle release est disponible.
          </p>
          <p v-else-if="etatMaj.derniereVersion" class="text-neutral-500">L'application est à jour.</p>
          <p v-if="!etatMaj.jenkinsConfigure" class="text-red-700 dark:text-red-500">
            Le déclenchement Jenkins n'est pas configuré sur le serveur.
          </p>
        </template>
      </div>
      <div class="mt-6 flex items-center justify-end gap-3">
        <p v-if="mutationMaj.isPending.value" class="text-sm text-neutral-500">Déclenchement en cours…</p>
        <p v-else-if="mutationMaj.isSuccess.value" class="text-sm text-green-700 dark:text-green-500">
          Mise à jour lancée.
        </p>
        <p v-else-if="mutationMaj.isError.value" class="text-sm text-red-700 dark:text-red-500">
          Le déclenchement a échoué.
        </p>
        <Bouton
          :desactive="mutationMaj.isPending.value || !etatMaj?.jenkinsConfigure"
          @click="mettreAJour"
        >
          Mettre à jour l'application
        </Bouton>
      </div>
    </section>
  </div>
</template>
