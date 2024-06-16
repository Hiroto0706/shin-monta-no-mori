import { Category, ChildCategory } from "../category";

export interface FetchCategoriesResponse {
  categories: Category[];
}

export interface FetchChildCategoriesResponse {
  child_categories: ChildCategory[];
}

export interface GetChildCategoryResponse {
  child_category: ChildCategory | null;
}
