// components/SearchBar.tsx

'use client';

import { useState } from "react"; //this let us remember what we type

type Props = {
    onSearch:(query: string) => void;
};

export default function SearchBar({onSearch}:Props) {  //this create a searchbar component which accept onsearch fucntion
    const[query, setQuery] = useState("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault(); //stop the page from refreshing
        onSearch(query);
    };

    return (
        <form onSubmit={handleSubmit} className="flex gap-2">
        <input
        type="text"
        placeholder="Search a movie"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        className="border rounded px-3 py-2 w-full"
        />
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">
             Search 
            </button>
            </form>

    );
}