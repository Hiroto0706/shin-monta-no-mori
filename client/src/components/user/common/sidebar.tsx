"use client";

import Image from "next/image";
import { usePathname } from "next/navigation";

function UserSidebar() {
  const pathname = usePathname();

  const links = [
    // {
    //   href: "/illustrations",
    //   icon: "/icon/illustration.png",
    //   icon_active: "/icon/illustration-active.png",
    //   text: "イラスト",
    //   sublinks: [
    //     {
    //       href: "/illustrations/new",
    //     },
    //     {
    //       href: "/illustrations/edit",
    //     },
    //   ],
    // },
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

  return (
    <>
      <div className="duration-200 transform -translate-x-full md:transform-none">
        <div className="w-16 h-full fixed inset-0 z-30 bg-gray-100 border-r border-gray-200">
          <div className="pt-16">
            <ul className="flex flex-col items-center mt-2">
              {links.map((link, index) => {
                // const isActive =
                //   pathname === link.href ||
                //   link.sublinks.some((sublink) => pathname.startsWith(sublink.href));

                return (
                  <li
                    className="mt-2 p-1 w-14 duration-200 rounded-lg hover:bg-gray-200 cursor-pointer"
                    key={index}
                  >
                    <div className="flex flex-col items-center">
                      <Image
                        src={link.icon}
                        alt={`${link.text}アイコン`}
                        height={28}
                        width={28}
                      />
                      <span className={`text-xs text-gray-600`}>
                        {link.text}
                      </span>
                    </div>
                  </li>
                );
              })}
            </ul>
          </div>
        </div>

        <div className="w-48 h-full bg-gray-50 border-r border-gray-200 fixed top-0 left-0">
          a
        </div>
      </div>
    </>
  );
}

export default UserSidebar;
