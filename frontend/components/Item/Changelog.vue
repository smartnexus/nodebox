<script setup lang="ts">
  import { toast } from "@/components/ui/sonner";
  import { useI18n } from "vue-i18n";
  import type { ChangelogEntry, ChangelogTagOut } from "~~/lib/api/types/data-contracts";
  import MdiHistory from "~icons/mdi/history";
  import MdiPlus from "~icons/mdi/plus";
  import DateTime from "~/components/global/DateTime.vue";
  import BaseCard from "@/components/Base/Card.vue";
  import { Button } from "@/components/ui/button";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
  import { Badge } from "@/components/ui/badge";
  import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

  interface Props {
    itemId: string;
  }

  const props = defineProps<Props>();

  const { t } = useI18n();
  const api = useUserApi();

  const entries = ref<ChangelogEntry[]>([]);
  const loading = ref(false);
  const showForm = ref(false);
  const newSummary = ref("");
  const newTagId = ref("");
  const submitting = ref(false);

  const availableTags = ref<ChangelogTagOut[]>([]);

  async function loadTags() {
    const { data } = await api.changelogTags.getAll();
    availableTags.value = data ?? [];
  }

  async function loadEntries() {
    loading.value = true;
    try {
      const { data, error } = await api.changelog.getByItemId(props.itemId);
      if (error) {
        toast.error(t("changelog.toast.failed_load"));
        return;
      }
      entries.value = data ?? [];
    } finally {
      loading.value = false;
    }
  }

  async function createEntry() {
    const summary = newSummary.value.trim();
    if (!summary) {
      return;
    }

    submitting.value = true;
    try {
      const payload: { summary: string; tagId?: string } = { summary };
      if (newTagId.value) {
        payload.tagId = newTagId.value;
      }

      const { data, error } = await api.changelog.create(props.itemId, payload);
      if (error) {
        toast.error(t("changelog.toast.failed_create"));
        return;
      }
      if (data) {
        entries.value.unshift(data);
      }
      newSummary.value = "";
      newTagId.value = "";
      showForm.value = false;
      toast.success(t("changelog.toast.created"));
    } finally {
      submitting.value = false;
    }
  }

  /** Return badge colour classes from a hex/named colour, falling back to secondary. */
  function badgeStyle(color: string | undefined | null): Record<string, string> {
    if (!color) return {};
    return { backgroundColor: color, color: "#fff", borderColor: color };
  }

  onMounted(() => {
    loadEntries();
    loadTags();
  });
</script>

<template>
  <BaseCard collapsable>
    <template #title>
      <MdiHistory class="mr-2 size-5" />
      {{ $t("changelog.title") }}
    </template>
    <template #title-actions>
      <Button size="sm" variant="outline" @click="showForm = !showForm">
        <MdiPlus class="mr-1 size-4" />
        {{ $t("changelog.add_entry") }}
      </Button>
    </template>

    <!-- Add Entry Form -->
    <div v-if="showForm" class="border-b p-4">
      <div class="flex flex-col gap-3">
        <Label for="changelog-summary">{{ $t("changelog.summary_label") }}</Label>
        <Input
          id="changelog-summary"
          v-model="newSummary"
          :placeholder="$t('changelog.summary_placeholder')"
          maxlength="500"
          @keyup.enter="createEntry"
        />

        <!-- Tag selector -->
        <div v-if="availableTags.length > 0" class="flex flex-col gap-1">
          <Label>{{ $t("changelog.tag_label") }}</Label>
          <Select v-model="newTagId">
            <SelectTrigger>
              <SelectValue :placeholder="$t('changelog.tag_placeholder')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">{{ $t("changelog.tag_none") }}</SelectItem>
              <SelectItem v-for="tag in availableTags" :key="tag.id" :value="tag.id">
                {{ tag.name }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="flex gap-2">
          <Button size="sm" :disabled="submitting || !newSummary.trim()" @click="createEntry">
            {{ $t("global.save") }}
          </Button>
          <Button
            size="sm"
            variant="outline"
            @click="
              showForm = false;
              newSummary = '';
              newTagId = '';
            "
          >
            {{ $t("global.cancel") }}
          </Button>
        </div>
      </div>
    </div>

    <!-- Entries list -->
    <div v-if="loading" class="px-4 py-6 text-center text-sm text-foreground/60">
      {{ $t("global.loading") }}
    </div>
    <div v-else-if="entries.length === 0" class="px-4 py-6 text-center text-sm text-foreground/60">
      {{ $t("changelog.no_entries") }}
    </div>
    <ul v-else class="divide-y">
      <li v-for="entry in entries" :key="entry.id" class="flex flex-col gap-1 px-4 py-3">
        <div class="flex flex-wrap items-center gap-2">
          <p class="flex-1 text-sm">{{ entry.summary }}</p>
          <Badge v-if="entry.tag" variant="secondary" class="shrink-0 text-xs" :style="badgeStyle(entry.tag.color)">
            {{ entry.tag.name }}
          </Badge>
        </div>
        <p class="text-xs text-foreground/50">
          <DateTime :date="entry.createdAt" />
        </p>
      </li>
    </ul>
  </BaseCard>
</template>
