export const FetchAllCharactersAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "characters/list/all";
};

export const SearchCharactersAPI = (
  page: number,
  query: string | null
): string => {
  let url = process.env.NEXT_PUBLIC_BASE_API + "characters/search?p=" + page;
  if (query) {
    url += `&q=${encodeURIComponent(query)}`;
  }
  return url;
};

export const GetCharacterAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "characters/" + id;
};
