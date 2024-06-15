import { ChildCategory } from "../category";

export interface FetchChildCategoriesResponse {
  child_categories: ChildCategory[];
}

export interface GetChildCategoryResponse {
  child_category: ChildCategory | null;
}
