import { Category } from "@/types/category";
import { Character } from "@/types/character";

export interface Image {
  id: number;
  title: string;
  original_src: string;
  original_filename: string;
  simple_src: {
    String: string;
    Valid: boolean;
  };
  simple_filename: {
    String: string;
    Valid: boolean;
  };
  created_at: string;
  updated_at: string;
}

export interface Illustration {
  Image: Image;
  Characters: { Character: Character }[];
  Categories: Category[];
}

// fetchIllustrations で期待する型
export interface FetchIllustrationsResponse {
  illustrations: Illustration[];
  total_pages: number;
  total_count: number;
}

export interface GetIllustrationResponse {
  illustration: Illustration | null;
}
