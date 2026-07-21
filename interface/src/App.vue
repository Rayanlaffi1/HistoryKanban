<script setup lang="ts">
import EnteteApplication from "@/composants/EnteteApplication.vue"
import MessagesEphemeres from "@/composants/MessagesEphemeres.vue"
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <EnteteApplication />
    <main class="flex-1">
      <RouterView />
    </main>
    <MessagesEphemeres />
  </div>
</template>
