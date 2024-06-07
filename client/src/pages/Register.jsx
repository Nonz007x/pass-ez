import { useNavigate, Navigate } from 'react-router-dom';
import { registerHandler, validateEmail } from '../utils/auth'
import React, { useState, useEffect } from 'react';

export default function Register() {
  const [email, setEmail] = useState('')
  const [emailError, setEmailError] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()

  const handleRegister = async (e) => {
    e.preventDefault()
    if (!validateEmail(email)) {
      setEmailError('Invalid email format.')
      return
    }

    setEmailError('')
    try {
      const response = await registerHandler(email, password)
      console.log(response)
      navigate('/login');
    } catch (error) {
      if (error.status == 409) {
        setEmailError(error.data.message)
      } else {
        console.log(error.message)
      }
    }

  }

  return (
    <>
      <a href="/login">Login</a>
      {/* <Navigate to="/login" replace /> */}
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
    </>
  )
}