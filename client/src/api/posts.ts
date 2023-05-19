import axios from "axios";

interface crosshairPosts {
    added: string;
    code: string;
    note: string;
}

export function getCrosshairPosts() {
    return axios
        .get<crosshairPosts[]>("/api/crosshair", {params: {_sort: "added:desc"} })
        .then((res) => res.data);
}

