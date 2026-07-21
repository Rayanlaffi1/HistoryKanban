<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from "vue"
import dayjs from "dayjs"
import { utiliserProjets, utiliserStatistiques } from "@/api/requetes"
import type { ParametresStatistiques } from "@/api/types"
import { abonnerProjet, desabonnerProjet } from "@/tempsreel/prise"
import Avatar from "@/composants/ui/Avatar.vue"
import Selection from "@/composants/ui/Selection.vue"

const proprietes = defineProps<{
  groupe: string
}>()

const identifiant = computed(() => proprietes.groupe)
const { data: projets } = utiliserProjets(identifiant)

watch(
  () => projets.value,
  (nouveaux, anciens) => {
    for (const projet of anciens ?? []) desabonnerProjet(projet.id)
    for (const projet of nouveaux ?? []) abonnerProjet(projet.id)
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  for (const projet of projets.value ?? []) desabonnerProjet(projet.id)
})

type Periode = "7j" | "30j" | "90j" | "12m" | "libre"

const periodes: { valeur: Periode; libelle: string; jours: number; granularite: ParametresStatistiques["granularite"] }[] = [
  { valeur: "7j", libelle: "7 jours", jours: 6, granularite: "jour" },
  { valeur: "30j", libelle: "30 jours", jours: 29, granularite: "jour" },
  { valeur: "90j", libelle: "3 mois", jours: 89, granularite: "semaine" },
  { valeur: "12m", libelle: "12 mois", jours: 364, granularite: "mois" },
  { valeur: "libre", libelle: "Personnalisé", jours: 29, granularite: "jour" },
]

const periode = ref<Periode>("30j")
const parametres = ref<ParametresStatistiques>({
  debut: dayjs().subtract(29, "day").format("YYYY-MM-DD"),
  fin: dayjs().format("YYYY-MM-DD"),
  granularite: "jour",
  projet: "",
})

watch(periode, (valeur) => {
  const configuration = periodes.find((element) => element.valeur === valeur)
  if (!configuration || valeur === "libre") return
  parametres.value = {
    ...parametres.value,
    debut: dayjs().subtract(configuration.jours, "day").format("YYYY-MM-DD"),
    fin: dayjs().format("YYYY-MM-DD"),
    granularite: configuration.granularite,
  }
})

const { data: statistiques, isLoading: chargement } = utiliserStatistiques(identifiant, parametres)

const maximumPoints = computed(() =>
  Math.max(1, ...(statistiques.value?.classement ?? []).map((ligne) => ligne.points)),
)
const maximumSerie = computed(() =>
  Math.max(1, ...(statistiques.value?.serie ?? []).map((point) => point.points)),
)

function libellePeriode(valeur: string): string {
  const moment = dayjs(valeur)
  if (parametres.value.granularite === "mois") return moment.format("MMM YYYY")
  if (parametres.value.granularite === "semaine") return "sem. " + moment.format("DD MMM")
  return moment.format("DD MMM")
}

const tuiles = computed(() => [
  { libelle: "Points réalisés", valeur: statistiques.value?.totaux.points ?? 0 },
  { libelle: "Tâches terminées", valeur: statistiques.value?.totaux.terminees ?? 0 },
  { libelle: "Tâches en retard", valeur: statistiques.value?.totaux.enretard ?? 0 },
  {
    libelle: "Total du périmètre",
    valeur: statistiques.value?.totaux.total ?? 0,
    complement: `${statistiques.value?.totaux.totalpoints ?? 0} pts`,
  },
])

const classeChampDate =
  "h-9 rounded-lg border border-neutral-300 bg-white px-2.5 text-sm shadow-sm focus:border-neutral-500 focus:outline-none focus:ring-2 focus:ring-neutral-900/10 dark:border-neutral-700 dark:bg-neutral-900 dark:focus:border-neutral-400"
</script>

