export const FetchIllustrationsAPI = (page: number = 0): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "illustrations/list?p=" + page;
};

export const SearchIllustrationsAPI = (
  page: number = 0,
  query: string
): string => {
  let url = process.env.NEXT_PUBLIC_BASE_API + "illustrations/search?p=" + page;
  if (query != "") {
    url += `&q=${encodeURIComponent(query)}`;
  }
  return url;
};

export const FetchIllustrationsByCategoryAPI = (
  category_id: number,
  page: number
): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API +
    "illustrations/category/child/" +
    category_id +
    "?p=" +
    page
  );
};

export const FetchIllustrationsByCharacterAPI = (
  character_id: number,
  page: number
): string => {
  return (
    process.env.NEXT_PUBLIC_BASE_API +
    "illustrations/character/" +
    character_id +
    "?p=" +
    page
  );
};

export const GetIllustrationAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "illustrations/" + id;
};
