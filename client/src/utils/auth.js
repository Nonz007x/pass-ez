import * as crypto from './crypto'
import { toBase64, base64ToUint8Array, strToArrayBuffer } from './utils'
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

  const response = await axios({
    method: 'post',
    url: 'http://localhost:4090/api/v1/register',
    headers: {},
    data: {
      email: email,
      password: b64HashedPassword,
      salt: b64Salt,
      vault_key: b64ProtectedKey
    }
  })
  return response
}

export const loginHandler = async (email, password) => {
  const masterKey = await crypto.deriveKey(password, email)

  const hashedPassword = await crypto.deriveKey(password, masterKey)
  const b64HashedPassword = await crypto.exportKey(hashedPassword)

  let response = await axios({
    method: 'post',
    url: 'http://localhost:4090/api/v1/login',
    data: {
      email: email,
      password: b64HashedPassword,
    }
  })

  const { salt, key } = response.data
  const stretchedKey = await crypto.stretchKey(masterKey, base64ToUint8Array(salt))


  const index = key.indexOf('|')
  const iv = key.slice(0, index)
  const sliced_key = key.slice(index + 1)

  const vault_key = await crypto.decryptData(
    stretchedKey,
    base64ToUint8Array(sliced_key),
    base64ToUint8Array(iv)
  )

  // TODO 
  // fetch items
  // decrypt items with vault_key

  return response
}

export const validateEmail = (email) => {
  const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return re.test(String(email).toLowerCase());
};
