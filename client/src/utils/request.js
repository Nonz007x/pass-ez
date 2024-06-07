import axios from 'axios';

export const sendRequest = async (options) => {
  try {
    const response = await axios({
      ...options,
      url: `${import.meta.env.VITE_API_URL}${options.endpoint}`,
      timeout: 10000,
    })
    return response.data
  } catch (error) {
    if (error.response) {
      throw error.response;
    } else if (error.request) {
      throw new Error('Network error: No response received from server');
    } else {
      throw new Error(error);
    }
  }
}

export default sendRequest