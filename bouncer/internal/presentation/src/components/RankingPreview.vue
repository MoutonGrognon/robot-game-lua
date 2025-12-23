<script setup lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'

/* import all the icons in Free Solid, Free Regular, and Brands styles */
import { fas } from '@fortawesome/free-solid-svg-icons'
// import { far } from '@fortawesome/free-regular-svg-icons'
// import { fab } from '@fortawesome/free-brands-svg-icons'

library.add(fas)

// TODO: type
defineProps<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  preview: any
  index: number
}>()
</script>

<template>
  <div class="wrapper">
    <div class="count">{{ preview.rank }} .</div>
    <div class="elo">
      {{ preview.elo }}
    </div>
    <div class="name">
      <div>
        {{ preview.botName }}
      </div>
      <div class="win-streak" v-if="preview.streak > 2">
        <FontAwesomeIcon :icon="['fas', 'fire']" />
        <div class="win-streak-count">×{{ preview.streak }}</div>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.wrapper {
  display: flex;
  flex-direction: row;
  padding: 8px;
  column-gap: 32px;
  width: calc(100% - 16px);

  color: black;
  text-decoration: none;
}

.count {
  width: 36px;
  text-align: right;
}

.elo {
  font-weight: 700;
  width: 48px;
  text-align: right;
}

.win-streak {
  display: flex;
  align-content: center;
  justify-content: center;
  margin-left: 8px;
  font-size: 20px;
  color: orange;
}

.win-streak-count {
  font-size: 12px;
  font-weight: 700;
  color: black;
  display: flex;
  align-self: center;
}

.name {
  display: flex;
  flex-direction: row;
  width: 150px;
}
</style>
