"use client";

import axios from "axios";
import { Illustration } from "@/types/illustration";
import IllustrationCard from "../illustrationCard";
import { useEffect, useState } from "react";
import { FetchIllustrationsResponse } from "@/types/user/illustration";
import { FetchRandomIllustrationsAPI } from "@/api/user/illustration";
import Loader from "@/components/common/loader";

interface Props {
  exclusion_id: number;
}

const fetchRandomIllustrations = async (
  exclusion_id: number
): Promise<FetchIllustrationsResponse> => {
  try {
    const response = await axios.get(FetchRandomIllustrationsAPI(exclusion_id));

    return response.data;
  } catch (error) {
    console.error(error);
    return { illustrations: [] };
  }
};

const RandomIllustrations: React.FC<Props> = ({ exclusion_id }) => {
  const [illustrations, setIllustrations] = useState<Illustration[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      const res = await fetchRandomIllustrations(exclusion_id);
      setIllustrations(res.illustrations);
    };

    fetchData();
  }, [exclusion_id]);

  return (
    <>
      {illustrations.length > 0 ? (
        <>
          <div className="mt-8 grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 lg:grid-cols-4">
            {illustrations.map((illustration) => (
              <IllustrationCard
                key={illustration.Image.id}
                illustration={illustration}
              />
            ))}
          </div>
        </>
      ) : (
        <>
          <Loader />
        </>
      )}
    </>
  );
};

export default RandomIllustrations;
