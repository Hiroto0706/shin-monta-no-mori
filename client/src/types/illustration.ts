import { Category } from "@/types/category";
import { Character } from "@/types/character";


export interface Image {
  id: number;
  title: string;
  original_src: string;
  original_filename: string;
  simple_src: string | null;
  simple_filename: string | null;
  created_at: string;
  updated_at:string;
}

export interface Illustration {
  Image: Image;
  Character: Character[];
  Category: Category[];
}