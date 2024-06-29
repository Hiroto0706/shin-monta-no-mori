import axios from "axios";
import UserHeader from "@/components/user/common/header/header";
import { FetchCategoriesAllAPI } from "@/api/user/category";
import { FetchCategoriesResponse } from "@/types/user/categories";
import { FetchCharactersResponse } from "@/types/user/characters";
import { FetchAllCharactersAPI } from "@/api/user/character";

const fetchCharacters = async (): Promise<FetchCharactersResponse> => {
  try {
    const response = await axios.get(FetchAllCharactersAPI());
    return response.data;
  } catch (error) {
    console.error("キャラクターの取得に失敗しました", error);
    return { characters: [] };
  }
};

const fetchCategories = async (): Promise<FetchCategoriesResponse> => {
  try {
    const response = await axios.get(FetchCategoriesAllAPI());
    return response.data;
  } catch (error) {
    console.error("カテゴリの取得に失敗しました", error);
    return { categories: [] };
  }
};

const Header = async () => {
  const fetchCharactersRes = await fetchCharacters();
  const fetchCategoriesRes = await fetchCategories();

  return (
    <header>
      <UserHeader
        characters={fetchCharactersRes.characters}
        categories={fetchCategoriesRes.categories}
      />
    </header>
  );
};

export default Header;
