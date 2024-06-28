"use client";

import axios from "axios";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { AuthLoginAPI } from "@/api/auth";

const LoginPage = () => {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordCheck, setPasswordCheck] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (password !== passwordCheck) {
      setError("パスワードが一致しません");
      return;
    }

    const formData = {
      email: email,
      password: password,
    };
    try {
      const response = await axios.post(AuthLoginAPI(), formData, {
        withCredentials: true,
      });
      if (response.status == 200) {
        const accessToken = response.data.access_token;
        // クッキーを特定の日付まで有効にする
        const expirationDate = new Date();
        expirationDate.setDate(expirationDate.getDate() + 2);
        Cookies.set("access_token", accessToken, { expires: expirationDate });

        router.push("/admin");
      } else {
        router.push("/login");
      }
    } catch (error: any) {
      if (error.code === "ERR_NETWORK") {
        setError("予期せぬエラーが発生しました");
      } else {
        setError("メールまたはパスワードが間違っています");
      }
      return;
    }
  };

  return (
    <>
      <div className="flex flex-col items-center justify-center p-12">
        <div className="p-6 md:p-12 border-gray-200 border-2 rounded-xl bg-white bg-opacity-70 w-2/4 min-w-[350px] md:min-w-[450px]">
          <h1 className="text-2xl font-bold mb-8">管理者ログインフォーム</h1>
          <form onSubmit={handleSubmit} className="w-full">
            <div className="mb-6">
              <label htmlFor="email" className="block text-gray-700">
                Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full p-2 border-2 border-gray-200 rounded mt-1"
                required
              />
            </div>
            <div className="mb-6">
              <label htmlFor="password" className="block text-gray-700">
                Password
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full p-2 border-2 border-gray-200 rounded mt-1"
                required
              />
            </div>
            <div className="mb-6">
              <label htmlFor="passwordCheck" className="block text-gray-700">
                Password Check
              </label>
              <input
                id="passwordCheck"
                type="password"
                value={passwordCheck}
                onChange={(e) => setPasswordCheck(e.target.value)}
                className="w-full p-2 border-2 border-gray-200 rounded mt-1"
                required
              />
            </div>
            {error !== "" && <p className="text-red-700">{error}</p>}
            <button
              type="submit"
              className="w-full bg-green-600 text-white p-2 rounded mt-6 hover:opacity-70 duration-200"
            >
              <span className="text-2xl">Login</span>
            </button>
          </form>
        </div>
      </div>
    </>
  );
};

export default LoginPage;
