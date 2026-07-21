import { createApp } from "vue"
import { createPinia } from "pinia"
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query"
import App from "@/App.vue"
import { routeur } from "@/routeur"
import { demarrerAuthentification } from "@/securite/keycloak"
import { demarrerTempsReel } from "@/tempsreel/prise"
import "@/style.css"

async function demarrer() {
  const connecte = await demarrerAuthentification()
  if (!connecte) return
  const clientRequetes = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 10000,
        retry: 1,
      },
    },
  })
  const application = createApp(App)
  application.use(createPinia())
  application.use(routeur)
  application.use(VueQueryPlugin, { queryClient: clientRequetes })
  application.mount("#application")
  demarrerTempsReel(clientRequetes)
}

demarrer()
