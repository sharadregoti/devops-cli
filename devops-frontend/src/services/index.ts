import axios from "axios";

const api = axios.create({
  baseURL: 'http://localhost:9753/v1',
});

export default api;