import Image from "@tiptap/extension-image"

export const ImageRedimensionnable = Image.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      largeur: {
        default: null,
        parseHTML: (element: HTMLElement) => element.getAttribute("width"),
        renderHTML: (attributs: Record<string, unknown>) =>
          attributs.largeur ? { width: String(attributs.largeur) } : {},
      },
    }
  },

  addNodeView() {
    return ({ node, editor, getPos }) => {
      let noeudActuel = node
      const conteneur = document.createElement("div")
      conteneur.className = "imageconteneur"
      const image = document.createElement("img")
      image.src = node.attrs.src
      if (node.attrs.alt) image.alt = node.attrs.alt
      if (node.attrs.largeur) image.setAttribute("width", String(node.attrs.largeur))
      conteneur.appendChild(image)

      if (editor.isEditable) {
        const poignee = document.createElement("span")
        poignee.className = "imagepoignee"
        conteneur.appendChild(poignee)

        let departX = 0
        let departLargeur = 0

        const surMouvement = (evenement: MouseEvent) => {
          const largeur = Math.round(Math.max(60, departLargeur + evenement.clientX - departX))
          image.setAttribute("width", String(largeur))
        }

        const surFin = () => {
          window.removeEventListener("mousemove", surMouvement)
          window.removeEventListener("mouseup", surFin)
          const largeur = image.getAttribute("width")
          const position = typeof getPos === "function" ? getPos() : null
          if (position !== null && position !== undefined) {
            const transaction = editor.state.tr.setNodeMarkup(position, undefined, {
              ...noeudActuel.attrs,
              largeur,
            })
            editor.view.dispatch(transaction)
          }
        }

        poignee.addEventListener("mousedown", (evenement) => {
          evenement.preventDefault()
          evenement.stopPropagation()
          departX = evenement.clientX
          departLargeur = image.getBoundingClientRect().width
          window.addEventListener("mousemove", surMouvement)
          window.addEventListener("mouseup", surFin)
        })
      }

      return {
        dom: conteneur,
        update: (nouveau) => {
          if (nouveau.type.name !== noeudActuel.type.name) return false
          noeudActuel = nouveau
          image.src = nouveau.attrs.src
          if (nouveau.attrs.largeur) {
            image.setAttribute("width", String(nouveau.attrs.largeur))
          } else {
            image.removeAttribute("width")
          }
          return true
        },
      }
    }
  },
})
