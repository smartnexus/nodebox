import { BaseAPI, route } from "../base";
import type { ChangelogEntry, ChangelogEntryCreate } from "../types/data-contracts";

export class ChangelogAPI extends BaseAPI {
  /**
   * Retrieve all changelog entries for a given item, ordered newest-first.
   */
  getByItemId(itemId: string) {
    return this.http.get<ChangelogEntry[]>({
      url: route(`/entities/${itemId}/changelog`),
    });
  }

  /**
   * Create a new changelog entry for the given item.
   * tagId is optional - if provided it must be a valid ChangelogTag ID for the current group.
   */
  create(itemId: string, data: { summary: string; tagId?: string }) {
    return this.http.post<Partial<ChangelogEntryCreate>, ChangelogEntry>({
      url: route(`/entities/${itemId}/changelog`),
      body: data,
    });
  }
}
