import { getServerAccessToken } from "@/utils/accessToken/server";

const CreateChildCategoryPage = async () => {
  const accessToken = getServerAccessToken();

  return <>new child</>;
};

export default CreateChildCategoryPage;
