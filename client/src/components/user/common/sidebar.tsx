"use client";

import Image from "next/image";
import { usePathname } from "next/navigation";

const links = [
  {
    href: "/illustrations",
    icon: "/icon/illustration.png",
    icon_active: "/icon/illustration-active.png",
    text: "イラスト",
    sublinks: [
      {
        href: "/illustrations/new",
      },
      {
        href: "/illustrations/edit",
      },
    ],
  },
  {
    href: "/characters",
    icon: "/icon/character.png",
    icon_active: "/icon/character-active.png",
    text: "キャラ",
    sublinks: [
      {
        href: "/characters/new",
      },
      {
        href: "/characters/edit",
      },
    ],
  },
  {
    href: "/categories",
    icon: "/icon/category.png",
    icon_active: "/icon/category-active.png",
    text: "カテゴリ",
    sublinks: [
      {
        href: "/categories/parent/",
      },
      {
        href: "/categories/child/",
      },
    ],
  },
];

function UserSidebar() {
  const pathname = usePathname();

  return (
    <>
      {links.map((link, index) => {
        const isActive =
          pathname === link.href ||
          link.sublinks.some((sublink) => pathname.startsWith(sublink.href));

        return (
          <li
            className={`mt-2 p-1 w-14 duration-200 rounded-lg ${
              isActive ? "" : "hover:bg-gray-200"
            }`}
            key={index}
          >
            <a href={link.href} className="flex flex-col items-center">
              <Image
                src={isActive ? link.icon_active : link.icon}
                alt={`${link.text}アイコン`}
                height={36}
                width={36}
              />
              <span
                className={`text-xs
                  ${isActive ? `text-green-600 font-bold` : `text-gray-600`}
                `}
              >
                {link.text}
              </span>
              {/* {isActive && link.sublinks.length > 0 && (
                <ul className="mt-4 pt-4 bg-gray-200 w-20 flex flex-col items-center">
                  {link.sublinks.map((sublink, subIndex) => (
                    <li key={subIndex} className="hover:opacity-50 mb-4">
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
                        <span className="mt-1">{sublink.text}</span>
                      </a>
                    </li>
                  ))}
                </ul>
              )} */}
            </a>
          </li>
        );
      })}
    </>
  );
}

export default UserSidebar;
