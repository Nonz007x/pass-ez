export const registerHandler = async (email, password) => {
  try {

    const masterKey = await deriveKey(password, email)

    const hashedPassword = await deriveKey(password, masterKey)

    const { salt, stretchedKey } = await generateSaltAndStretch(masterKey)

    const vaultKey = generateRandomBytes(32)

    const iv = generateRandomBytes(16)

    const cipherText = await encryptData(stretchedKey, vaultKey, iv) 

    const b64CipherText = toBase64(iv) + '|' + toBase64(cipherText)

    return b64CipherText
  } catch (error) {
    console.error('Error registering:', error)
    throw error
  }
}

const encryptData = async (key, data, iv) => {
  return window.crypto.subtle.encrypt(
    {
      name: 'AES-CBC',
      iv: iv,
    },
    key,
    data
  )
}

const generateRandomBytes = (length) => {
  const bytes = new Uint8Array(length);
  window.crypto.getRandomValues(bytes);
  return bytes;
}

const toBase64 = (bytes) => {
  if (bytes instanceof Uint8Array) {
    return btoa(String.fromCharCode(...bytes))
  } else {
    return btoa(String.fromCharCode(...new Uint8Array(bytes)))
  }
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

export const generateSaltAndStretch = async (secret) => {
  const encoder = new TextEncoder()
  const salt = generateRandomBytes(16)
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

  return { salt, stretchedKey }
}

const exportKey = async (key) => {
  const exportedKey = await window.crypto.subtle.exportKey('raw', key)
  const buffer = new Uint8Array(exportedKey)
  const base64Key = btoa(String.fromCharCode(...buffer))
  return base64Key
}
