"use client";

import { Illustration } from "@/types/illustration";
import SearchBox from "./searcBox";
import { Character } from "@/types/character";
import { Category } from "@/types/category";
import ListTable from "@/components/admin/illustrations/listTable";

type Props = {
  illustrations: Illustration[];
  characters: Character[];
  categories: Category[];
};

const Illustrations: React.FC<Props> = ({
  illustrations,
  characters,
  categories,
}) => {
  return (
    <div>
      <a
        href="illustrations/new"
        className="flex items-center bg-white hover:bg-green-600 border-2 border-green-600 text-green-600 hover:text-white rounded-lg py-2 font-bold mb-6 ml-auto w-full lg:w-36 justify-center duration-200"
      >
        + イラスト追加
      </a>

      <SearchBox
        illustrations={illustrations}
        characters={characters}
        categories={categories}
      />

      <ListTable illustrations={illustrations} />
    </div>
  );
};

export default Illustrations;
