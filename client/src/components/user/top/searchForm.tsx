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

  const searchIllustrations = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const queryParams: { [key: string]: string } = {};

    if (name) {
      queryParams.q = name;
    }

    const queryString = new URLSearchParams(queryParams).toString();
    router.push(`/illustrations?${queryString}`);
  };

  return (
    <>
      <p className="text-md md:text-xl mb-2 md:mb-4 font-bold px-4 mb:px-0">
        もんたの森はゆるーくてゆーもある無料イラストサイトです
      </p>
      <div className="w-full px-4">
        <form
          className="flex justify-between w-full md:max-w-[550px] mx-auto border-gray-200 rounded-md bg-white mb-2"
          onSubmit={(e) => searchIllustrations(e)}
        >
          <div className="w-full">
            <input
              className="pl-2 w-full h-full rounded-l-md text-gray-600"
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
      </div>
      <div className="flex flex-wrap items-center w-full md:max-w-[550px] px-4 md:px-0 mb-4 md:mb-12">
        <span className="text-sm my-1">おすすめかてごり : </span>
        {categories.slice(0, 5).map((category) => (
          <a
            href=""
            key={category.ParentCategory.id}
            className="text-gray-600 text-sm ml-2 my-1 py-1 px-2 rounded-lg border bg-white hover:bg-gray-200 duration-200 cursor-pointer shadow"
          >
            # {category.ParentCategory.name}
          </a>
        ))}
      </div>
    </>
  );
};

export default SearchFormTop;
