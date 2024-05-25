"use client";

import Image from "next/image";

const links = [
  { href: "/admin", icon: "/icon/home.png", text: "TOP" },
  { href: "/admin/illustrations", icon: "/icon/illustration.png", text: "イラスト" },
  { href: "/admin/characters", icon: "/icon/character.png", text: "キャラ" },
  { href: "/admin/categories", icon: "/icon/category.png", text: "カテゴリ" },
];

function Sidebar() {
  return (
    <div className="w-20 h-full fixed inset-0 z-30 border-r-2 border-gray-200 bg-gray-50">
      <div className="pt-16">
        <ul className="flex flex-col items-center">
          {links.map((link, index) => (
            <li className="my-4" key={index}>
              <a href={link.href} className="flex flex-col items-center">
                <Image
                  src={link.icon}
                  alt={`${link.text}アイコン`}
                  height={36}
                  width={36}
                />
                <span className="text-gray-600 text-xs">{link.text}</span>
              </a>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default Sidebar;
