"use client";

import { ChildCategory } from "@/types/category";
import SearchBox from "../common/searchBox";
import Link from "next/link";

type Props = {
  child_categories: ChildCategory[];
};

const SearchFormTop: React.FC<Props> = ({ child_categories }) => {
  return (
    <>
      <p className="text-md md:text-xl mb-2 md:mb-4 font-bold px-4 mb:px-0">
        もんたの森はゆるーくてゆーもあな無料イラストサイトです
      </p>
      <div className="w-full px-4">
        <SearchBox maxWidth={600} addClass="mb-2" />
      </div>
      <div className="flex flex-wrap items-center w-full md:max-w-[600px] px-4 md:px-0">
        {child_categories.slice(0, 5).map((child_category) => (
          <Link
            href={`illustrations/category/${child_category.id}`}
            key={child_category.id}
            className="text-gray-600 text-sm mr-1 my-1 py-1 px-2 rounded-full border bg-white hover:bg-gray-200 duration-200 cursor-pointer shadow"
          >
            # {child_category.name}
          </Link>
        ))}
      </div>
    </>
  );
};

export default SearchFormTop;
