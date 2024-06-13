export const FetchCategoriesAPI = (page: number = 0): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "categories/list?p=" + page;
};
