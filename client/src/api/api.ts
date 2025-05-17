import axios from "axios";



const api = axios.create({
    baseURL: "/",
    withCredentials: true, // Для работы с httpOnly cookies (если сервер использует)
});


export default api;