<template>
  <div class="space-y-6">
    <section class="rounded-xl border border-neutral-200 bg-white p-5 dark:border-neutral-800 dark:bg-neutral-900">
      <div class="flex flex-wrap items-end gap-3">
        <div class="flex rounded-lg border border-neutral-300 p-0.5 dark:border-neutral-700">
          <button
            v-for="option in periodes"
            :key="option.valeur"
            class="rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
            :class="
              periode === option.valeur
                ? 'bg-neutral-900 text-white dark:bg-neutral-100 dark:text-neutral-900'
                : 'text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100'
            "
            @click="periode = option.valeur"
          >
            {{ option.libelle }}
          </button>
        </div>
        <template v-if="periode === 'libre'">
          <label class="space-y-1">
            <span class="block text-xs text-neutral-500">Du</span>
            <input v-model="parametres.debut" type="date" :class="classeChampDate" />
          </label>
          <label class="space-y-1">
            <span class="block text-xs text-neutral-500">Au</span>
            <input v-model="parametres.fin" type="date" :class="classeChampDate" />
          </label>
        </template>
        <div class="w-full sm:w-36">
          <Selection v-model="parametres.granularite" etiquette="Granularité">
            <option value="jour">Par jour</option>
            <option value="semaine">Par semaine</option>
            <option value="mois">Par mois</option>
          </Selection>
        </div>
        <div class="w-full sm:w-48">
          <Selection v-model="parametres.projet" etiquette="Projet">
            <option value="">Tous les projets</option>
            <option v-for="projet in projets" :key="projet.id" :value="projet.id">{{ projet.nom }}</option>
          </Selection>
        </div>
      </div>
    </section>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="tuile in tuiles"
        :key="tuile.libelle"
        class="rounded-xl border border-neutral-200 bg-white p-4 dark:border-neutral-800 dark:bg-neutral-900"
      >
        <p class="text-sm text-neutral-500">{{ tuile.libelle }}</p>
        <p class="mt-1 text-2xl font-bold">
          {{ tuile.valeur }}
          <span v-if="tuile.complement" class="text-sm font-normal text-neutral-400">· {{ tuile.complement }}</span>
        </p>
      </div>
    </div>

    <section class="rounded-xl border border-neutral-200 bg-white p-5 dark:border-neutral-800 dark:bg-neutral-900">
      <h3 class="font-semibold">Points réalisés dans le temps</h3>
      <p class="mt-0.5 text-xs text-neutral-500">Tâches arrivées dans la dernière colonne de leur tableau.</p>
      <p v-if="chargement" class="py-10 text-center text-sm text-neutral-500">Chargement…</p>
      <div v-else-if="statistiques?.serie.length" class="mt-5 flex h-44 items-end gap-1.5 overflow-x-auto pb-8">
        <div
          v-for="point in statistiques.serie"
          :key="point.periode"
          class="group relative flex h-full min-w-9 flex-1 flex-col items-center justify-end"
        >
          <span class="mb-1 text-[11px] font-semibold text-neutral-500 opacity-0 transition-opacity group-hover:opacity-100">
            {{ point.points }} pts · {{ point.taches }} t.
          </span>
          <div
            class="w-full max-w-12 rounded-t-md bg-neutral-800 transition-colors group-hover:bg-neutral-600 dark:bg-neutral-300 dark:group-hover:bg-neutral-100"
            :style="{ height: Math.max(4, Math.round((point.points / maximumSerie) * 100)) + '%' }"
          ></div>
          <span class="absolute -bottom-7 whitespace-nowrap text-[10px] text-neutral-400">
            {{ libellePeriode(point.periode) }}
          </span>
        </div>
      </div>
      <p v-else class="rounded-lg border border-dashed border-neutral-300 p-8 text-center text-sm text-neutral-500 dark:border-neutral-700">
        Aucune tâche terminée sur la période.
      </p>
    </section>

    <section class="rounded-xl border border-neutral-200 bg-white p-5 dark:border-neutral-800 dark:bg-neutral-900">
      <h3 class="font-semibold">Classement des contributeurs</h3>
      <p class="mt-0.5 text-xs text-neutral-500">
        Classés par points de tâches terminées, puis par nombre de tâches. Une tâche sans affectation est créditée à sa créatrice ou son créateur.
      </p>
      <ol v-if="statistiques?.classement.length" class="mt-4 space-y-2">
        <li
          v-for="(ligne, indice) in statistiques.classement"
          :key="ligne.utilisateur"
          class="flex items-center gap-3 rounded-lg border border-neutral-200 bg-neutral-50 px-3 py-2.5 dark:border-neutral-700 dark:bg-neutral-800/50"
        >
          <span
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold"
            :class="
              indice === 0
                ? 'bg-neutral-900 text-white dark:bg-neutral-100 dark:text-neutral-900'
                : 'bg-neutral-200 text-neutral-600 dark:bg-neutral-700 dark:text-neutral-300'
            "
          >
            {{ indice + 1 }}
          </span>
          <Avatar :nom="ligne.nom" :prenom="ligne.prenom" />
          <div class="min-w-0 flex-1">
            <div class="flex items-baseline justify-between gap-2">
              <p class="truncate text-sm font-medium">{{ ligne.prenom }} {{ ligne.nom }}</p>
              <p class="shrink-0 text-sm font-bold">{{ ligne.points }} pts</p>
            </div>
            <div class="mt-1 h-1.5 overflow-hidden rounded-full bg-neutral-200 dark:bg-neutral-700">
              <div
                class="h-full rounded-full bg-neutral-800 dark:bg-neutral-300"
                :style="{ width: Math.round((ligne.points / maximumPoints) * 100) + '%' }"
              ></div>
            </div>
            <p class="mt-1 text-xs text-neutral-500">
              {{ ligne.taches }} terminée{{ ligne.taches > 1 ? "s" : "" }} ·
              {{ ligne.creees }} créée{{ ligne.creees > 1 ? "s" : "" }} ·
              {{ ligne.commentaires }} commentaire{{ ligne.commentaires > 1 ? "s" : "" }}
            </p>
          </div>
        </li>
      </ol>
      <p v-else class="mt-4 rounded-lg border border-dashed border-neutral-300 p-8 text-center text-sm text-neutral-500 dark:border-neutral-700">
        Aucune contribution sur la période.
      </p>
    </section>
  </div>
</template>
