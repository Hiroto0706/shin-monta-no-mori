import { cookies } from "next/headers";

export const GetAccessToken = () => {
  return cookies().get("access_token")?.value;
};

export const SetBearerToken = (accessToken: string | undefined) => {
  return `Bearer ${accessToken}`;
};
