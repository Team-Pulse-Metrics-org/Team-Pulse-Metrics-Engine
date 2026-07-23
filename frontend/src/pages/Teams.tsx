import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
interface Developer {
  id:string;
  name: string;
  commits: number;
  velocity: number;
  tasksResolved: number;
  openIssues: number;
}

const API = import.meta.env.VITE_API_URL;

function Teams() {
  const navigate = useNavigate();
   const [developers, setDevelopers] = useState<Developer[]>([]);
   useEffect(() => {
  const token = localStorage.getItem("app_token");

  fetch(`${API}/api/v1/teams`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
    .then(async (res) => {
      console.log("Status:", res.status);

      const data = await res.json();
      console.log("Response:", data);

      if (res.ok && Array.isArray(data)) {
        setDevelopers(data);
      } else {
        console.error("API Error:", data);
        setDevelopers([]);
      }
    })
    .catch((err) => console.error(err));
}, []);
  return (
           <div className="p-6">
  {/* Header */}
  <div className="mb-8">
    <h1 className="text-3xl font-bold text-white">Teams</h1>
    <p className="text-slate-400 mt-2">
      Developer productivity overview
    </p>
  </div>
 <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {developers.map((developer, index) => (
          <div
            key={index}
            className="
              bg-slate-900/70
              backdrop-blur-lg
              border border-slate-700
              rounded-2xl
              p-6
              shadow-2xl
              hover:-translate-y-2
              hover:border-cyan-400
              hover:shadow-cyan-500/20
              transition-all duration-300
            "
          >
            <h2 className="text-xl font-semibold text-white mb-4">
              👤 {developer.name}
            </h2>

            <div className="space-y-3 text-slate-300">
              <div className="flex justify-between">
                <span>Commits</span>
                <span className="font-semibold text-cyan-400">
                  {developer.commits}
                </span>
              </div>

              <div className="flex justify-between">
                <span>Velocity Score</span>
                <span className="font-semibold text-green-400">
                  {developer.velocity}
                </span>
              </div>

              <div className="flex justify-between">
                <span>Tasks Resolved</span>
                <span className="font-semibold text-purple-400">
                  {developer.tasksResolved}
                </span>
              </div>

              <div className="flex justify-between">
                <span>Open Issues</span>
                <span className="font-semibold text-red-400">
                  {developer.openIssues}
                </span>
              </div>
            </div>

            <button
 onClick={() => navigate("/profile")}
  className="
    w-full
    mt-6
    border border-cyan-500
    text-cyan-400
    py-2
    rounded-xl
    hover:bg-cyan-500
    hover:text-white
    transition-all duration-300
  "
>
  View Profile
</button>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Teams;
 