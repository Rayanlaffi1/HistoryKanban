import { ref } from "vue"
import type { QueryClient } from "@tanstack/vue-query"
import { useLocalStorage } from "@vueuse/core"
import { client } from "@/api/client"
import { urlWS } from "@/configuration"
import { jeton, seDeconnecter } from "@/securite/keycloak"
import { utiliserMagasinNotifications } from "@/magasins/notifications"
import { typesNotifications } from "@/api/types"

let prise: WebSocket | null = null
let actif = false
const abonnements = new Set<string>()
const projetsStatistiquesEnAttente = new Set<string>()
let minuterieStatistiques: ReturnType<typeof setTimeout> | null = null

export const notifierDeplacement = useLocalStorage("historykanban.notifdeplacement", true)
export const connectes = ref<string[]>([])

async function chargerPresence() {
  try {
    const reponse = await client.get("/presence")
    connectes.value = reponse.data.connectes ?? []
  } catch {
    connectes.value = []
  }
}

interface MessageServeur {
  type: string
  projet?: string
  acteurnom?: string
  titre?: string
  utilisateur?: string
  etat?: string
  donnees?: any
}

// Invalidation différée et ciblée des statistiques : une rafale d'événements de
// tâches ne déclenche qu'un seul recalcul, et seules les vues concernées (tous
// projets ou le projet touché) sont rafraîchies.
function planifierInvalidationStatistiques(clientRequetes: QueryClient, projet: string) {
  projetsStatistiquesEnAttente.add(projet)
  if (minuterieStatistiques) return
  minuterieStatistiques = setTimeout(() => {
    const projets = new Set(projetsStatistiquesEnAttente)
    projetsStatistiquesEnAttente.clear()
    minuterieStatistiques = null
    clientRequetes.invalidateQueries({
      predicate: (requete) => {
        if (requete.queryKey[0] !== "statistiques") return false
        const parametres = requete.queryKey[2] as { projet?: string } | undefined
        return !parametres?.projet || projets.has(parametres.projet)
      },
    })
  }, 2000)
}

export function demarrerTempsReel(clientRequetes: QueryClient) {
  if (actif) return
  actif = true
  const magasin = utiliserMagasinNotifications()

  function connecter() {
    const base = urlWS
    prise = new WebSocket(`${base}/ws`, ["historykanban", `bearer.${jeton()}`])
    prise.onopen = () => {
      for (const projet of abonnements) {
        prise?.send(JSON.stringify({ action: "abonner", projet }))
      }
      chargerPresence()
    }
    prise.onmessage = (evenement) => {
      let message: MessageServeur
      try {
        message = JSON.parse(evenement.data)
      } catch {
        return
      }
      if (message.type === "session.remplacee") {
        magasin.annoncer("Session terminée", "Une connexion a été ouverte sur un autre appareil.")
        setTimeout(() => seDeconnecter(), 1500)
        return
      }
      if (message.type === "presence" && message.utilisateur) {
        const sans = connectes.value.filter((identifiant) => identifiant !== message.utilisateur)
        connectes.value = message.etat === "enligne" ? [...sans, message.utilisateur] : sans
        return
      }
      if (message.type === "notification") {
        clientRequetes.invalidateQueries({ queryKey: ["notifications"] })
        const contenu = message.donnees?.contenu ?? {}
        magasin.annoncer(
          typesNotifications[message.donnees?.type] ?? "Notification",
          [contenu.titre, contenu.acteurnom ? `par ${contenu.acteurnom}` : ""].filter(Boolean).join(" — "),
        )
        return
      }
      if (message.projet) {
        if (message.type.startsWith("tache.")) {
          clientRequetes.invalidateQueries({ queryKey: ["taches", message.projet] })
        } else {
          clientRequetes.invalidateQueries({ queryKey: ["projet", message.projet] })
        }
        if (message.type === "tache.commentee" && message.donnees?.tache) {
          clientRequetes.invalidateQueries({ queryKey: ["commentaires", message.donnees.tache] })
        }
        const cibleActivite = message.donnees?.tache ?? message.donnees?.id
        if (message.type.startsWith("tache.")) {
          if (cibleActivite) {
            clientRequetes.invalidateQueries({ queryKey: ["activites", cibleActivite] })
          }
          planifierInvalidationStatistiques(clientRequetes, message.projet)
        }
        if (message.type === "tache.deplacee" && notifierDeplacement.value) {
          magasin.annoncer(
            "Tâche déplacée",
            [message.titre, message.acteurnom ? `par ${message.acteurnom}` : ""].filter(Boolean).join(" — "),
          )
        }
      }
    }
    prise.onclose = () => {
      prise = null
      setTimeout(connecter, 3000)
    }
  }

  connecter()
}

export function abonnerProjet(projet: string) {
  abonnements.add(projet)
  if (prise?.readyState === WebSocket.OPEN) {
    prise.send(JSON.stringify({ action: "abonner", projet }))
  }
}

export function desabonnerProjet(projet: string) {
  abonnements.delete(projet)
  if (prise?.readyState === WebSocket.OPEN) {
    prise.send(JSON.stringify({ action: "desabonner", projet }))
  }
}
