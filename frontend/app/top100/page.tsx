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

    if (loading) return <p className="p-4">Loading your top 100...</p>;

    return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Your Top 100 Movies</h1>
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
            <h2 className="text-lg font-semibold">{item.rank}. {item.movie.title}</h2>
            <p className="text-sm text-gray-600 mb-1">{item.notes || "No notes added"}</p>
          </div>
        ))}
      </div>
    </div>
  );
}