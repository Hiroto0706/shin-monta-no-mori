export const FetchIllustrationsAPI = (page: number): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/list?p=" + page
  );
};

export const SearchIllustrationsAPI = (
  page: number = 0,
  query: string | null,
  characters?: string | null,
  categories?: string | null
): string => {
  let url =
    process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/search?p=" + page;
  if (query) {
    url += `&q=${encodeURIComponent(query)}`;
  }
  if (characters) {
    url += `&characters=${encodeURIComponent(characters)}`;
  }
  if (categories) {
    url += `&categories=${encodeURIComponent(categories)}`;
  }
  return url;
};

export const GetIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/" + id;
};

export const CreateIllustrationAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/create";
};

export const EditIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/" + id;
};

export const DeleteIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/illustrations/" + id;
};
