<template>
  <article class="feature-card match-card">
    <div class="feature-card__title">
      <strong>{{ match.provider }} × {{ match.learner }}</strong>
      <el-tag type="warning">{{ match.score }}%</el-tag>
    </div>
    <p class="muted">{{ match.offerSkill }} ↔ {{ match.wantedSkill }}</p>
    <p>{{ match.recommendation }}</p>

    <div class="match-meta">
      <el-tag size="small" effect="plain">{{ match.campus }}</el-tag>
      <el-tag size="small" type="success" effect="plain">回报：{{ match.reward }}</el-tag>
    </div>

    <div class="basis-row">
      <span class="muted">匹配依据：</span>
      <el-tag v-for="item in match.basis" :key="item" size="small" class="basis-tag">{{ item }}</el-tag>
    </div>

    <div class="blocks-row">
      <span class="muted">未来一周连续可交换时段：</span>
      <el-tag
        v-for="block in match.commonBlocks"
        :key="block.slotText"
        size="small"
        type="primary"
        effect="dark"
      >
        {{ block.slotText }}
        <small v-if="block.parts.length > 1">（连续 {{ block.parts.length }} 段）</small>
      </el-tag>
    </div>

    <div v-if="canInvite" class="invite-row">
      <el-select v-model="chosenSlot" size="small" style="width: 180px">
        <el-option
          v-for="block in match.commonBlocks"
          :key="block.slotText"
          :label="block.slotText"
          :value="block.slotText"
        />
      </el-select>
      <el-button type="primary" size="small" :loading="acting" @click="emitInvite">
        向 {{ match.learner }} 发起邀请
      </el-button>
    </div>
    <el-alert
      v-else-if="inviteHint"
      :title="inviteHint"
      :type="readOnly ? 'error' : 'info'"
      :closable="false"
      show-icon
      class="invite-hint"
    />
  </article>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import type { Match } from '../../types/domain';

const props = defineProps<{
  match: Match;
  viewer: string;
  readOnly: boolean;
  acting: boolean;
  activeByNeed: Record<string, number>;
  maxActive: number;
}>();

const emit = defineEmits<{
  (event: 'invite', match: Match, slot: string): void;
}>();

const chosenSlot = ref(props.match.commonBlocks[0]?.slotText ?? props.match.commonSlots[0] ?? '');
watch(
  () => props.match.id,
  () => {
    chosenSlot.value = props.match.commonBlocks[0]?.slotText ?? props.match.commonSlots[0] ?? '';
  },
);

const activeCount = computed(
  () => props.activeByNeed[String(props.match.needId)] ?? 0,
);

const canInvite = computed(() => props.match.viewerCanInvite);

const inviteHint = computed(() => {
  if (props.readOnly) {
    return '信用分低于 80，仅保留浏览权限，不能发起邀请';
  }
  // 查看者是技能提供方：不能发起仅可能因为该需求已满三项或本人已邀请。
  if (props.match.provider === props.viewer) {
    if (activeCount.value >= props.maxActive) {
      return `该需求已有 ${props.maxActive} 项有效邀请，达到上限`;
    }
    return '你已对该需求发起有效邀请，等待对方确认';
  }
  // 查看者是需求方：邀请在下方列表中接受。
  if (props.match.learner === props.viewer) {
    return '你是该需求的需求方，可在下方「邀请与终态」中接受待确认邀请';
  }
  return '你不是该匹配的当事人，仅可浏览';
});

function emitInvite(): void {
  emit('invite', props.match, chosenSlot.value);
}
</script>
