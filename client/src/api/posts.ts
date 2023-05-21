import axios from "axios";
import { successResponse } from "./types";


export function getCrosshairPosts() {
    const axiosInstance = axios.create({
        baseURL: import.meta.env.VITE_BASE_URL,
        withCredentials: true,
        headers: {"Origin": "http://localhost:5173"}
    })

    return axiosInstance
        .get<successResponse[]>("/api/crosshairs")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

