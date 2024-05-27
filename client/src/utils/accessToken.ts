import { cookies } from "next/headers";

export const getAccessToken = () => {
  return cookies().get("access_token")?.value;
};
