"use client";

import { Illustration } from "@/types/illustration";
import SearchBox from "./searcBox";
import { Character } from "@/types/character";
import { Category } from "@/types/category";

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

      <div className="my-12">
        <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
          {illustrations.map((illustration, index) => (
            <li key={index} className="p-2 border-2 border-gray-200 rounded-xl">
              <div className="flex items-center">
                <span className="ml-4 font-bold text-2xl">
                  {illustration.Image.title}
                </span>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default Illustrations;
