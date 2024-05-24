import { useState } from 'react'
import './App.css'
import axios from 'axios'

function App() {
  const [count, setCount] = useState(0)
  const [text, setText] = useState("")

  const getText = () => {
    axios({
      method: 'get',
      url: 'http://localhost:4090/api/v1/',
    })
      .then(response => {
        setText(response.data.hello);
      })
    // .then(data => setText(data))
  }

  return (
    <>
      <h1>{text}</h1>
      <button onClick={() => getText()}>Press me!</button>
    </>
  )
}

export default App
