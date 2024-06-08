import axios from "axios";
import { getServerAccessToken } from "@/utils/accessToken/server";

const EditChildCategoryPage = async ({
  params,
}: {
  params: { id: number };
}) => {
  const accessToken = getServerAccessToken();

  return <>child : id = {params.id}</>;
};

export default EditChildCategoryPage;
