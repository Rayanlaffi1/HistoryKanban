<script setup lang="ts">
import { computed, ref } from "vue"
import { useDark, useToggle, onClickOutside } from "@vueuse/core"
import { utiliserProfil, utiliserNotifications } from "@/api/requetes"
import { seDeconnecter } from "@/securite/keycloak"
import Avatar from "@/composants/ui/Avatar.vue"
import PanneauNotifications from "@/composants/PanneauNotifications.vue"
import RechercheGlobale from "@/composants/RechercheGlobale.vue"

const sombre = useDark()
const basculerTheme = useToggle(sombre)

const { data: profil } = utiliserProfil()
const { data: notifications } = utiliserNotifications()

const nonLues = computed(() => notifications.value?.filter((notification) => !notification.lue).length ?? 0)

const panneauOuvert = ref(false)
const zonePanneau = ref<HTMLElement | null>(null)
onClickOutside(zonePanneau, () => (panneauOuvert.value = false))
</script>

<template>
  <header
    class="sticky top-0 z-40 border-b border-neutral-200 bg-white/90 backdrop-blur dark:border-neutral-800 dark:bg-neutral-950/90"
  >
    <div class="mx-auto flex h-14 max-w-7xl items-center justify-between gap-2 px-3 sm:px-4">
      <div class="flex min-w-0 items-center gap-2 sm:gap-6">
        <RouterLink to="/" class="flex shrink-0 items-center gap-2 text-base font-bold tracking-tight sm:text-lg">
          <img src="/icone-192.png" alt="" class="h-7 w-7" />
          <span class="hidden min-[420px]:inline">HistoryKanban</span>
        </RouterLink>
        <nav class="flex items-center gap-0.5 text-sm sm:gap-1">
          <RouterLink
            to="/"
            class="rounded-md px-2 py-1.5 text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900 dark:text-neutral-400 dark:hover:bg-neutral-800 dark:hover:text-neutral-100 sm:px-3"
            active-class="bg-neutral-100 !text-neutral-900 dark:bg-neutral-800 dark:!text-neutral-100"
          >
            Accueil
          </RouterLink>
          <RouterLink
            to="/parametres"
            class="rounded-md px-2 py-1.5 text-neutral-600 hover:bg-neutral-100 hover:text-neutral-900 dark:text-neutral-400 dark:hover:bg-neutral-800 dark:hover:text-neutral-100 sm:px-3"
            active-class="bg-neutral-100 !text-neutral-900 dark:bg-neutral-800 dark:!text-neutral-100"
          >
            Paramètres
          </RouterLink>
        </nav>
      </div>
      <div class="flex shrink-0 items-center gap-1 sm:gap-2">
        <RechercheGlobale />
        <button
          class="rounded-md p-2 text-neutral-600 hover:bg-neutral-100 dark:text-neutral-400 dark:hover:bg-neutral-800"
          :title="sombre ? 'Mode clair' : 'Mode sombre'"
          @click="basculerTheme()"
        >
          <svg v-if="sombre" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2m0 16v2M4.9 4.9l1.4 1.4m11.4 11.4 1.4 1.4M2 12h2m16 0h2M4.9 19.1l1.4-1.4m11.4-11.4 1.4-1.4" />
          </svg>
          <svg v-else class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8Z" />
          </svg>
        </button>
        <div ref="zonePanneau" class="relative">
          <button
            class="relative rounded-md p-2 text-neutral-600 hover:bg-neutral-100 dark:text-neutral-400 dark:hover:bg-neutral-800"
            title="Notifications"
            @click="panneauOuvert = !panneauOuvert"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9" />
              <path d="M10.3 21a1.9 1.9 0 0 0 3.4 0" />
            </svg>
            <span
              v-if="nonLues > 0"
              class="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-neutral-900 px-1 text-[10px] font-bold text-white dark:bg-neutral-100 dark:text-neutral-900"
            >
              {{ nonLues }}
            </span>
          </button>
          <PanneauNotifications v-if="panneauOuvert" @fermer="panneauOuvert = false" />
        </div>
        <div class="ml-1 flex items-center gap-2 sm:ml-2 sm:gap-3">
          <Avatar v-if="profil" class="hidden sm:inline-flex" :nom="profil.nom" :prenom="profil.prenom" />
          <button
            class="text-xs text-neutral-600 hover:text-neutral-900 hover:underline dark:text-neutral-400 dark:hover:text-neutral-100 sm:text-sm"
            @click="seDeconnecter()"
          >
            Déconnexion
          </button>
        </div>
      </div>
    </div>
  </header>
</template>
