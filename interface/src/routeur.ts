import { createRouter, createWebHistory } from "vue-router"

export const routeur = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "accueil",
      component: () => import("@/pages/PageAccueil.vue"),
    },
    {
      path: "/groupes/:id",
      name: "groupe",
      component: () => import("@/pages/PageGroupe.vue"),
    },
    {
      path: "/projets/:id",
      name: "projet",
      component: () => import("@/pages/PageProjet.vue"),
    },
    {
      path: "/parametres",
      name: "parametres",
      component: () => import("@/pages/PageParametres.vue"),
    },
    {
      path: "/:inconnu(.*)*",
      redirect: "/",
    },
  ],
})
