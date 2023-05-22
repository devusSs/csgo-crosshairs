import axios from "axios";
import { User, userSuccessResponse } from "./types";

const base_URL = import.meta.env.VITE_BASE_URL;

const axiosInstance = axios.create({
    baseURL: base_URL,
    withCredentials: true,
})

export function getCrosshairPosts() {

    return axiosInstance
        .get("/api/crosshairs")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

export async function createUser(user: User) {
        try {
            const res = await axios.post('https://api.dropawp.com/api/users/register', JSON.stringify(user))

            return res
        } catch (error) {
   
            return error
        }
  }; 

export async function loginUser(user) {
      
    return await axiosInstance
        .post<User[]>("/api/users/login", user)
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}

export async function logoutUser() {

    return await axiosInstance
        .get("/api/users/logout")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}


