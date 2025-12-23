<script setup lang="ts">
import PaginationSettings from './PaginationSettings.vue'

// TODO: type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
defineProps<{ page: any }>()
const emits = defineEmits<{
  (e: 'updatePaginationSize', size: number): void
  (e: 'updatePaginationStart', start: number): void
}>()

function updatePaginationStart(start: number) {
  emits('updatePaginationStart', start)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="pagination">
    <pagination-settings
      :page="page"
      @updatePaginationSize="emits('updatePaginationSize', $event)"
      @updatePaginationStart="updatePaginationStart"
    />

    <div v-if="page.elements" class="elements">
      <div
        v-for="index in page.elements.length ?? 0"
        :key="index"
        class="element-wrapper"
        :class="{
          'clip-top-border-radius ': index == 1,
          'clip-bottom-border-radius ': index == page.elements.length,
        }"
      >
        <template v-if="index > 1">
          <slot name="separator" />
        </template>
        <slot name="data" v-bind="{ index }" />
      </div>
      <div v-if="page.elements.length === 0">
        <slot name="empty">
          <div>No element available</div>
        </slot>
      </div>
    </div>

    <pagination-settings
      v-if="(page?.elements?.length ?? 0) > 10"
      :page="page"
      @updatePaginationSize="emits('updatePaginationSize', $event)"
      @updatePaginationStart="updatePaginationStart"
    />
  </div>
</template>

<style lang="css" scoped>
.pagination {
  display: flex;
  flex-direction: column;
  padding: 16px;
}

.elements {
  margin: 16px;
  border: 1px solid black;
  border-radius: 8px;
}

.clip-top-border-radius {
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
  overflow-x: hidden;
}

.clip-bottom-border-radius {
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
  overflow-x: hidden;
}

.element-wrapper {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}
</style>
