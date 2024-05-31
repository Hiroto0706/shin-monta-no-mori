import axios from "axios";
import { Illustration } from "@/types/illustration";
import Illustrations from "@/components/admin/illustrations/illustrations";
import { GetAccessToken, SetBearerToken } from "@/utils/accessToken";
import { Character } from "@/types/character";
import { Category } from "@/types/category";

const fetchIllustrations = async (
  page: number = 0
): Promise<Illustration[]> => {
  const accessToken = GetAccessToken();

  try {
    const response = await axios.get(
      "http://localhost:8080/api/v1/admin/illustrations/list?p=" + page,
      {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      }
    );
    return response.data;
  } catch (error) {
    console.error(error);
    return [];
  }
};

const fetchCharacters = async (page: number = 0): Promise<Character[]> => {
  const accessToken = GetAccessToken();

  try {
    const response = await axios.get(
      "http://localhost:8080/api/v1/admin/characters/list?p=" + page,
      {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      }
    );
    return response.data.characters;
  } catch (error) {
    console.error(error);
    return [];
  }
};

const fetchCategories = async (): Promise<Category[]> => {
  const accessToken = GetAccessToken();

  try {
    const response = await axios.get(
      "http://localhost:8080/api/v1/admin/categories/list",
      {
        headers: {
          Authorization: SetBearerToken(accessToken),
        },
      }
    );
    console.log(response.data.categories);
    return response.data.categories;
  } catch (error) {
    console.error(error);
    return [];
  }
};

export default async function IllustrationsPage({
  searchParams,
}: {
  searchParams: { p: string };
}) {
  const page = searchParams.p ? parseInt(searchParams.p, 10) : 0;
  const illustrations = await fetchIllustrations(page);
  const characters = await fetchCharacters(page);
  const categories = await fetchCategories();
  return (
    <Illustrations
      illustrations={illustrations}
      characters={characters}
      categories={categories}
    />
  );
}
