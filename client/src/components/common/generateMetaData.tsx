import axios from "axios";
import type { Metadata } from "next";
import { GetIllustrationAPI } from "@/api/user/illustration";
import { GetIllustrationResponse } from "@/types/user/illustration";

const getIllustration = async (
  id: number
): Promise<GetIllustrationResponse> => {
  try {
    const response = await axios.get(GetIllustrationAPI(id));

    return response.data;
  } catch (error) {
    console.error(error);
    return { illustration: null };
  }
};

export async function generateMetaData({ params }: { params: { id: number } }) {
  const response = await getIllustration(params.id);

  return {
    title: response.illustration?.Image.title,
    description: `『${response.illustration?.Image.title}』だよ。もんたの森では他にも可愛くてクセのある画像がたくさんあるよ。`,
    icons: {
      icon: [
        {
          url: "/src/app/favicon.ico",
          href: "/src/app/favicon.ico",
        },
      ],
    },
    openGraph: {
      images: [
        {
          url:
            response.illustration?.Image.original_src != undefined
              ? response.illustration.Image.original_src
              : "https://storage.googleapis.com/shin-monta-no-mori/montanomori-top-image.png",
          alt:
            response.illustration?.Image.original_src != undefined
              ? response.illustration.Image.title
              : "もんたの森のイメージ画像",
        },
      ],
    },
  };
}
