// Sans variables Vite (build Docker), les URL publiques sont derivees de
// l'adresse consultee : l'application fonctionne ainsi indifferemment via
// historykanban.localhost ou historykanban.<IP>.sslip.io sans reconstruction.
const hote = window.location.host
const securise = window.location.protocol === "https:"

export const urlAPI = import.meta.env.VITE_URLAPI || window.location.origin
export const urlWS = import.meta.env.VITE_URLWS || `${securise ? "wss" : "ws"}://${hote}`
export const urlKeycloak = import.meta.env.VITE_KEYCLOAKURL || `https://auth.${hote}`
