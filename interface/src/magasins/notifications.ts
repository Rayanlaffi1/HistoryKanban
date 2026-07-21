import { defineStore } from "pinia"
import { ref } from "vue"

export interface MessageEphemere {
  id: number
  titre: string
  detail: string
}

export const utiliserMagasinNotifications = defineStore("notifications", () => {
  const messages = ref<MessageEphemere[]>([])
  let compteur = 0

  function annoncer(titre: string, detail: string) {
    const id = ++compteur
    messages.value.push({ id, titre, detail })
    setTimeout(() => {
      messages.value = messages.value.filter((message) => message.id !== id)
    }, 5000)
  }

  function fermer(id: number) {
    messages.value = messages.value.filter((message) => message.id !== id)
  }

  return { messages, annoncer, fermer }
})
