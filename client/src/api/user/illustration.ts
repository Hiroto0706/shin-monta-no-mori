export const FetchIllustrationsAPI = (page: number = 0): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "illustrations/list?p=" + page;
};

export const SearchIllustrationsAPI = (
  page: number = 0,
  query: string | null
): string => {
  let url = process.env.NEXT_PUBLIC_BASE_API + "illustrations/search?p=" + page;
  if (query) {
    url += `&q=${encodeURIComponent(query)}`;
  }
  return url;
};

export const GetIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "illustrations/" + id;
};
