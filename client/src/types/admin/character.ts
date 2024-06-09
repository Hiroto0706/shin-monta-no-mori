import { Character } from "../character";

export interface FetchCharactersResponse {
  characters: Character[];
  total_pages: number;
  total_count: number;
}

export interface GetCharacterResponse {
  character: Character | null;
}