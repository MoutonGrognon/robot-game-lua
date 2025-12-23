<!-- eslint-disable vue/valid-v-slot -->
<script setup lang="ts">
import { computed } from 'vue'

// TODO: type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const props = defineProps<{ page: any }>()

const emits = defineEmits<{
  (e: 'updatePaginationSize', size: number): void
  (e: 'updatePaginationStart', start: number): void
}>()

const sizes = [10, 20, 50]

const currentPageIndex = computed(() => Math.ceil(props.page.start / props.page.size))

function getStartIndexFromPageIndex(pageIndex: number): number {
  return (pageIndex - 1) * props.page.size + 1
}
</script>

<template>
  <div class="pagination-settings">
    <div class="sizes">
      <div>rows :</div>
      <button
        v-for="size in sizes"
        :key="size"
        @click="emits('updatePaginationSize', size)"
        class="button"
        :class="{ current: size === page.size }"
        v-bind:disabled="size === page.size"
      >
        {{ size }}
      </button>
    </div>

    <div v-if="page.last !== -1" class="starts">
      <div v-if="currentPageIndex !== 1" class="row-count-selector">
        <button
          class="button"
          @click="emits('updatePaginationStart', getStartIndexFromPageIndex(currentPageIndex - 1))"
        >
          previous
        </button>
        <button class="button" @click="emits('updatePaginationStart', 1)">1</button>
      </div>
      <div v-if="currentPageIndex - 3 > 1">...</div>
      <div v-for="offset in 2" :key="offset">
        <button
          class="button"
          v-if="currentPageIndex - 3 + offset > 1"
          @click="
            emits(
              'updatePaginationStart',
              getStartIndexFromPageIndex(currentPageIndex - 3 + offset),
            )
          "
        >
          {{ currentPageIndex - 3 + offset }}
        </button>
      </div>
      <button class="button current" disabled>{{ currentPageIndex }}</button>
      <div v-for="offset in 2" :key="offset" class="row-count-selector">
        <button
          class="button"
          v-if="currentPageIndex + offset < page.last"
          @click="
            emits('updatePaginationStart', getStartIndexFromPageIndex(currentPageIndex + offset))
          "
        >
          {{ currentPageIndex + offset }}
        </button>
      </div>
      <div v-if="page.last - 3 > currentPageIndex">...</div>
      <div v-if="currentPageIndex !== page.last" class="row-count-selector">
        <button
          class="button"
          @click="emits('updatePaginationStart', getStartIndexFromPageIndex(page.last))"
        >
          {{ page.last }}
        </button>
        <button
          class="button"
          @click="emits('updatePaginationStart', getStartIndexFromPageIndex(currentPageIndex + 1))"
        >
          next
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.button {
  border: none;
  background-color: transparent;
  color: black;
  font-family: inherit;
  font-size: inherit;
}

.button.current {
  text-decoration: underline;
}

.button:hover {
  text-decoration: underline;
}

.pagination-settings {
  display: flex;
  flex-direction: row;
  padding: 0px 32px;
}

.sizes {
  display: flex;
  flex-direction: row;
}

.starts {
  display: flex;
  flex-direction: row;
  margin-left: auto;
}

.row-count-selector {
  display: flex;
  flex-direction: row;
}
</style>
