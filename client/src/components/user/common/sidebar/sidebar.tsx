import axios from "axios";
import SidebarMain from "./sidebarMain";
import { Character } from "@/types/character";
import { FetchAllCharactersAPI } from "@/api/user/character";
import { FetchCategoriesAllAPI } from "@/api/user/category";
import { Category } from "@/types/category";

export const fetchCategories = async (): Promise<Category[] | undefined> => {
  try {
    const response = await axios.get(FetchCategoriesAllAPI());
    return response.data.categories;
  } catch (error) {
    console.error("カテゴリの取得に失敗しました", error);
    return undefined;
  }
};

export const fetchCharacters = async (): Promise<Character[] | undefined> => {
  try {
    const response = await axios.get(FetchAllCharactersAPI());

    return response.data.characters;
  } catch (error) {
    console.error("キャラクターの取得に失敗しました", error);
    return undefined;
  }
};

const UserSidebar = async () => {
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
        <SidebarMain links={links} />
      </div>
    </>
  );
};

export default UserSidebar;
