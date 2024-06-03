// クライアントサイドでのトークン取得
export  const getClientAccessToken = (): string | undefined => {
  const match = document.cookie.match(new RegExp("(^| )access_token=([^;]+)"));
  return match ? match[2] : undefined;
};