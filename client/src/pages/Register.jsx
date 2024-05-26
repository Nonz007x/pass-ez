import React, { useState, useEffect } from 'react';
import { registerHandler } from '../utils/crypto'
import axios from 'axios'

export default function Register() {
  const [a, setA] = useState("a");


  useEffect(() => {
    const generateKey = async () => {
      const key = await registerHandler("nonz007x@gmail.com", "securepass")
      setA(key);
    };

    // axios({
    //   method: 'get',
    //   url: 'http://localhost:4090/api/v1/',
    // }).then(response => {
    //   console.log(response.data);
    // })

    generateKey();
  }, []);

  return (
    <>
      <h1>
        {a}
      </h1>
    </>
  );
}
