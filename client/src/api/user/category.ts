export const FetchChildCategoriesAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "categories/child/list";
};
