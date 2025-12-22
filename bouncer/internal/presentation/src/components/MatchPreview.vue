<script setup lang="ts">
import { computed } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { library } from '@fortawesome/fontawesome-svg-core';

/* import all the icons in Free Solid, Free Regular, and Brands styles */
import { fas } from '@fortawesome/free-solid-svg-icons';
import { far } from '@fortawesome/free-regular-svg-icons';
import { fab } from '@fortawesome/free-brands-svg-icons';
library.add(fas, far, fab);

// TODO: type
const props = defineProps<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  preview: any;
  index: number;
}>();

function getResultIcon(result: number): string {
  if (result == 0) {
    return "equals";
  }
  return result > 0 ? "check" : "xmark";
}

const trueBlueResult = computed(() => {
  return props.preview.blueScore - props.preview.redScore;
});
</script>

<template>
  <RouterLink :to="`matches/${preview.id}`" class="wrapper">
    <div class="count">{{ index }}</div>
    <div class="blue-name" :class="{ 'winner': trueBlueResult > 0, 'loser': trueBlueResult < 0 }">
      {{ preview.blueBotName }}
    </div>
    <FontAwesomeIcon :icon="['fas', getResultIcon(trueBlueResult)]"
      :class="{ 'winner': trueBlueResult > 0, 'loser': trueBlueResult < 0 }" />
    <div class="score">
      <span class="blue-score">{{ preview.blueScore }}</span>
      <span>-</span>
      <span class="red-score">{{ preview.redScore }}</span>
    </div>
    <FontAwesomeIcon :icon="['fas', getResultIcon(-trueBlueResult)]"
      :class="{ 'winner': trueBlueResult < 0, 'loser': trueBlueResult > 0 }" />
    <div class="red-name" :class="{ 'winner': trueBlueResult < 0, 'loser': trueBlueResult > 0 }">
      {{ preview.redBotName }}
    </div>
  </RouterLink>
</template>

<style lang="css" scoped>
.wrapper {
  display: flex;
  flex-direction: row;
  padding: 8px;
  column-gap: 32px;
  width: 100%;

  color: black;
  text-decoration: none;
}

.count {
  width: 50px
}

.blue-name,
.red-name {
  width: 150px;
}

.blue-name {
  text-align: right;
}

.red-name {
  text-align: left;
}

.blue-score,
.red-score {
  padding: 0 8px;
  width: 20px;
  display: inline-block;
}

.blue-score {
  text-align: right;
}

.red-score {
  text-align: left;
}

.score {
  font-weight: 700;
  text-align: center;
}

.winner {
  color: green;
}

.loser {
  color: red;
}
</style>
