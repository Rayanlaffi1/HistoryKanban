import { marked } from "marked"
import { assainir } from "@/utilitaires/html"

// Le HTML brut éventuellement présent (anciens commentaires) traverse marked
// tel quel puis passe par DOMPurify, ce qui garde l'affichage rétrocompatible.
export function rendreMarkdown(texte: string): string {
  return assainir(marked.parse(texte, { gfm: true, breaks: true, async: false }))
}
