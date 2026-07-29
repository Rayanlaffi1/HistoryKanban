import Keycloak from "keycloak-js"
import { urlKeycloak } from "@/configuration"

const keycloak = new Keycloak({
  url: urlKeycloak,
  realm: import.meta.env.VITE_KEYCLOAKREALM ?? "historykanban",
  clientId: import.meta.env.VITE_KEYCLOAKCLIENT ?? "interface",
})

export async function demarrerAuthentification(): Promise<boolean> {
  const connecte = await keycloak.init({
    onLoad: "login-required",
    pkceMethod: "S256",
    checkLoginIframe: false,
  })
  if (connecte) {
    setInterval(async () => {
      try {
        await keycloak.updateToken(60)
      } catch {
        keycloak.login()
      }
    }, 30000)
  }
  return connecte
}

export function jeton(): string {
  return keycloak.token ?? ""
}

export function seDeconnecter(): void {
  keycloak.logout({ redirectUri: window.location.origin })
}

export default keycloak
