import axios from "axios";



const api = axios.create({
    baseURL: "/",
    withCredentials: true, // Для работы с httpOnly cookies (если сервер использует)
});

api.interceptors.request.use((config) => {
    const token = sessionStorage.getItem("access_token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});


export default api;