<script setup lang="ts">
import FullGame from '@/components/FullGame.vue'
import { useGameStore } from '@/stores/game.ts'
import { ref, watch } from 'vue'
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
</script>

<template>
  <div class="header">
    <!-- eslint-disable-next-line vue/no-parsing-error -->
    <RouterLink to="/matches" class="back-button"> < back</RouterLink>
    <div class="title">
      <h1 class="blue-name">
        {{ gameStore.gameMetadata?.blueBotName }}
      </h1>
      <h1>VS</h1>
      <h1 class="red-name">
        {{ gameStore.gameMetadata?.redBotName }}
      </h1>
    </div>
  </div>
  <FullGame :game="gameStore.game" :metadata="gameStore.gameMetadata" />
</template>

<style lang="css" scoped>
.header {
  position: relative;
  width: 100%;
  display: flex;
  flex-direction: row;
}

.back-button {
  margin: 16px;
  position: absolute;
  color: black;
  text-decoration: none;
}

.back-button:hover {
  text-decoration: underline;
}

.title {
  padding-top: 16px;
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
}

.blue-name,
.red-name {
  display: flex;
  flex-grow: 1;
  flex-basis: 50%;
  margin: 16px;
}

.blue-name {
  justify-content: end;
}

.red-name {
  justify-content: start;
}
</style>
