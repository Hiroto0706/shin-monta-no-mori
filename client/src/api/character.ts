export const FetchCharactersAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/list";
};
