"use client";

import Image from "next/image";
import { usePathname } from "next/navigation";

const links = [
  {
    href: "/illustrations",
    icon: "/icon/illustration.png",
    icon_active: "/icon/illustration-active.png",

    text: "イラスト",
  },
  {
    href: "/characters",
    icon: "/icon/character.png",
    icon_active: "/icon/character-active.png",
    text: "キャラ",
  },
  {
    href: "/categories",
    icon: "/icon/category.png",
    icon_active: "/icon/category-active.png",
    text: "カテゴリ",
  },
];

function UserSidebar() {
  const pathname = usePathname();

  return (
    <>
      {links.map((link, index) => (
        <li className="my-4" key={index}>
          <a href={link.href} className="flex flex-col items-center">
            <Image
              src={pathname == link.href ? link.icon_active : link.icon}
              alt={`${link.text}アイコン`}
              height={36}
              width={36}
            />
            <span
              className={
                pathname == link.href
                  ? `text-green-600 font-bold`
                  : `text-gray-600` + `text-xs`
              }
            >
              {link.text}
            </span>
          </a>
        </li>
      ))}
    </>
  );
}

export default UserSidebar;
