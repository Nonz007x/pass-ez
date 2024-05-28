import React, { useState, useEffect } from 'react';
import { loginHandler, registerHandler, validateEmail } from '../utils/auth'

export default function Register() {
  const [email, setEmail] = useState('');
  const [emailError, setEmailError] = useState('');
  const [password, setPassword] = useState('');


  useEffect(() => {
    // axios({
    //   method: 'get',
    //   url: 'http://localhost:4090/api/v1/',
    // }).then(response => {
    //   console.log(response.data);
    // })
  }, []);

  const handleRegister = async (e) => {
    e.preventDefault()
    if (!validateEmail(email)) {
      setEmailError('Invalid email format.')
      return
    }

    setEmailError('')
    try {
      const response = await registerHandler(email, password)
      console.log(response.data)
    } catch (error) {
      if (error.response) {
        setEmailError(error.response.data.message)
      } else if (error.request) {
        console.error(error.request)
      } else {
        console.error('Error', error.message)
      }
    }

  }

  const handleLogin = async (e) => {
    e.preventDefault()

    if (!validateEmail(email)) {
      setEmailError('Invalid email format.')
      return
    }

    setEmailError('')
    try {
      const response = await loginHandler(email, password)
      // console.log(response.data)
    } catch (error) {
      if (error.response) {
        setEmailError(error.response.data.message)
      } else if (error.request) {
        console.error(error.request)
      } else {
        console.error('Error', error.message)
      }
    }
  }

  return (
    <>
      <h1>Register</h1>
      <form onSubmit={handleRegister}>
        <label>Email</label><br />
        <input
          type="text"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        /><br />
        {emailError ? <><span style={{ color: 'red' }}>{emailError}</span><br /></> : ''}
        <label>Password</label><br />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        /><br />
        <button type="submit">Submit</button>
      </form>
      <button onClick={handleLogin}>Submit</button>
    </>
  )
}