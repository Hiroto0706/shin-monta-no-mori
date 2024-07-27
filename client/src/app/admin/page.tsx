"use client";

import Image from "next/image";
import Link from "next/link";

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
    <div>
      <h1 className="text-2xl font-bold">管理者画面トップ</h1>

      <div className="my-12">
        <ul
          className="
        grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8
        "
        >
          {links.map((link, index) => (
            <li
              key={index}
              className="
              p-2 border-2 border-gray-200 rounded-xl min-w-[220px] bg-white hover:scale-105 duration-200 cursor-pointer
              "
            >
              <Link href={link.href} className="flex items-center">
                <Image
                  src={link.icon}
                  alt={`${link.text}アイコン`}
                  height={40}
                  width={40}
                />
                <span className="ml-4 text-xl">{link.text}</span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
