// app/search/page.tsx

"use client";

import { useState } from "react"; //heps to remember data
import SearchBar from "@/components/SearchBar"; // the search inpyut component
import { searchMovie } from "@/services/api"; // a function to talk to movie
import Image from "next/image";
//import { headers } from "next/headers";

//defining what movie looks like
type Movie = {
  id: number;
  title: string;
  poster_url: string | null;
  release_date: string;
};

export default function SearchPage() {
  const [results, setResults] = useState<Movie[]>([]); // stors the list of the movie found
  const [loading, setLoading] = useState(false); //remeber if we are waiting for the rsults

  //the search function
  const handleSearch = async (query: string) => {
    try {
      setLoading(true);
      const data = await searchMovie(query);
      console.log("API Data:", data);
      setResults(data);
    } catch (err) {
      console.error("Search error:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToTop100 = async (movie:Movie) => {
    try {
      const res = await fetch("http://localhost:8080/top100/add", {
        method: "POST",
        headers: {
          "Content-Type":"application/json",
        },
        body: JSON.stringify({
          user_id:1,
          tmdb_id: movie.id,
          title: movie.title,
          overview:"",
          poster_url: movie.poster_url,
          release_date: movie.release_date,
         // rank: 0, //backend can assign based on logic
          notes:""
        }),
      });

      if(!res.ok) {
        throw new Error("Failed to add movie to top 100");
      }

      alert(`${movie.title} added to your top 100!`);

    }catch(err) {
      console.error("Add error",err)
      alert("Something went wrong while adding the movie.");
    }
  };

  //displaying the page
  
  return (
    <div className="p-6 max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Search Movies</h1>
      <SearchBar onSearch={handleSearch} />

      {loading && <p className="mt-4">Loading...</p>}

      <div className="grid grid-cols-2 gap-4 mt-6">
        {results.map((movie) => (
          <div key={movie.id} className="border p-4 rounded shadow">
            <Image
              src={movie.poster_url || 'https://via.placeholder.com/150' }
              alt={movie.title}
              height={230}
              width={200}
              className="w-full h-auto mb-2"
            />

            <h2 className="text-lg font-semibold">{movie.title}</h2>
            <p className="text-sm text-gray-600">{movie.release_date}</p>

            <button
              onClick={() => handleAddToTop100(movie)}
              className="mt-2 px-4 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
            >
              Add to Top 100
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}