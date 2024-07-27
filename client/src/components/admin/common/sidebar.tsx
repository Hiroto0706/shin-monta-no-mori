"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";

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
      },
      {
        href: "/admin/illustrations/edit",
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
        href: "/admin/characters/new",
      },
      {
        href: "/admin/characters/edit",
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
        href: "/admin/categories/parent/",
      },
      {
        href: "/admin/categories/child/",
      },
    ],
  },
];

function AdminSidebar() {
  const pathname = usePathname();

  return (
    <>
      <div className="w-16 h-full fixed inset-0 z-30 bg-gray-100">
        <div className="pt-16">
          <ul className="flex flex-col items-center mt-2">
            {links.map((link, index) => {
              const isActive =
                pathname === link.href ||
                link.sublinks.some((sublink) =>
                  pathname.startsWith(sublink.href)
                );

              return (
                <li
                  className={`mt-2 p-1 w-14 duration-200 rounded-lg ${
                    isActive ? "" : "hover:bg-gray-200"
                  }`}
                  key={index}
                >
                  <Link href={link.href} className="flex flex-col items-center">
                    <Image
                      src={isActive ? link.icon_active : link.icon}
                      alt={`${link.text}アイコン`}
                      height={28}
                      width={28}
                    />
                    <span
                      className={`text-xs
                  ${isActive ? `text-green-600 font-bold` : `text-gray-600`}
                `}
                    >
                      {link.text}
                    </span>
                  </Link>
                </li>
              );
            })}
          </ul>
        </div>
      </div>
    </>
  );
}

export default AdminSidebar;
