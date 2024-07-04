import { Category, ChildCategory } from "../category";

export interface FetchCategoriesResponse {
  categories: Category[];
  total_pages: number;
  total_count: number;
}

export interface GetCategoryResponse {
  category: Category | null;
}

export interface GetChildCategoryResponse {
  child_category: ChildCategory | null;
}
