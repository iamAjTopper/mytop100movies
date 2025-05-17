//services/api.ts

import { error } from "console";

export async function searchMovie(query:string) {
    const res = await fetch(`http://localhost:8080/search?query=${query}`);
    if(!res.ok) throw new Error('Failed to fetch');
    return res.json()

}