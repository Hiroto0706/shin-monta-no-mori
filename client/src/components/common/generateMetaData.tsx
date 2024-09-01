import axios from "axios";
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

  const title = response.illustration?.Image.title;
  const description = `『${title}』だよ。もんたの森では他にも可愛くてクセのある画像がたくさんあるよ。`;
  const imageUrl =
    response.illustration?.Image.original_src != undefined
      ? response.illustration.Image.original_src
      : "/site-image.png";
  const imageAlt =
    response.illustration?.Image.original_src != undefined
      ? title
      : "もんたの森のイメージ画像";
  return {
    title,
    description,
    openGraph: {
      images: [
        {
          url: imageUrl,
          alt: imageAlt,
        },
      ],
    },
    twitter: {
      card: "summary",
      title,
      description,
      image: imageUrl,
      image_alt: imageAlt,
    },
  };
}
