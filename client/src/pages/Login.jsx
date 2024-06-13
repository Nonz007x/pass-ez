import { isAuthenticated, loginHandler, validateEmail } from '../utils/auth';
import React, { useState, useEffect } from 'react';
import { useNavigate, Navigate } from 'react-router-dom';

export default function Login() {
  const [email, setEmail] = useState('');
  const [emailError, setEmailError] = useState('');
  const [password, setPassword] = useState('');
  const [isAuth, setIsAuth] = useState(null)
  const [loading, setLoading] = useState(true);

  const navigate = useNavigate()

  const handleLogin = async (e) => {
    e.preventDefault()

    if (!validateEmail(email)) {
      setEmailError('Invalid email format.')
      return
    }

    setEmailError('')
    try {
      await loginHandler(email, password)
      navigate('/vault')
    } catch (error) {
      if (error.status == 404) {
        setEmailError(error.data.message)
      } else {
        console.error(error)
      }
    }
  }

  useEffect(() => {
    const checkAuth = async () => {
      const authStatus = await isAuthenticated()
      setIsAuth(authStatus)
      setLoading(false)
    }
    checkAuth()
  }, [])

  if (loading) {
    return <div>Loading...</div>
  }

  if (isAuth) {
    return <Navigate to="/vault" replace />
  }

  return (
    <>
      <a href="/register">Register</a>
      <h1>Login</h1>
      <form onSubmit={handleLogin}>
        <label>Email</label><br />
        <input
          type="text"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        /><br />
        {emailError && <><span style={{ color: 'red' }}>{emailError}</span><br /></>}
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
