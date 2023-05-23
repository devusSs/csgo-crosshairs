import axios from "axios";
import { User } from "./types";

const BASE_URL = import.meta.env.VITE_BASE_URL;


export async function getMe() {
    try {
        const res = await axios.get(`${BASE_URL}/api/users/me` , {withCredentials:true})

        return res
    } catch (error) {

        return error
    }
}

export async function getCrosshairsUser() {
    try {
        const res = await axios.get(`${BASE_URL}/api/crosshairs` , {withCredentials:true})

        return res
    } catch (error) {

        return error
    }

}

export async function createUser(user: User) {
    try {
        const res = await axios.post(`${BASE_URL}/api/users/register`, JSON.stringify(user))
        return res
    } catch (error) {

        return error
    }
}; 
    
export async function loginUser(user: any) {
        
    try {
        const res = await axios.post(`${BASE_URL}/api/users/login`, JSON.stringify(user), {withCredentials:true})
        
        return res
    } catch (error) {
        
        return error
    }
}

export async function logoutUser() {
    try {
        const res = await axios.get(`${BASE_URL}/api/users/logout`, {withCredentials:true})
 
        return res
    } catch (error) {
        
        return error
    }
}

export async function verifyMail(url: any) {
    try {
        const res = await axios.get(url)

        return res
    } catch (error) {

        return error
    }
}



