import axios from "axios";
import { User, successResponse } from "./types";

const base_URL = import.meta.env.VITE_BASE_URL;

const axiosInstance = axios.create({
    baseURL: base_URL,
    withCredentials: true,
})

export function getCrosshairPosts() {

    return axiosInstance
        .get<successResponse[]>("/api/crosshairs")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

export async function createAnUser(user) {

    return await axios
        .post<User[]>("/api/users/register", user)
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

export async function loginUser(user) {
      
    return await axios
        .post<User[]>("/api/users/login", user)
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

export async function logoutUser() {

    return await axiosInstance
        .get<successResponse[]>("/api/users/logout")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}
