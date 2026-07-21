<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from "vue"
import { useEditor, EditorContent } from "@tiptap/vue-3"
import StarterKit from "@tiptap/starter-kit"
import Placeholder from "@tiptap/extension-placeholder"
import { ImageRedimensionnable } from "@/composants/ui/imageredimensionnable"

const proprietes = defineProps<{
  etiquette?: string
  indication?: string
  erreur?: string
  compact?: boolean
  televerser?: (fichier: File) => Promise<string>
}>()

const modele = defineModel<string>({ default: "" })
const champImage = ref<HTMLInputElement | null>(null)
const envoiEnCours = ref(false)

function deposerFichiers(fichiers: File[], position: number | null) {
  if (!proprietes.televerser) return
  for (const fichier of fichiers) {
    envoiEnCours.value = true
    proprietes
      .televerser(fichier)
      .then((url) => {
        const instance = editeur.value
        if (!instance) return
        if (position !== null) {
          const noeud = instance.state.schema.nodes.image.create({ src: url, alt: fichier.name })
          instance.view.dispatch(instance.state.tr.insert(Math.min(position, instance.state.doc.content.size), noeud))
        } else {
          instance.chain().focus().setImage({ src: url, alt: fichier.name }).run()
        }
      })
      .finally(() => {
        envoiEnCours.value = false
      })
  }
}

const editeur = useEditor({
  content: modele.value,
  extensions: [
    StarterKit,
    ImageRedimensionnable,
    Placeholder.configure({ placeholder: proprietes.indication ?? "" }),
  ],
  editorProps: {
    attributes: {
      class:
        "texteriche focus:outline-none px-3 py-2 text-sm " +
        (proprietes.compact ? "min-h-16" : "min-h-28"),
    },
    handleDrop: (vue, evenement, _tranche, deplace) => {
      if (deplace || !proprietes.televerser) return false
      const fichiers = Array.from(evenement.dataTransfer?.files ?? []).filter((fichier) =>
        fichier.type.startsWith("image/"),
      )
      if (!fichiers.length) return false
      evenement.preventDefault()
      const cible = vue.posAtCoords({ left: evenement.clientX, top: evenement.clientY })
      deposerFichiers(fichiers, cible?.pos ?? null)
      return true
    },
    handlePaste: (_vue, evenement) => {
      if (!proprietes.televerser) return false
      const fichiers = Array.from(evenement.clipboardData?.files ?? []).filter((fichier) =>
        fichier.type.startsWith("image/"),
      )
      if (!fichiers.length) return false
      evenement.preventDefault()
      deposerFichiers(fichiers, null)
      return true
    },
  },
  onUpdate: ({ editor }) => {
    modele.value = editor.isEmpty ? "" : editor.getHTML()
  },
})

watch(modele, (valeur) => {
  const instance = editeur.value
  if (!instance) return
  const actuel = instance.isEmpty ? "" : instance.getHTML()
  if (valeur !== actuel) {
    instance.commands.setContent(valeur || "")
  }
})

onBeforeUnmount(() => {
  editeur.value?.destroy()
})

async function insererImage(evenement: Event) {
  const fichier = (evenement.target as HTMLInputElement).files?.[0]
  ;(evenement.target as HTMLInputElement).value = ""
  if (!fichier || !proprietes.televerser || !editeur.value) return
  envoiEnCours.value = true
  try {
    const url = await proprietes.televerser(fichier)
    editeur.value.chain().focus().setImage({ src: url, alt: fichier.name }).run()
  } finally {
    envoiEnCours.value = false
  }
}

const classeBouton =
  "rounded px-1.5 py-0.5 text-xs font-semibold text-neutral-600 hover:bg-neutral-200 dark:text-neutral-300 dark:hover:bg-neutral-700"
const classeActive = "bg-neutral-300 text-neutral-900 dark:bg-neutral-600 dark:text-neutral-100"
</script>

<template>
  <div class="block space-y-1.5">
    <span v-if="etiquette" class="text-sm font-medium text-neutral-700 dark:text-neutral-300">{{ etiquette }}</span>
    <div
      class="rounded-lg border bg-white dark:bg-neutral-900"
      :class="erreur ? 'border-red-500 dark:border-red-700' : 'border-neutral-300 dark:border-neutral-700'"
    >
      <div
        v-if="editeur"
        class="flex flex-wrap items-center gap-1 border-b border-neutral-200 px-2 py-1.5 dark:border-neutral-700"
      >
        <button type="button" :class="[classeBouton, editeur.isActive('bold') && classeActive]" title="Gras" @click="editeur.chain().focus().toggleBold().run()">G</button>
        <button type="button" :class="[classeBouton, editeur.isActive('italic') && classeActive]" title="Italique" @click="editeur.chain().focus().toggleItalic().run()"><span class="italic">I</span></button>
        <button type="button" :class="[classeBouton, editeur.isActive('strike') && classeActive]" title="Barré" @click="editeur.chain().focus().toggleStrike().run()"><span class="line-through">S</span></button>
        <span class="mx-1 h-4 w-px bg-neutral-300 dark:bg-neutral-600"></span>
        <button type="button" :class="[classeBouton, editeur.isActive('heading', { level: 2 }) && classeActive]" title="Titre" @click="editeur.chain().focus().toggleHeading({ level: 2 }).run()">T</button>
        <button type="button" :class="[classeBouton, editeur.isActive('bulletList') && classeActive]" title="Liste à puces" @click="editeur.chain().focus().toggleBulletList().run()">•</button>
        <button type="button" :class="[classeBouton, editeur.isActive('orderedList') && classeActive]" title="Liste numérotée" @click="editeur.chain().focus().toggleOrderedList().run()">1.</button>
        <button type="button" :class="[classeBouton, editeur.isActive('blockquote') && classeActive]" title="Citation" @click="editeur.chain().focus().toggleBlockquote().run()">❝</button>
        <template v-if="televerser">
          <span class="mx-1 h-4 w-px bg-neutral-300 dark:bg-neutral-600"></span>
          <button type="button" :class="classeBouton" :disabled="envoiEnCours" title="Insérer une image" @click="champImage?.click()">
            {{ envoiEnCours ? "Envoi…" : "Image" }}
          </button>
          <input ref="champImage" type="file" accept="image/*" class="hidden" @change="insererImage" />
        </template>
        <span class="mx-1 h-4 w-px bg-neutral-300 dark:bg-neutral-600"></span>
        <button type="button" :class="classeBouton" title="Annuler" @click="editeur.chain().focus().undo().run()">↶</button>
        <button type="button" :class="classeBouton" title="Rétablir" @click="editeur.chain().focus().redo().run()">↷</button>
      </div>
      <EditorContent :editor="editeur" />
    </div>
    <span v-if="erreur" class="block text-xs text-red-700 dark:text-red-500">{{ erreur }}</span>
  </div>
</template>
