import { useEffect } from 'react'
import secureSessionStorage from '../../utils/secureSessionStorage'

const ClearSessionOnRefresh = () => {
  useEffect(() => {
    const handleBeforeUnload = () => {
      secureSessionStorage.clearKeys()
      sessionStorage.clear()
    }

    window.addEventListener('beforeunload', handleBeforeUnload)

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload)
    }
  }, [])

  return null
}

export default ClearSessionOnRefresh
