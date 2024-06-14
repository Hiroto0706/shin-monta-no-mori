import { Category, ChildCategory } from "../category";

export interface FetchCategoriesResponse {
  categories: Category[];
}

export interface GetCategoryResponse {
  category: Category | null;
}

export interface GetChildCategoryResponse {
  child_category: ChildCategory | null;
}
