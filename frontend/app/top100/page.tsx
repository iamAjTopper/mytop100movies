"use client";

import { useEffect, useState } from "react";
import Image from "next/image";

type MovieItem = {
    user_movie_id:  number;
    rank:           number;
    notes:          string;
    movie: {
        id:         number;
        tmdb_id:    number;
        title:      string;
        overview:   string;
        poster_url: string

    };
};

export default function Top100Page() {
    const [topMovies, setTopMovies] = useState<MovieItem[]>([]);
    const [loading, setloading] = useState(true);
    const [clearing, setClearing] = useState(false);


    useEffect(()=> {
        const fetchTop100 = async () => {
            try {
                const res = await fetch("http://localhost:8080/top100/get?user_id=1");
                const data = await res.json();
                setTopMovies(data);
            } catch(err) {
                console.error("Failed to fetch top 100 list",err);
            } finally {
                setloading(false)
            }
        };
        fetchTop100();
    },[]);

    const handleClearTop100 = async () =>{
      if(!window.confirm("Are you sure you want to clear yout top 100 list?")) return;

      setClearing(true);
      try {
        const res = await fetch("http://localhost:8080/top100/delete?user_id=1",{
          method: "DELETE",
        });

        if (res.ok){
          alert("Top 100 list cleared!");
          setTopMovies([]);
        } else {
          const text = await res.text();
          alert("failed to clear:"+text);
        }

      }catch(err) {
        alert("Error clearing list: " + err );
      }finally {
        setClearing(false)
      }
    };

    if (loading) return <p className="p-4">Loading your top 100...</p>;

    return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Your Top 100 Movies</h1>
        <button
          onClick={handleClearTop100}
          disabled={clearing}
          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          {clearing ? "Clearing..." : "Clear Top 100"}
        </button>
      </div>

      {topMovies.length === 0 ? (
        <p className="text-gray-600">You haven't added any movies yet.</p>
      ) : (
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {topMovies.map((item) => (
            <div key={item.user_movie_id} className="border p-4 rounded shadow">
              <Image
                src={item.movie.poster_url || "https://via.placeholder.com/150"}
                alt={item.movie.title}
                width={200}
                height={300}
                className="w-full h-auto mb-2"
              />
              <h2 className="text-lg font-semibold">
                {item.rank}. {item.movie.title}
              </h2>
              <p className="text-sm text-gray-600 mb-1">
                {item.notes || "No notes added"}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}