<script setup lang="ts">
import router from '@/router'
import type { RouteRecordRaw } from 'vue-router'
defineProps<{
  tabs: RouteRecordRaw[]
}>()

function pathMatchCurrentRoute(path: string): boolean {
  return router.currentRoute.value.path.startsWith(path)
}
</script>

<template>
  <nav>
    <div class="tab">
      <RouterLink
        v-for="tab in tabs"
        :key="tab.path"
        :to="tab.path"
        class="link"
        :class="{
          'selected-link': pathMatchCurrentRoute(tab.path),
          'active-link': !pathMatchCurrentRoute(tab.path),
        }"
      >
        {{ tab.name }}
      </RouterLink>
    </div>
  </nav>
</template>

<style lang="css" scoped>
nav {
  width: 100%;
  background-color: white;
}

.tab {
  display: flex;
}

.link {
  display: flex;
  margin: 0 24px;
  padding-top: 24px;
  padding-bottom: 20px;
  margin-bottom: 4px;
  width: fit-content;
  text-decoration: none;
  color: black;

  background:
    linear-gradient(black 0 0) bottom/ 0% 2px no-repeat,
    white;
  transition-property: background;
  transition-duration: 0.1s;
}

.selected-link,
.active-link:hover {
  background:
    linear-gradient(black 0 0) bottom/ 100% 2px no-repeat,
    white;
}
</style>
