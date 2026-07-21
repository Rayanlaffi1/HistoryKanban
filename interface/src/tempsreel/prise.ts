import type { QueryClient } from "@tanstack/vue-query"
import { jeton, seDeconnecter } from "@/securite/keycloak"
import { utiliserMagasinNotifications } from "@/magasins/notifications"
import { typesNotifications } from "@/api/types"

let prise: WebSocket | null = null
let actif = false
const abonnements = new Set<string>()

interface MessageServeur {
  type: string
  projet?: string
  donnees?: any
}

export function demarrerTempsReel(clientRequetes: QueryClient) {
  if (actif) return
  actif = true
  const magasin = utiliserMagasinNotifications()

  function connecter() {
    const base = import.meta.env.VITE_URLWS ?? "wss://api.historykanban.localhost"
    prise = new WebSocket(`${base}/ws?jeton=${jeton()}`)
    prise.onopen = () => {
      for (const projet of abonnements) {
        prise?.send(JSON.stringify({ action: "abonner", projet }))
      }
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
