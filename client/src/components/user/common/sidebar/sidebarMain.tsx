"use client";

import Image from "next/image";
import SidebarSub from "@/components/user/common/sidebar/sidebarSub";
import { Character } from "@/types/character";
import { useEffect, useState } from "react";
import { Category } from "@/types/category";
import {
  fetchCategories,
  fetchCharacters,
} from "@/components/user/common/sidebar/sidebar";
import useSidebarStore from "@/store/sidebar";

type Props = {
  links: {
    id: number;
    href: string;
    icon: string;
    icon_active: string;
    text: string;
  }[];
};

const SidebarMain: React.FC<Props> = ({ links }) => {
  const [selectedLink, setSelectedLink] = useState<number>(0);
  const [categories, setCategories] = useState<Category[] | undefined>(
    undefined
  );
  const [characters, setCharacters] = useState<Character[] | undefined>(
    undefined
  );
  const { isShow, toggleIsShow } = useSidebarStore();

  const setSidebarStatus = (id: number) => {
    setSelectedLink(id);
    toggleIsShow(true);
  };

  useEffect(() => {
    if (selectedLink === 0 && categories == undefined) {
      fetchCategories().then((data) => {
        if (data) {
          setCategories(data);
        }
      });
    }

    if (selectedLink === 1 && characters == undefined) {
      fetchCharacters().then((data) => {
        if (data) {
          setCharacters(data);
        }
      });
    }
  }, [selectedLink]);

  return (
    <>
      <div className="w-16 h-full fixed inset-0 z-30 bg-gray-100 border-r border-gray-200">
        <div className="pt-16 z-30">
          <ul className="flex flex-col items-center mt-2">
            {links.map((link, index) => {
              return (
                <li
                  key={index}
                  className="mt-2 p-1 w-14 duration-200 rounded-lg hover:bg-gray-200 cursor-pointer"
                  onClick={() => setSidebarStatus(link.id)}
                >
                  <div className="flex flex-col items-center">
                    <Image
                      src={
                        selectedLink == link.id && isShow
                          ? link.icon_active
                          : link.icon
                      }
                      alt={`${link.text}アイコン`}
                      height={28}
                      width={28}
                    />
                    <span
                      className={`text-xs text-gray-600 ${
                        selectedLink == link.id && isShow
                          ? "text-green-600 font-bold"
                          : ""
                      }`}
                    >
                      {link.text}
                    </span>
                  </div>
                </li>
              );
            })}
          </ul>
        </div>
      </div>

      <SidebarSub
        links={links}
        selectedLink={selectedLink}
        characters={characters}
        categories={categories}
      />
    </>
  );
};

export default SidebarMain;
