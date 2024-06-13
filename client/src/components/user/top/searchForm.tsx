"use client";

import { FormEvent, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Category } from "@/types/category";

type Props = {
  categories: Category[];
};

const SearchFormTop: React.FC<Props> = ({ categories }) => {
  const router = useRouter();
  const [name, setName] = useState("");

  const searchCharacters = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const queryParams: { [key: string]: string } = {};

    if (name) {
      queryParams.q = name;
    }

    const queryString = new URLSearchParams(queryParams).toString();
    router.push(`/admin/categories?${queryString}`);
    router.refresh();
  };

  return (
    <>
      <p className="text-lg mb-4 font-bold">
        もんたの森はゆるーくてゆーもある無料イラストサイトです
      </p>
      <form
        className="flex justify-between w-full max-w-[550px] mx-auto border-gray-200 rounded-md bg-white mb-2"
        onSubmit={(e) => searchCharacters(e)}
      >
        <div className="w-full">
          <input
            className="pl-2 w-full h-full rounded-l-md"
            type="text"
            placeholder="いらすとを検索する"
            onChange={(e) => setName(e.target.value)}
          />
        </div>
        <button className="p-2">
          <Image
            src="/icon/search_gray.png"
            alt="searchアイコン"
            width={24}
            height={24}
          />
        </button>
      </form>
      <div className="flex flex-wrap items-center w-full max-w-[550px] mb-12">
        <span className="text-sm my-1">おすすめかてごり : </span>
        {categories.slice(0, 5).map((category) => (
          <div
            key={category.ParentCategory.id}
            className="text-gray-600 text-sm ml-2 my-1 py-1 px-2 rounded-lg bg-gray-200 hover:bg-gray-300 duration-200 cursor-pointer shadow"
          >
            # {category.ParentCategory.name}
          </div>
        ))}
      </div>
    </>
  );
};

export default SearchFormTop;
