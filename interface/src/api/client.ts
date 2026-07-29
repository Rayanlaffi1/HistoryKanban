import axios from "axios"
import { urlAPI } from "@/configuration"
import { jeton, seDeconnecter } from "@/securite/keycloak"
import { extraireErreur } from "@/utilitaires/erreurs"
import { utiliserMagasinNotifications } from "@/magasins/notifications"

export const client = axios.create({
  baseURL: urlAPI + "/api",
})

client.interceptors.request.use((requete) => {
  requete.headers.Authorization = `Bearer ${jeton()}`
  return requete
})

client.interceptors.response.use(
  (reponse) => reponse,
  (erreur) => {
    if (erreur.response?.status === 401 && erreur.response?.data?.code === "session.remplacee") {
      seDeconnecter()
      return Promise.reject(erreur)
    }
    if (erreur.config?.method !== "get") {
      utiliserMagasinNotifications().annoncer("Action impossible", extraireErreur(erreur))
    }
    return Promise.reject(erreur)
  },
)
