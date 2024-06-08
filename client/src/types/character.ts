export interface Character {
  id: number;
  name: string;
  src: string;
  filename: {
    String: string;
    Valid: boolean;
  };
  created_at: string;
  updated_at: string;
}

export interface FetchCharactersResponse {
  characters: Character[];
  total_pages: number;
  total_count: number;
}

export interface GetCharacterResponse {
  character: Character | null;
}