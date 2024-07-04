export const VerifyAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "auth/verify";
};

export const AuthLoginAPI = (): string => {
  return process.env.NEXT_PUBLIC_BASE_API + "auth/login";
};
