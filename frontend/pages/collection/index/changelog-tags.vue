<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import type { ChangelogTagOut } from "~~/lib/api/types/data-contracts";
  import MdiPlus from "~icons/mdi/plus";
  import MdiDelete from "~icons/mdi/delete";
  import MdiTag from "~icons/mdi/tag";
  import { Button } from "@/components/ui/button";
  import { Input } from "@/components/ui/input";
  import { Label } from "@/components/ui/label";
  import { Badge } from "@/components/ui/badge";
  import { Card } from "@/components/ui/card";
  import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

  const { t } = useI18n();
  const api = useUserApi();
  const confirm = useConfirm();

  const tags = ref<ChangelogTagOut[]>([]);
  const loadingTags = ref(false);

  async function loadTags() {
    loadingTags.value = true;
    try {
      const { data, error } = await api.changelogTags.getAll();
      if (error) {
        toast.error(t("changelog_tags.toast.failed_load"));
        return;
      }
      tags.value = data ?? [];
    } finally {
      loadingTags.value = false;
    }
  }

  const createForm = reactive({ name: "", color: "#6366f1" });
  const creating = ref(false);
  const showCreateForm = ref(false);

  function resetCreateForm() {
    createForm.name = "";
    createForm.color = "#6366f1";
    showCreateForm.value = false;
  }

  async function createTag() {
    if (!createForm.name.trim()) {
      toast.error(t("changelog_tags.toast.name_required"));
      return;
    }

    creating.value = true;
    try {
      const { error } = await api.changelogTags.create({
        name: createForm.name.trim(),
        color: createForm.color,
      });
      if (error) {
        toast.error(t("changelog_tags.toast.failed_create"));
        return;
      }
      toast.success(t("changelog_tags.toast.created"));
      resetCreateForm();
      await loadTags();
    } finally {
      creating.value = false;
    }
  }

  async function deleteTag(tag: ChangelogTagOut) {
    const { isCanceled } = await confirm.open(t("changelog_tags.delete_confirm", { name: tag.name }));
    if (isCanceled) return;

    const { error } = await api.changelogTags.delete(tag.id);
    if (error) {
      toast.error(t("changelog_tags.toast.failed_delete"));
      return;
    }
    toast.success(t("changelog_tags.toast.deleted"));
    await loadTags();
  }

  /** Returns an inline style for the badge from the stored colour value. */
  function badgeStyle(color: string | undefined | null): Record<string, string> {
    if (!color) return {};
    return { backgroundColor: color, color: "#fff", borderColor: color };
  }

  onMounted(loadTags);
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-4 flex items-center justify-between">
      <h3 class="text-lg font-medium">{{ $t("changelog_tags.title") }}</h3>
      <Button size="sm" @click="showCreateForm = !showCreateForm">
        <MdiPlus class="mr-1 size-4" />
        {{ $t("changelog_tags.add") }}
      </Button>
    </div>

    <!-- Create form -->
    <Card v-if="showCreateForm" class="mb-4 p-4">
      <form class="flex flex-col gap-3" @submit.prevent="createTag">
        <div class="flex flex-col gap-1">
          <Label for="tag-name">{{ $t("changelog_tags.name_label") }}</Label>
          <Input
            id="tag-name"
            v-model="createForm.name"
            :placeholder="$t('changelog_tags.name_placeholder')"
            maxlength="100"
            autofocus
          />
        </div>
        <div class="flex flex-col gap-1">
          <Label for="tag-color">{{ $t("changelog_tags.color_label") }}</Label>
          <div class="flex items-center gap-3">
            <input
              id="tag-color"
              v-model="createForm.color"
              type="color"
              class="size-9 cursor-pointer rounded border border-input p-0.5"
            />
            <span class="text-sm text-muted-foreground">{{ createForm.color }}</span>
          </div>
        </div>
        <div class="flex gap-2">
          <Button type="submit" size="sm" :disabled="creating || !createForm.name.trim()">
            {{ $t("global.save") }}
          </Button>
          <Button type="button" size="sm" variant="outline" @click="resetCreateForm">
            {{ $t("global.cancel") }}
          </Button>
        </div>
      </form>
    </Card>

    <!-- Tags list -->
    <div v-if="loadingTags" class="py-6 text-center text-sm text-muted-foreground">
      {{ $t("global.loading") }}
    </div>
    <div v-else-if="tags.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
      <MdiTag class="mb-3 size-10 text-muted-foreground/50" />
      <p class="mb-4 text-muted-foreground">{{ $t("changelog_tags.empty") }}</p>
      <Button @click="showCreateForm = true">
        <MdiPlus class="mr-2 size-4" />
        {{ $t("changelog_tags.add") }}
      </Button>
    </div>
    <div v-else class="space-y-2">
      <Card v-for="tag in tags" :key="tag.id" class="p-4">
        <div class="flex items-center gap-3">
          <Badge variant="secondary" class="shrink-0 text-xs" :style="badgeStyle(tag.color)">
            {{ tag.name }}
          </Badge>
          <span class="mr-auto text-sm font-medium">{{ tag.name }}</span>
          <TooltipProvider :delay-duration="0">
            <Tooltip>
              <TooltipTrigger as-child>
                <Button variant="ghost" size="icon" class="size-8 text-destructive" @click="deleteTag(tag)">
                  <MdiDelete class="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>{{ $t("global.delete") }}</TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </Card>
    </div>
  </div>
</template>
