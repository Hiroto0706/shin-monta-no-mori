import { cookies } from "next/headers";

// サーバーサイドでのトークン取得
export const getServerAccessToken = (): string | undefined => {
  return cookies().get("access_token")?.value;
};
