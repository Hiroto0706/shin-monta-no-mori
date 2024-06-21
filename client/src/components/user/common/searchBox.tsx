"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import React, { FormEvent, useState } from "react";

type Props = {
  query?: string;
  maxWidth?: string;
  placeHolder?: string;
  addClass?: string;
};

const SearchBox: React.FC<Props> = ({
  query = "",
  maxWidth = "550px",
  placeHolder = "いらすとを検索する",
  addClass,
}) => {
  const router = useRouter();
  const [name, setName] = useState(query);

  const searchIllustrations = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const queryParams: { [key: string]: string } = {};

    if (name) {
      queryParams.q = name;
    }

    const queryString = new URLSearchParams(queryParams).toString();
    router.push(`/illustrations?${queryString}`);
    router.refresh();
  };

  return (
    <form
      className={`flex justify-between w-full md:max-w-[${maxWidth}] mx-auto border-gray-200 rounded-full bg-white ${addClass}`}
      onSubmit={(e) => searchIllustrations(e)}
    >
      <div className="w-full">
        <input
          className="pl-2 w-full h-full rounded-full text-gray-600"
          type="text"
          placeholder={placeHolder}
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
  );
};

export default SearchBox;
