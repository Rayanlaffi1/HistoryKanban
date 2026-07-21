import axios from "axios"
import { jeton, seDeconnecter } from "@/securite/keycloak"

export const client = axios.create({
  baseURL: (import.meta.env.VITE_URLAPI ?? "https://api.historykanban.localhost") + "/api",
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
    }
    return Promise.reject(erreur)
  },
)
