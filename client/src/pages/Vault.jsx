import React, { useEffect, useState } from 'react'
import sendRequest from '../utils/request'
import { Navigate, useNavigate } from 'react-router-dom'
import { isAuthenticated } from '../utils/auth'
import { createItem } from '../utils/ciphers'
import secureSessionStorage from '../utils/secureSessionStorage'

export default function Vault() {
  const [test, setTest] = useState('')
  const [isAuth, setIsAuth] = useState(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  const handleCreateItem = async (e) => {
    e.preventDefault()

    try {
      await createItem()
    } catch (error) {
      console.error(error)
    }
  }

  const handleLogout = () => {
    sessionStorage.clear()
    secureSessionStorage.clearKeys();
  }

  useEffect(() => {
    const setThing = async () => {
      let response
      try {
        response = await sendRequest({
          method: 'get',
          endpoint: '/v1/test',
          headers: {
            'Authorization': `Bearer ${sessionStorage.getItem('access_token')}`,
          }
        })
      } catch (error) {
        console.error(error.data)
      }
      setTest(response)
    }

    const checkAuth = async () => {
      const authStatus = await isAuthenticated()
      setIsAuth(authStatus)
      setLoading(false)
    }

    checkAuth()
    if (isAuth) {
      setThing()
      console.log(secureSessionStorage.getKeys())
    }

  }, [isAuth])

  if (loading) {
    return <div>Loading...</div>
  }

  if (!isAuth) {
    return <Navigate to="/login" replace />
  }

  return (
    <>
    <a href="" onClick={handleLogout}>Logout</a>
      <h2>{test}</h2>
      <div>Vault</div>
      <form onSubmit={handleCreateItem}>
        <button type='submit'>create</button>
      </form>
    </>
  )
}
