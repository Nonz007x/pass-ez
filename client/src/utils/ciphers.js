import sendRequest from "./request"

export const createItem = async () => {

  try {
    await sendRequest({
      method: 'post',
      endpoint: '/v1/ciphers',
      headers: {
        Authorization: `Bearer ${sessionStorage.getItem('access_token')}`,
        'Content-Type': 'application/json',
      },
      data: {
        vault_id: '044c35e9-6903-4df7-809d-7d128d1d414e',
        name: "gayer",
        note: "note test",
        type: 1
      }
    })
  } catch (error) {
    throw error
  }
}