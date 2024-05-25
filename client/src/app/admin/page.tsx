"use client";

import axios from "axios";
import { useState, useEffect } from "react";
// import FetchData from "@/utils/helloWorld";

const FetchData = async () => {
  try {
    const response = await axios.get(
      "http://localhost:8080/api/v1/illustrations/list?p=0"
    );
    console.log(response.data);
    return response.data;
  } catch (error) {
    console.error(error);
    return null;
  }
};

export default function Home() {
  const [data, setData] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const result = await FetchData();
      setData(result);
    };
    fetchData();
  }, []);

  return (
    <main>
      <h1 className="m-10 text-4xl text-red-700">Hello World from admin!!</h1>
      {data.length > 0 ? (
        <>
          {data.map((item) => (
            <>
              <div key={item.Image.id}>image : {item.Image.title}</div>
              {item.Character.map((c) => (
                <>
                  <div key={c.id}>character : {c.name}</div>
                  <br></br>
                </>
              ))}
              {item.Category.map((c) => (
                <>
                  <div key={c.ParentCategory.id}>
                    category.parent : {c.ParentCategory.name}
                  </div>
                  <br></br>
                  <div key={c.ChildCategory[0].id}>
                    category.child : {c.ChildCategory[0].name}
                  </div>
                </>
              ))}
            </>
          ))}
        </>
      ) : (
        <p>loading</p>
      )}
    </main>
  );
}
