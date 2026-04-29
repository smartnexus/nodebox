import { BaseAPI, route } from "../base";
import type { ChangelogTagCreate, ChangelogTagOut } from "../types/data-contracts";

export class ChangelogTagsAPI extends BaseAPI {
  /**
   * Retrieve all changelog tags for the current group.
   */
  getAll() {
    return this.http.get<ChangelogTagOut[]>({
      url: route("/changelog-tags"),
    });
  }

  /**
   * Create a new changelog tag for the current group.
   */
  create(data: ChangelogTagCreate) {
    return this.http.post<ChangelogTagCreate, ChangelogTagOut>({
      url: route("/changelog-tags"),
      body: data,
    });
  }

  /**
   * Delete a changelog tag by ID.
   */
  delete(id: string) {
    return this.http.delete<void>({
      url: route(`/changelog-tags/${id}`),
    });
  }
}
