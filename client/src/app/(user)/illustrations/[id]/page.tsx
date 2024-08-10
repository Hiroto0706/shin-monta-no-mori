import axios from "axios";
import { GetIllustrationAPI } from "@/api/user/illustration";
import { GetIllustrationResponse } from "@/types/user/illustration";
import IllustrationDetailTemplate from "@/components/user/illustrations/detail/illustrationDetailTemplate";

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

const IllustrationDetailPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const getIllustrationRes = await getIllustration(params.id);

  return (
    <>
      <div className="w-full max-w-[1100px] m-auto">
        <IllustrationDetailTemplate
          id={params.id}
          illustration={getIllustrationRes.illustration}
        />
      </div>
    </>
  );
};

export default IllustrationDetailPage;
