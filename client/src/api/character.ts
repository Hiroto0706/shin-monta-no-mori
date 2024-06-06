export const FetchAllCharactersAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/list/all";
};

export const FetchCharactersAPI = (page: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/list?p=" + page;
};
