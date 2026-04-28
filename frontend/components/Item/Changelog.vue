<script setup lang="ts">
  import { toast } from "@/components/ui/sonner";
  import { useI18n } from "vue-i18n";
  import type { ChangelogEntry } from "~~/lib/api/types/data-contracts";
  import MdiHistory from "~icons/mdi/history";
  import MdiPlus from "~icons/mdi/plus";
  import DateTime from "~/components/global/DateTime.vue";
  import BaseCard from "@/components/Base/Card.vue";
  import { Button } from "@/components/ui/button";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";

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
  const submitting = ref(false);

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
      const { data, error } = await api.changelog.create(props.itemId, { summary });
      if (error) {
        toast.error(t("changelog.toast.failed_create"));
        return;
      }
      if (data) {
        entries.value.unshift(data);
      }
      newSummary.value = "";
      showForm.value = false;
      toast.success(t("changelog.toast.created"));
    } finally {
      submitting.value = false;
    }
  }

  onMounted(() => {
    loadEntries();
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
        <p class="text-sm">{{ entry.summary }}</p>
        <p class="text-xs text-foreground/50">
          <DateTime :date="entry.createdAt" />
        </p>
      </li>
    </ul>
  </BaseCard>
</template>
