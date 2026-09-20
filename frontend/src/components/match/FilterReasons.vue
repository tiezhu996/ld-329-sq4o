<template>
  <div class="filter-list">
    <el-empty v-if="items.length === 0" description="没有被过滤的候选" :image-size="60" />
    <el-table v-else :data="items" size="small">
      <el-table-column label="候选交换对" min-width="200">
        <template #default="{ row }">
          <strong>{{ row.provider }}</strong>
          <span class="muted"> 提供「{{ row.offerSkill }}」</span>
          <br />
          <strong>{{ row.learner }}</strong>
          <span class="muted"> 需求「{{ row.wantedSkill }}」</span>
        </template>
      </el-table-column>
      <el-table-column prop="campus" label="校区" width="150" />
      <el-table-column label="共同时隙" min-width="160">
        <template #default="{ row }">
          <template v-if="row.commonSlots.length">
            <el-tag v-for="slot in row.commonSlots" :key="slot" size="small" effect="plain" class="gap-tag">
              {{ slot }}
            </el-tag>
          </template>
          <span v-else class="muted">无交集</span>
        </template>
      </el-table-column>
      <el-table-column label="过滤原因" width="230">
        <template #default="{ row }">
          <el-tag :type="tagType(row.reasonCode)" size="small">
            {{ reasonLabel(row.reasonCode) }}
          </el-tag>
          <p class="filter-reason">{{ row.reason }}</p>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { FILTER_REASON_LABEL, FILTER_REASON_TAG_TYPE } from '../../constants/match.constants';
import type { FilteredMatch } from '../../types/domain';

defineProps<{ items: FilteredMatch[] }>();

function reasonLabel(code: string): string {
  return FILTER_REASON_LABEL[code] ?? code;
}
function tagType(code: string): 'danger' | 'warning' | 'info' {
  return FILTER_REASON_TAG_TYPE[code] ?? 'info';
}
</script>

<style scoped>
.gap-tag {
  margin: 2px 4px 2px 0;
}
.filter-reason {
  margin: 4px 0 0;
  color: #667085;
  font-size: 12px;
  line-height: 1.4;
}
</style>
