"use client";

import { useState } from "react";
import { Illustration } from "@/types/illustration";
import Image from "next/image";
import InfiniteScroll from "react-infinite-scroller";
import axios from "axios";
import {
  FetchIllustrationsAPI,
  FetchIllustrationsByCategoryAPI,
  FetchIllustrationsByCharacterAPI,
  SearchIllustrationsAPI,
} from "@/api/user/illustration";
import Loader from "@/components/common/loader";
import IllustrationCard from "./illustrationCard";

interface Props {
  initialIllustrations: Illustration[];
  fetchType: {
    query?: string;
    categoryID?: number;
    characterID?: number;
  };
}

const ListIllustrations: React.FC<Props> = ({
  initialIllustrations,
  fetchType,
}) => {
  const [illustrations, setIllustrations] =
    useState<Illustration[]>(initialIllustrations);
  const [hasMore, setHasMore] = useState(
    !(
      initialIllustrations.length <
      Number(process.env.NEXT_PUBLIC_IMAGE_FETCH_LIMIT)
    )
  );
  const [page, setPage] = useState(1);

  const fetchUrlGenerator = (page: number) => {
    if (fetchType.query) {
      return SearchIllustrationsAPI(page, fetchType.query);
    } else if (fetchType.categoryID) {
      return FetchIllustrationsByCategoryAPI(fetchType.categoryID, page);
    } else if (fetchType.characterID) {
      return FetchIllustrationsByCharacterAPI(fetchType.characterID, page);
    } else {
      return FetchIllustrationsAPI(page);
    }
  };

  const handleFetchMore = async () => {
    try {
      const url = fetchUrlGenerator(page);
      const response = await axios.get(url);
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
        loadMore={handleFetchMore}
        hasMore={hasMore}
        loader={
          <div key={0}>
            <Loader />
          </div>
        }
      >
        <div className="mt-8 grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 lg:grid-cols-4">
          {illustrations.map((illustration) => (
            <IllustrationCard
              key={illustration.Image.id}
              illustration={illustration}
            />
          ))}
        </div>
      </InfiniteScroll>
    </>
  );
};

export default ListIllustrations;
