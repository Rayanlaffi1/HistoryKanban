<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue"

export interface ElementMenu {
  id: string
  libelle: string
  danger?: boolean
}

const proprietes = defineProps<{
  ouvert: boolean
  x: number
  y: number
  elements: ElementMenu[]
}>()

const emissions = defineEmits<{
  choisir: [string]
  fermer: []
}>()

const panneau = ref<HTMLElement | null>(null)
const position = ref({ x: 0, y: 0 })

watch(
  () => [proprietes.ouvert, proprietes.x, proprietes.y] as const,
  async ([ouvert]) => {
    if (!ouvert) return
    position.value = { x: proprietes.x, y: proprietes.y }
    await nextTick()
    const boite = panneau.value?.getBoundingClientRect()
    if (!boite) return
    position.value = {
      x: Math.min(proprietes.x, window.innerWidth - boite.width - 8),
      y: Math.min(proprietes.y, window.innerHeight - boite.height - 8),
    }
  },
)

function surPointeur(evenement: PointerEvent) {
  if (!proprietes.ouvert) return
  if (panneau.value && panneau.value.contains(evenement.target as Node)) return
  emissions("fermer")
}

function surTouche(evenement: KeyboardEvent) {
  if (proprietes.ouvert && evenement.key === "Escape") {
    emissions("fermer")
  }
}

function surDefilement() {
  if (proprietes.ouvert) {
    emissions("fermer")
  }
}

onMounted(() => {
  window.addEventListener("pointerdown", surPointeur, true)
  window.addEventListener("keydown", surTouche)
  window.addEventListener("scroll", surDefilement, true)
  window.addEventListener("resize", surDefilement)
})

onBeforeUnmount(() => {
  window.removeEventListener("pointerdown", surPointeur, true)
  window.removeEventListener("keydown", surTouche)
  window.removeEventListener("scroll", surDefilement, true)
  window.removeEventListener("resize", surDefilement)
})

function choisir(id: string) {
  emissions("choisir", id)
  emissions("fermer")
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="ouvert"
      ref="panneau"
      class="fixed z-50 min-w-48 overflow-hidden rounded-lg border border-neutral-200 bg-white py-1 shadow-xl dark:border-neutral-700 dark:bg-neutral-900"
      :style="{ left: position.x + 'px', top: position.y + 'px' }"
      @contextmenu.prevent
    >
      <button
        v-for="element in elements"
        :key="element.id"
        class="block w-full px-3 py-2 text-left text-sm transition-colors"
        :class="
          element.danger
            ? 'text-red-700 hover:bg-red-50 dark:text-red-500 dark:hover:bg-red-950/40'
            : 'text-neutral-800 hover:bg-neutral-100 dark:text-neutral-200 dark:hover:bg-neutral-800'
        "
        @click="choisir(element.id)"
      >
        {{ element.libelle }}
      </button>
    </div>
  </Teleport>
</template>
