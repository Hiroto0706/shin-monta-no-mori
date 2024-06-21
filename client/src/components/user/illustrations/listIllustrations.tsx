"use client";

import { useEffect, useState } from "react";
import { Illustration } from "@/types/illustration";
import Image from "next/image";
import InfiniteScroll from "react-infinite-scroller";
import axios from "axios";
import {
  FetchIllustrationsAPI,
  SearchIllustrationsAPI,
} from "@/api/user/illustration";

interface Props {
  initialIllustrations: Illustration[];
  query: string | null;
}

const ListIllustrations: React.FC<Props> = ({
  initialIllustrations,
  query,
}) => {
  const [illustrations, setIllustrations] =
    useState<Illustration[]>(initialIllustrations);
  const [hasMore, setHasMore] = useState(true);
  const [page, setPage] = useState(1);

  const fetchMoreIllustrations = async (page: number) => {
    try {
      const response = query
        ? await axios.get(SearchIllustrationsAPI(page, query))
        : await axios.get(FetchIllustrationsAPI(page));
      const newIllustrations = response.data.illustrations;

      if (newIllustrations.length === 0) {
        setHasMore(false);
      } else {
        setIllustrations((prevIllustrations) => [
          ...prevIllustrations,
          ...newIllustrations,
        ]);
        setPage(page + 1);
      }
    } catch (error) {
      console.error("Failed to fetch more illustrations", error);
    }
  };

  return (
    <>
      <InfiniteScroll
        pageStart={0}
        loadMore={fetchMoreIllustrations}
        hasMore={hasMore}
        loader={<div key={0}>Loading...</div>}
      >
        <div className="mt-8 grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6">
          {illustrations.map((illustration) => (
            <a
              key={illustration.Image.id}
              href={`/illustrations/${illustration.Image.id}`}
              className="group cursor-pointer"
            >
              <div
                className="mb-2 border-2 border-gray-200 rounded-xl bg-white relative w-full overflow-hidden"
                style={{ paddingTop: "100%" }}
              >
                <div className="absolute inset-0 m-4">
                  <Image
                    className="group-hover:scale-110 duration-200 absolute top-0 left-0 w-full h-full object-cover"
                    src={illustration.Image.original_src}
                    alt={illustration.Image.title}
                    fill
                  />
                </div>
              </div>
              <span className="break-words group-hover:text-green-600 group-hover:font-bold duration-200">
                {illustration.Image.title}
              </span>
            </a>
          ))}
        </div>
      </InfiniteScroll>
    </>
  );
};

export default ListIllustrations;
