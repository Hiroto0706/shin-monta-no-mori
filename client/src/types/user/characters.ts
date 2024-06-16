import { Character } from "../character";

export interface FetchCharactersResponse {
  characters: Character[];
}

export interface GetCharacterResponse {
  character: Character | null;
}
