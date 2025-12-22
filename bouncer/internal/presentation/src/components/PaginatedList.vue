<script setup lang="ts">
import { computed } from 'vue'

// TODO: type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const props = defineProps<{ page: any }>()
defineEmits<{
  (e: 'updatePaginationSize', size: number): void
  (e: 'updatePaginationStart', start: number): void
}>()

const sizes = [10, 20, 50]

const currentPage = computed(() => Math.ceil(props.page.start / props.page.size))

function getStartIndexFromPageIndex(pageIndex: number): number {
  return (pageIndex - 1) * props.page.size + 1
}
</script>

<template>
  <div class="header">
    <div class="sizes">
      <div>rows :</div>
      <button
        v-for="size in sizes"
        :key="size"
        @click="$emit('updatePaginationSize', size)"
        class="button"
        :class="{ current: size === page.size }"
        v-bind:disabled="size === page.size"
      >
        {{ size }}
      </button>
    </div>

    <div v-if="page.last !== -1" class="starts">
      <div v-if="currentPage !== 1" class="row">
        <button
          class="button"
          @click="$emit('updatePaginationStart', getStartIndexFromPageIndex(currentPage - 1))"
        >
          previous
        </button>
        <button class="button" @click="$emit('updatePaginationStart', 1)">1</button>
      </div>
      <div v-if="currentPage - 3 > 1">...</div>
      <div v-for="offset in 2" :key="offset">
        <button
          class="button"
          v-if="currentPage - 3 + offset > 1"
          @click="
            $emit('updatePaginationStart', getStartIndexFromPageIndex(currentPage - 3 + offset))
          "
        >
          {{ currentPage - 3 + offset }}
        </button>
      </div>
      <button class="button current" disabled>{{ currentPage }}</button>
      <div v-for="offset in 2" :key="offset" class="row">
        <button
          class="button"
          v-if="currentPage + offset < page.last"
          @click="$emit('updatePaginationStart', getStartIndexFromPageIndex(currentPage + offset))"
        >
          {{ currentPage + offset }}
        </button>
      </div>
      <div v-if="page.last - 3 > currentPage">...</div>
      <div v-if="currentPage !== page.last" class="row">
        <button
          class="button"
          @click="$emit('updatePaginationStart', getStartIndexFromPageIndex(page.last))"
        >
          {{ page.last }}
        </button>
        <button
          class="button"
          @click="$emit('updatePaginationStart', getStartIndexFromPageIndex(currentPage + 1))"
        >
          next
        </button>
      </div>
    </div>
  </div>

  <div v-if="page.elements">
    <slot name="separator"></slot>
    <div v-for="index in page.elements.length ?? 0" :key="index">
      <slot name="data" v-bind="{ index }" />
      <slot name="separator"></slot>
    </div>
    <div v-if="page.elements.length === 0">
      <slot name="empty">
        <div>No element available</div>
      </slot>
      <slot name="separator"></slot>
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

.header {
  display: flex;
  flex-direction: row;
  padding: 16px 32px;
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

.row {
  display: flex;
  flex-direction: row;
}
</style>
