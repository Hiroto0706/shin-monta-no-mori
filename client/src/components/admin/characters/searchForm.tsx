"use client";

import { FormEvent, useEffect, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";

const CharactersSearchForm: React.FC = () => {
  const router = useRouter();
  const displayLimit = 3;
  const [title, setTitle] = useState("");

  const searchIllustrations = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const queryParams: { [key: string]: string } = {};

    if (title) {
      queryParams.q = title;
    }

    const queryString = new URLSearchParams(queryParams).toString();
    router.push(`/admin/illustrations?${queryString}`);
    router.refresh();
  };

  return (
    <div className="flex flex-col flex-col-reverse lg:flex-row justify-between">
      <form className="flex flex-wrap" onSubmit={(e) => searchIllustrations(e)}>
        <div className="lg:mr-3 mb-6 lg:mb-3 w-full lg:w-80">
          <input
            type="text"
            placeholder="タイトル検索"
            onChange={(e) => setTitle(e.target.value)}
            className="border-2 border-gray-200 py-2.5 px-4 rounded-md w-full"
          />
        </div>

        <button className="flex justify-center items-center lg:justify-start bg-green-600 text-white rounded-md font-bold py-2.5 border-2 border-green-600 pl-4 pr-3 lg:mb-6 w-full lg:w-auto hover:opacity-70 duration-200">
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
