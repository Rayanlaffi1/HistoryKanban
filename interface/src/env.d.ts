/// <reference types="vite/client" />

declare module "*.vue" {
  import type { DefineComponent } from "vue"
  const composant: DefineComponent<{}, {}, any>
  export default composant
}

interface ImportMetaEnv {
  readonly VITE_URLAPI: string
  readonly VITE_URLWS: string
  readonly VITE_KEYCLOAKURL: string
  readonly VITE_KEYCLOAKROYAUME: string
  readonly VITE_KEYCLOAKCLIENT: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
