<script setup lang="ts">
// TODO: type
defineProps<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  preview: any
  index: number
}>()
</script>

<template>
  <div class="wrapper">
    <div class="count">{{ preview.rank }}</div>
    <div class="elo">
      {{ preview.elo }}
    </div>
    <div class="name">
      {{ preview.botName }}
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
  width: 50px;
}

.name {
  width: 150px;
}

.elo {
  font-weight: 700;
  width: 60px;
}
</style>
