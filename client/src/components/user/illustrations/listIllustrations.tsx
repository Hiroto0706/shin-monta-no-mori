"use client";

import { Illustration } from "@/types/illustration";
import Image from "next/image";
import React from "react";

interface Props {
  illustrations: Illustration[];
}

const ListIllustrations: React.FC<Props> = ({ illustrations }) => {
  return (
    <>
      <div className="mt-8">
        <div className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6">
          {illustrations.map((illustration) => (
            <div key={illustration.Image.id} className="group cursor-pointer">
              <div className="mb-2 p-4 border-2 border-gray-200 rounded-xl bg-white">
                <Image
                  className="group-hover:scale-110 duration-200 image"
                  src={illustration.Image.original_src}
                  alt={illustration.Image.title}
                  fill
                />
              </div>
              <span className="group-hover:text-green-600 group-hover:font-bold duration-200">
                {illustration.Image.title}
              </span>
            </div>
          ))}
        </div>
      </div>
    </>
  );
};

export default ListIllustrations;
