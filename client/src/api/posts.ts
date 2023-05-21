import axios from "axios";
import { successResponse } from "./types";


export function getCrosshairPosts() {
    const axiosInstance = axios.create({
        baseURL: import.meta.env.VITE_BASE_URL,
        withCredentials: true,
    })

    return axiosInstance
        .get<successResponse[]>("/api/admins/events/user_registered")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

