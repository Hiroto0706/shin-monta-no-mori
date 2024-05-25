"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

export default function TOP() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordCheck, setPasswordCheck] = useState("");

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (password !== passwordCheck) {
      alert("パスワードが一致しません");
      return;
    }

    // フォームの送信処理
    console.log("Username:", username);
    console.log("Email:", email);
    console.log("Password:", password);
    console.log("Password Check:", passwordCheck);

    // TODO: ここでログインのロジックを追加
    router.push("/admin");
  };

  return (
    <div className="m-0 md:m-12 p-12 border-gray-200 border-2 rounded-xl bg-white bg-opacity-50">
      <h1 className="text-2xl font-bold mb-8">管理者ログインフォーム</h1>
      <form onSubmit={handleSubmit} className="w-full">
        <div className="mb-6">
          <label htmlFor="username" className="block text-gray-700">
            Username
          </label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full p-2 border-2 border-gray-200 rounded mt-1"
            required
          />
        </div>
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
        <div className="mb-2">
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
        <div className="mb-12">
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
        <button
          type="submit"
          className="w-full bg-green-600 text-white p-2 rounded"
        >
          <span className="text-2xl">Login</span>
        </button>
      </form>
    </div>
  );
}
