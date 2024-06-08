export const FetchAllCharactersAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/list/all";
};

export const FetchCharactersAPI = (page: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/list?p=" + page;
};

export const CreateCharacterAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/create";
};

export const GetCharacterAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/" + id;
};

export const EditCharacterAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/" + id;
};

export const DeleteCharacterAPI = (id: number): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "admin/characters/" + id;
};
