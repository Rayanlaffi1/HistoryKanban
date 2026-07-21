<script setup lang="ts">
import { computed } from "vue"
import { classes } from "@/utilitaires/classes"

const proprietes = withDefaults(
  defineProps<{
    variante?: "primaire" | "secondaire" | "fantome" | "danger"
    taille?: "petite" | "normale"
    type?: "button" | "submit"
    desactive?: boolean
  }>(),
  { variante: "primaire", taille: "normale", type: "button", desactive: false },
)

const apparence = computed(() =>
  classes(
    "inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-neutral-400 disabled:pointer-events-none disabled:opacity-50",
    proprietes.taille === "petite" ? "h-8 px-3 text-sm" : "h-10 px-4 text-sm",
    {
      primaire:
        "bg-neutral-900 text-white hover:bg-neutral-700 dark:bg-neutral-100 dark:text-neutral-900 dark:hover:bg-neutral-300",
      secondaire:
        "border border-neutral-300 bg-white text-neutral-900 hover:bg-neutral-100 dark:border-neutral-700 dark:bg-neutral-900 dark:text-neutral-100 dark:hover:bg-neutral-800",
      fantome:
        "text-neutral-700 hover:bg-neutral-200 dark:text-neutral-300 dark:hover:bg-neutral-800",
      danger:
        "bg-red-700 text-white hover:bg-red-800 dark:bg-red-800 dark:hover:bg-red-700",
    }[proprietes.variante],
  ),
)
</script>

<template>
  <button :type="proprietes.type" :class="apparence" :disabled="proprietes.desactive">
    <slot />
  </button>
</template>
