"use client";

import Image from "next/image";
import { usePathname } from "next/navigation";
import { useState } from "react";

const links = [
  {
    href: "/admin",
    icon: "/icon/home.png",
    icon_active: "/icon/home-active.png",
    text: "TOP",
    sublinks: [],
  },
  {
    href: "/admin/illustrations",
    icon: "/icon/illustration.png",
    icon_active: "/icon/illustration-active.png",
    text: "イラスト",
    sublinks: [
      {
        href: "/admin/illustrations/new",
        icon: "/icon/create.png",
        text: "新規作成",
      },
      {
        href: "/admin/illustrations/",
        icon: "/icon/list.png",
        text: "一覧",
      },
    ],
  },
  {
    href: "/admin/characters",
    icon: "/icon/character.png",
    icon_active: "/icon/character-active.png",
    text: "キャラ",
    sublinks: [
      {
        href: "/admin/illustrations/new",
        icon: "/icon/create.png",
        text: "新規作成",
      },
      {
        href: "/admin/illustrations/",
        icon: "/icon/list.png",
        text: "一覧",
      },
    ],
  },
  {
    href: "/admin/categories",
    icon: "/icon/category.png",
    icon_active: "/icon/category-active.png",
    text: "カテゴリ",
    sublinks: [
      {
        href: "/admin/illustrations/new",
        icon: "/icon/create.png",
        text: "新規作成",
      },
      {
        href: "/admin/illustrations/",
        icon: "/icon/list.png",
        text: "一覧",
      },
    ],
  },
];

function AdminSidebar() {
  const pathname = usePathname();

  return (
    <>
      {links.map((link, index) => (
        <li className="mt-4" key={index}>
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
            {pathname === link.href && link.sublinks.length > 0 && (
              <ul className="mt-2 pt-4 bg-gray-200 w-20 flex flex-col items-center">
                {link.sublinks.map((sublink, subIndex) => (
                  <li key={subIndex} className="hover:opacity-50 mb-4 ">
                    <a
                      href={sublink.href}
                      className="text-gray-600 text-xs flex flex-col items-center"
                    >
                      <Image
                        src={sublink.icon}
                        alt={`${sublink.text}アイコン`}
                        height={32}
                        width={32}
                      />
                      <span>{sublink.text}</span>
                    </a>
                  </li>
                ))}
              </ul>
            )}
          </a>
        </li>
      ))}
    </>
  );
}

export default AdminSidebar;
