export const FetchCategoriesAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/list";
};

export const SearchCategoriesAPI = (query: string): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/categories/search?q=" + query
  );
};
