<script setup lang="ts">
import { ref } from 'vue'
import PageTitle from '@/components/PageTitle.vue'
import PaginatedList from '@/components/PaginatedList.vue'
import RankingPreview from '@/components/RankingPreview.vue'

// TODO: type
const page = ref({ start: 1, size: 10, last: -1, elements: undefined })

function fetchRanking(): void {
  // No need for pagination the total amount of bot should be quite low
  fetch(`http://localhost:5555/ranking`, {
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      // ranking is updated frequently, caching is annoying for the user
      cache: 'no-store',
    },
  }).then((resp) => {
    resp.json().then((rankingResp) => {
      if (rankingResp.ranking) {
        // TODO: type
        const sortedRanking = rankingResp.ranking.toSorted(
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          (rankA: any, rankB: any) => rankB.elo - rankA.elo,
        )
        // TODO: type
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const elements = sortedRanking.reduce((ranking: any[], rank: any) => {
          const prevRank = ranking.length > 0 ? ranking[ranking.length - 1] : undefined
          const prevElo = prevRank?.elo ?? null
          ranking.push({
            ...rank,
            rank: rank.elo === prevElo ? prevRank.rank : ranking.length + 1,
          })
          return ranking
        }, [])
        page.value = {
          ...page.value,
          elements,
          last: Math.ceil(rankingResp.ranking.length / page.value.size),
        }
      }
    })
  })
}

// TODO: type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function updatePage(updatedPage: any): void {
  page.value = { ...updatedPage }
}

function updatePaginationStart(start: number): void {
  updatePage({ ...page.value, start })
}

function updatePaginationSize(size: number): void {
  updatePage({ ...page.value, size, start: 1 })
}

fetchRanking()
</script>

<template>
  <PageTitle title="Ranking" />
  <PaginatedList
    :page="page"
    @updatePaginationSize="updatePaginationSize"
    @updatePaginationStart="updatePaginationStart"
  >
    <template #data="{ index }">
      <RankingPreview
        :index="page.start + index - 1"
        :preview="page.elements?.[index - 1]"
        :class="{ contrast: !!(index % 2) }"
      />
    </template>
    <template v-slot:separator>
      <div class="separator"></div>
    </template>
    <template v-slot:empty>
      <div>No ranking available</div>
    </template>
  </PaginatedList>
</template>

<style lang="css" scoped>
.contrast {
  background-color: white;
}

.separator {
  width: 100%;
  height: 1px;
  background-color: black;
}
</style>
