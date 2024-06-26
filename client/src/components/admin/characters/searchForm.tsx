"use client";

import { FormEvent, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";

const CharactersSearchForm: React.FC = () => {
  const router = useRouter();
  const [name, setName] = useState("");

  const searchCharacters = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const queryParams: { [key: string]: string } = {};

    if (name) {
      queryParams.q = name;
    }

    const queryString = new URLSearchParams(queryParams).toString();
    router.push(`/admin/characters?${queryString}`);
    router.refresh();
  };

  return (
    <div className="flex flex-col flex-col-reverse lg:flex-row justify-between">
      <form className="flex flex-wrap" onSubmit={(e) => searchCharacters(e)}>
        <div className="lg:mr-3 mb-6 lg:mb-3 w-full lg:w-80">
          <input
            type="text"
            placeholder="名前検索"
            onChange={(e) => setName(e.target.value)}
            className="border-2 border-gray-200 py-2.5 px-4 rounded-md w-full"
          />
        </div>

        <button className="flex justify-center items-center lg:justify-start bg-green-600 text-white rounded-md font-bold py-2.5 pl-4 pr-3 lg:mb-6 w-full lg:w-auto hover:opacity-70 duration-200">
          <span className="mr-1">検索</span>
          <Image
            src="/icon/search.png"
            alt="searchアイコン"
            width={24}
            height={24}
          />
        </button>
      </form>
    </div>
  );
};

export default CharactersSearchForm;
