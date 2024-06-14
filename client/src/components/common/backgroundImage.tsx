"use client";

import { usePathname } from "next/navigation";

function BackgroundImage() {
  const pathname = usePathname();

  return (
    <>{pathname != "/" ? <div className="background-image"></div> : <></>}</>
  );
}

export default BackgroundImage;
