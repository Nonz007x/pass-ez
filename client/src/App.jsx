import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const [text, setText] = useState("")

  const getText = () => {
    setText("Hello")
    // fetch('localhost:4090/api/v1/home')
    // .then(response => response.json())
    // .then(data => setText(data))
  }

  return (
    <>
      <h1>Vite + React</h1>
      <h1>{text}</h1>
      <button onClick={() => getText()}>Press me!</button>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.jsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
