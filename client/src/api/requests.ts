import axios from "axios";
import { User, userSuccessResponse } from "./types";

const base_URL = import.meta.env.VITE_BASE_URL;


export async function getMe() {
    try {
        const res = await axios.get('https://api.dropawp.com/api/users/me' , {withCredentials:true})

        return res
    } catch (error) {

        return error
    }
}

export async function getCrosshairsUser() {
    try {
        const res = await axios.get('https://api.dropawp.com/api/crosshairs' , {withCredentials:true})

        return res
    } catch (error) {

        return error
    }

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
      
    try {
        const res = await axios.post('https://api.dropawp.com/api/users/login', JSON.stringify(user), {withCredentials:true})
 
        return res
    } catch (error) {

        return error
    }
}

export async function verifyMail(url) {
    try {
        const res = await axios.get(url)

        return res
    } catch (error) {

        return error
    }
}

export async function logoutUser() {

    return await axios
        .get("/api/users/logout")
        .then((res) => res.status === 200 ? res.data : Promise.reject(res));
}


