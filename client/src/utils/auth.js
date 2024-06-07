import { toBase64, base64ToUint8Array } from './utils'
import { exportKey } from './crypto'
import sendRequest from './request'
import * as crypto from './crypto'
import secureSessionStorage from './secureSessionStorage'
import axios from 'axios'

export const registerHandler = async (email, password) => {
  const masterKey = await crypto.deriveKey(password, email)

  const hashedPassword = await crypto.deriveKey(password, masterKey)
  const b64HashedPassword = await crypto.exportKey(hashedPassword)

  const salt = crypto.generateRandomBytes(16)
  const b64Salt = toBase64(salt)

  const stretchedKey = await crypto.stretchKey(masterKey, salt)
  const vaultKey = crypto.generateRandomBytes(32)
  const iv = crypto.generateRandomBytes(16)

  const protectedKey = await crypto.encryptData(stretchedKey, vaultKey, iv)
  const b64ProtectedKey = toBase64(iv) + '|' + toBase64(protectedKey)

  let response

  try {
    response = await sendRequest({
      method: 'post',
      endpoint: '/v1/register',
      headers: {},
      data: {
        email: email,
        password: b64HashedPassword,
        salt: b64Salt,
        vault_key: b64ProtectedKey
      }
    })
  } catch (error) {
    throw error
  }

  return response
}

export const loginHandler = async (email, password) => {
  const masterKey = await crypto.deriveKey(password, email)

  const hashedPassword = await crypto.deriveKey(password, masterKey)
  const b64HashedPassword = await crypto.exportKey(hashedPassword)

  let response

  try {
    response = await sendRequest({
      method: 'post',
      endpoint: '/v1/login',
      data: {
        email: email,
        password: b64HashedPassword,
      }
    })
  } catch (error) {
    throw error
  }

  const { salt, key, token } = response
  const stretchedKey = await crypto.stretchKey(masterKey, base64ToUint8Array(salt))
  
  const vault_key = await crypto.decryptKey(key, stretchedKey)
  
  try {
    secureSessionStorage.setKeys(await exportKey(stretchedKey), toBase64(vault_key))
  } catch (error) {
    throw error
  }
  sessionStorage.setItem("refresh_token", token.refresh_token)
  sessionStorage.setItem("access_token", token.access_token)

  
  // TODO:
  // fetch items
  // decrypt items with vault_key

  return response
}

export const validateEmail = (email) => {
  const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return re.test(String(email).toLowerCase());
};

export const isAuthenticated = async () => {
  const token = sessionStorage.getItem('access_token')
  if (token == null) {
    return false
  }
  try {
    await sendRequest(
      {
        endpoint: '/v1/validate-token',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }
    )
    return true
  } catch (error) {
    return false
  }
}