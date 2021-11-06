import React, { useState, useEffect } from "react";
import Slider from "@mui/material/Slider";

export default function Distance() {
  let [distance, setDistance] = useState(200);

  useEffect(() => {
    const timer = setInterval(() => {
      fetch("/api/distance")
        .then((resp) => {
          if (resp.ok) {
            return resp.json();
          }
          throw new Error(`request failed with ${resp.statusText}`);
        })
        .then((data) => setDistance(data))
        .catch((err) => console.log(err));
    }, 1000);

    return () => clearInterval(timer);
  });

  return <Slider disabled value={distance} step={10} marks min={0} max={200} />;
}
