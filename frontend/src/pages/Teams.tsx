//import { useEffect, useState } from "react";
/*interface Developer {
  name: string;
  commits: number;
  velocity: number;
  tasksResolved: number;
  openIssues: number;
}*/
interface Developer {
  name: string;
  commits: number;
  velocity: number;
  tasksResolved: number;
  openIssues: number;
}

const developers: Developer[] = [
  {
    name: "Harshitha V",
    commits: 45,
    velocity: 8.5,
    tasksResolved: 12,
    openIssues: 3,
  },
  {
    name: "John Doe",
    commits: 38,
    velocity: 7.9,
    tasksResolved: 10,
    openIssues: 2,
  },
  {
    name: "Jane Smith",
    commits: 52,
    velocity: 9.1,
    tasksResolved: 15,
    openIssues: 1,
  },
];
function Teams() {
  /*  const [developers, setDevelopers] = useState<Developer[]>([]);
    useEffect(() => {
  fetch("http://localhost:8080/api/v1/teams")
    .then((response) => response.json())
    .then((data) => setDevelopers(data))
    .catch((error) => console.error("Error fetching teams:", error));
}, []);*/
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
 