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

export const GetCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/" + id;
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

export const GetChildCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/child/" + id;
};

export const EditChildCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/child/" + id;
};

export const DeleteChildCategoryAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/categories/child/" + id;
};
