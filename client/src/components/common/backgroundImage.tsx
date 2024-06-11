"use client";

import { usePathname } from "next/navigation";

function BackgroundImage() {
  const pathname = usePathname();

  return (
    <>
      <div className={pathname != "/" ? "background-image" : ""}></div>
    </>
  );
}

export default BackgroundImage;
