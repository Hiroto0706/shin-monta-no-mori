"use client";

import axios from "axios";
import { useState, useEffect } from "react";
// import FetchData from "@/utils/helloWorld";

const FetchData = async () => {
  try {
    const response = await axios.get("http://localhost:8080/");
    return response.data;
  } catch (error) {
    console.error(error);
    return null;
  }
};

export default function Home() {
  const [data, setData] = useState();

  useEffect(() => {
    const fetchData = async () => {
      const result = await FetchData();
      setData(result);
    };
    fetchData();
  }, []);

  return (
    <main>
      <h1 className="m-10 text-4xl text-red-700">Hello World from client!!</h1>
      <h1 className="m-10 text-4xl text-blue-500">{data}</h1>
    </main>
  );
}
