export const FetchCategoriesAllAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "categories/list/all";
};

export const FetchChildCategoriesAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "categories/child/list";
};

export const GetChildCategoryAPI = (child_category_id: number): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "categories/child/" + child_category_id
  );
};
