import axios from "axios";
import { Illustration } from "@/types/illustration";
import Illustrations from "@/components/admin/illustrations/illustrations";
import { getAccessToken } from "@/utils/accessToken";

const fetchIllustrations = async (): Promise<Illustration[]> => {
  const accessToken = getAccessToken();

  try {
    const response = await axios.get(
      "http://localhost:8080/api/v1/admin/illustrations/list/?p=0",
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
    return response.data;
  } catch (error) {
    console.error(error);
    return [];
  }
};

export default async function IllustrationsPage() {
  const illustrations = await fetchIllustrations();
  return <Illustrations illustrations={illustrations} />;
}
