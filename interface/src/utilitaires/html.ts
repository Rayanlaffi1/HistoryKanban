import DOMPurify from "dompurify"

export function assainir(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: [
      "p", "br", "strong", "em", "s", "u", "h1", "h2", "h3",
      "ul", "ol", "li", "blockquote", "code", "pre", "a", "img",
    ],
    ALLOWED_ATTR: ["href", "src", "alt", "title", "target", "rel", "width", "height"],
  })
}

export function texteBrut(html: string): string {
  const zone = document.createElement("div")
  zone.innerHTML = assainir(html)
  return (zone.textContent ?? "").trim()
}

export function contenuVide(html: string): boolean {
  return texteBrut(html) === "" && !html.includes("<img")
}
