export function strToArrayBuffer(str) {
  var buf = new ArrayBuffer(str.length * 2); // 2 bytes for each char
  var bufView = new Uint16Array(buf);
  for (var i = 0, strLen = str.length; i < strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return buf;
}

export const toBase64 = (bytes) => {
  if (bytes instanceof Uint8Array) {
    return btoa(String.fromCharCode(...bytes))
  } else {
    return btoa(String.fromCharCode(...new Uint8Array(bytes)))
  }
}

export function base64ToUint8Array(base64) {
  // Decode Base64 to binary string
  const binaryString = atob(base64);
  // Create a Uint8Array with the same length as the binary string
  const bytes = new Uint8Array(binaryString.length);
  // Convert each character in the binary string to a byte
  for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes;
}