export const FetchCategoriesAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/list";
};

export const SearchCategoriesAPI = (query: string): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/categories/search?q=" + query
  );
};

export const CreateParentCategoryAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/parent/create";
};

export const GetParentCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/parent/" + id;
};

export const EditParentCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/parent/" + id;
};

export const DeleteParentCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/parent/" + id;
};

export const CreateChildCategoryAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/child/create";
};
