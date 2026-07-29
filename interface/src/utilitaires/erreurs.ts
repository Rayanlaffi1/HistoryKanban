import { AxiosError } from "axios"
import { urlAPI } from "@/configuration"

const baseAPI = urlAPI

export function extraireErreur(erreur: unknown): string {
  if (erreur instanceof AxiosError) {
    if (!erreur.response) {
      return `API injoignable. Ouvrez ${baseAPI}/sante dans un onglet, acceptez le certificat, puis réessayez.`
    }
    const message = (erreur.response.data as { erreur?: string } | undefined)?.erreur
    if (message) {
      return `${message} (code ${erreur.response.status})`
    }
    return `Le serveur a renvoyé une erreur ${erreur.response.status}.`
  }
  if (erreur instanceof Error) {
    return erreur.message
  }
  return "Erreur inattendue."
}
