import { base64ToUint8Array } from "./utils"

export const encryptData = async (key, data, iv) => {
  return window.crypto.subtle.encrypt(
    {
      name: 'AES-CBC',
      iv: iv,
    },
    key,
    data
  )
}

export const decryptData = async (key, data, iv) => {
  return window.crypto.subtle.decrypt(
    {
      name: 'AES-CBC',
      iv: iv,
    },
    key,
    data
  )
}

export const generateRandomBytes = (length) => {
  const bytes = new Uint8Array(length);
  window.crypto.getRandomValues(bytes);
  return bytes;
}


export const deriveKey = async (secret, inputSalt) => {
  const encoder = new TextEncoder()
  const salt = encoder.encode(inputSalt)
  const keyMaterial = await window.crypto.subtle.importKey(
    'raw',
    encoder.encode(secret),
    { name: 'PBKDF2' },
    false,
    ['deriveBits', 'deriveKey']
  )

  const key = await window.crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt: salt,
      iterations: 600000,
      hash: 'SHA-256',
    },
    keyMaterial,
    { name: 'AES-CBC', length: 256 },
    true,
    ['encrypt', 'decrypt']
  )

  return key
}

export const stretchKey = async (secret, inputSalt) => {
  const encoder = new TextEncoder()
  const salt = encoder.encode(inputSalt)
  const keyMaterial = await window.crypto.subtle.importKey(
    'raw',
    encoder.encode(secret),
    { name: 'HKDF' },
    false,
    ['deriveBits', 'deriveKey']
  )

  const stretchedKey = await window.crypto.subtle.deriveKey(
    {
      name: 'HKDF',
      hash: 'SHA-256',
      salt: salt,
      info: new Uint8Array(),
    },
    keyMaterial,
    { name: 'AES-CBC', length: 256 },
    true,
    ['encrypt', 'decrypt']
  )

  return stretchedKey
}

export const exportKey = async (key) => {
  const exportedKey = await window.crypto.subtle.exportKey('raw', key)
  const buffer = new Uint8Array(exportedKey)
  const base64Key = btoa(String.fromCharCode(...buffer))
  return base64Key
}

export const decryptKey = async (key, stretchedKey) => {
  const index = key.indexOf('|')
  const iv = key.slice(0, index)
  const sliced_key = key.slice(index + 1)

  const vault_key = await decryptData(
    stretchedKey,
    base64ToUint8Array(sliced_key),
    base64ToUint8Array(iv)
  )

  return vault_key
}


export const secureSessionStorage = (() => {
  let masterKey = null
  let userKey = null

  return {
    setKeys: (decryptedMasterKey, decryptedUserKey) => {
      masterKey = decryptedMasterKey
      userKey = decryptedUserKey
    },
    getKeys: () => {
      if (!masterKey || !userKey) {
        throw new Error("Keys are not available in memory")
      }
      return { masterKey, userKey }
    },
    clearKeys: () => {
      masterKey = null
      userKey = null
    }
  }
})()
