"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import React, { FormEvent, useState } from "react";

type Props = {
  maxWidth: number;
  query?: string;
  placeHolder?: string;
  addClass?: string;
};

const SearchBox: React.FC<Props> = ({
  query = "",
  maxWidth,
  placeHolder = "いらすとを検索する",
  addClass,
}) => {
  const router = useRouter();
  const [name, setName] = useState(query);

  const searchIllustrations = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (name != "") {
      router.push(`/illustrations/search/${name}`);
    } else {
      router.push("/illustrations");
    }
  };

  return (
    <form
      className={`flex justify-between w-full md:max-w-[${maxWidth}px] mx-auto border-gray-200 rounded-full bg-white ${addClass}`}
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
