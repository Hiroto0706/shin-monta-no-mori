import axios from "axios";
import SidebarMain from "./sidebarMain";
import { Character } from "@/types/character";
import { FetchAllCharactersAPI } from "@/api/user/character";
import { FetchCategoriesAllAPI } from "@/api/user/category";
import { FetchCategoriesResponse } from "@/types/user/categories";

const fetchCategories = async (): Promise<FetchCategoriesResponse> => {
  try {
    const response = await axios.get(FetchCategoriesAllAPI());
    return response.data;
  } catch (error) {
    console.error("カテゴリの取得に失敗しました", error);
    return { categories: [] };
  }
};

export const fetchCharacters = async (): Promise<Character[] | undefined> => {
  try {
    const response = await axios.get(FetchAllCharactersAPI());

    return response.data.characters;
  } catch (error) {
    console.error(error);
    return undefined;
  }
};

const UserSidebar = async () => {
  const fetchCategoriesRes = await fetchCategories();
  const links = [
    {
      id: 0,
      href: "/categories",
      icon: "/icon/category.png",
      icon_active: "/icon/category-active.png",
      text: "カテゴリ",
    },
    {
      id: 1,
      href: "/characters",
      icon: "/icon/character.png",
      icon_active: "/icon/character-active.png",
      text: "キャラ",
    },
  ];

  return (
    <>
      <div className="duration-200 transform -translate-x-full md:transform-none">
        <SidebarMain links={links} categories={fetchCategoriesRes.categories} />
      </div>
    </>
  );
};

export default UserSidebar;
