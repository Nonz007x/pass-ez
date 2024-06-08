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
        vault_id: '9c9eedc7-0bc3-4659-b7c9-4852c091a49d',
        name: "cunny",
        note: "note test",
        type: 1
      }
    })
  } catch (error) {
    throw error
  }
}