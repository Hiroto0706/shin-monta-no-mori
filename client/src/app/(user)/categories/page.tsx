"use client";

import Image from "next/image";

const links = [
  {
    href: "/admin/illustrations",
    icon: "/icon/illustration.png",
    text: "イラスト",
  },
  { href: "/admin/characters", icon: "/icon/character.png", text: "キャラ" },
  { href: "/admin/categories", icon: "/icon/category.png", text: "カテゴリ" },
];

export default function TOP() {
  return (
    <>
      <h1 className="text-4xl font-bold">管理者ログインフォーム</h1>

      <div className="my-12">
        <input />
      </div>
    </>
  );
}
