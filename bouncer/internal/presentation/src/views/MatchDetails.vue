<script setup lang="ts">
import PageTitle from '@/components/PageTitle.vue'
import FullGame from '@/components/FullGame.vue'
import { useGameStore } from '@/stores/game.ts'
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const gameStore = useGameStore()

const route = useRoute()

const gameId = ref('')

watch(
  () => route.params.id,
  (newId, oldId) => {
    if (Array.isArray(newId)) {
      return
    }
    if (newId !== oldId || !gameId.value.length) {
      gameId.value = newId
    }
    if (gameStore.gameMetadata?.id !== gameId.value) {
      gameStore.resetGame()
      gameStore.loadGame(gameId.value)
    }
  },
  { immediate: true },
)

const subtitle = computed(
  () => `${gameStore.gameMetadata?.blueBotName} VS ${gameStore.gameMetadata?.redBotName}`,
)
</script>

<template>
  <PageTitle title="Matches" :subtitles="[subtitle]" />
  <FullGame :game="gameStore.game" :metadata="gameStore.gameMetadata" />
</template>

<style lang="css" scoped></style>
