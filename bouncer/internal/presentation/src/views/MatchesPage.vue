<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PageTitle from '@/components/PageTitle.vue'
import PaginatedList from '@/components/PaginatedList.vue'
import MatchPreview from '@/components/MatchPreview.vue'

// TODO: type
const page = ref({ start: 1, size: 10, last: -1, elements: undefined })

function fetchPagination(): void {
  fetch(`http://localhost:5555/matches?start=${page.value.start}&size=${page.value.size}`, {
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      // match list is updated frequently, cache create UI inconsistencies
      cache: 'no-store',
    },
  }).then((resp) => {
    resp.json().then((matchesResp) => {
      if (matchesResp.summaries) {
        page.value = {
          ...page.value,
          elements: matchesResp.summaries,
          last: Math.ceil(matchesResp.total / matchesResp.size),
        }
      }
    })
  })
}

// TODO: type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function updatePage(updatedPage: any): void {
  page.value = { ...updatedPage }
  fetchPagination()
}

function updatePaginationStart(start: number): void {
  updatePage({ ...page.value, start })
}

function updatePaginationSize(size: number): void {
  updatePage({ ...page.value, size, start: 1 })
}

onMounted(() => {
  fetchPagination()
})
</script>

<template>
  <PageTitle title="Matches" />
  <PaginatedList
    :page="page"
    @updatePaginationSize="updatePaginationSize"
    @updatePaginationStart="updatePaginationStart"
  >
    <template #data="{ index }">
      <MatchPreview
        :index="page.start + index - 1"
        :preview="page.elements?.[index - 1]"
        :class="{ contrast: !!(index % 2) }"
      />
    </template>
    <template v-slot:separator>
      <div class="separator"></div>
    </template>
    <template v-slot:empty>
      <div>No match available</div>
    </template>
  </PaginatedList>
</template>

<style lang="css" scoped>
.preview {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.contrast {
  background-color: white;
}

.separator {
  width: 100%;
  height: 1px;
  background-color: black;
}
</style>
