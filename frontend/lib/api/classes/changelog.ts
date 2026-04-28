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
   */
  create(itemId: string, data: ChangelogEntryCreate) {
    return this.http.post<ChangelogEntryCreate, ChangelogEntry>({
      url: route(`/entities/${itemId}/changelog`),
      body: data,
    });
  }
}
