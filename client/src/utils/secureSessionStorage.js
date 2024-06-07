class SecureSessionStorage {
  constructor() {
    if (SecureSessionStorage.instance) {
      return SecureSessionStorage.instance
    }
    this.masterKey = null
    this.userKey = null
    SecureSessionStorage.instance = this
  }

  setKeys(decryptedMasterKey, decryptedUserKey) {
    this.masterKey = decryptedMasterKey
    this.userKey = decryptedUserKey
  }

  getKeys() {
    if (!this.masterKey || !this.userKey) {
      return {}
    }
    return { masterKey: this.masterKey, userKey: this.userKey }
  }

  clearKeys() {
    this.masterKey = null
    this.userKey = null
  }
}

const secureSessionStorage = new SecureSessionStorage()
// Object.freeze(secureSessionStorage)
window.addEventListener('beforeunload', () => {
  secureSessionStorage.clearKeys();
  sessionStorage.clear();
});

export default secureSessionStorage
