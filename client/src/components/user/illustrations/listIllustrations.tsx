"use client";

import { Illustration } from "@/types/illustration";
import React from "react";

interface Props {
  illustrations: Illustration[];
}

const ListIllustrations: React.FC<Props> = ({ illustrations }) => {
  return (
    <>
      {illustrations.map((illustration) => (
        <div key={illustration.Image.id}>{illustration.Image.title}</div>
      ))}
    </>
  );
};

export default ListIllustrations;
