import { getServerAccessToken } from "@/utils/accessToken/server";

const CreateParentCategoryPage = async () => {
  const accessToken = getServerAccessToken();

  return <>new parent</>;
};

export default CreateParentCategoryPage;
