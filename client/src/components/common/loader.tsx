"use client";

import { Oval } from "react-loader-spinner";

type Props = {
  height?: string;
  size?: number;
};

const Loader: React.FC<Props> = ({ height = "h-40", size = 50 }) => {
  return (
    <>
      <div className={`w-full ${height}`}>
        <div className="h-full flex items-center justify-center">
          <Oval color="#17A34A" height={size} width={size} />
        </div>
      </div>
    </>
  );
};

export default Loader;
