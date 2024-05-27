import axios from "axios";
import { Illustration } from "@/types/illustration";
import Illustrations from "@/components/admin/illustrations/illustrations";
import { GetAccessToken, SetBearerToken } from "@/utils/accessToken";

const fetchIllustrations = async (
  page: number = 0
): Promise<Illustration[]> => {
  const accessToken = GetAccessToken();

  try {
    const response = await axios.get(
      "http://localhost:8080/api/v1/admin/illustrations/list/?p=" + page,
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

export default async function IllustrationsPage({
  searchParams,
}: {
  searchParams: { p: string };
}) {
  const page = searchParams.p ? parseInt(searchParams.p, 10) : 0;
  console.log("page", page);
  const illustrations = await fetchIllustrations(page);
  return <Illustrations illustrations={illustrations} />;
}
