import React, { useEffect,useState } from "react";


export default function Activity() {
  
const [activities, setActivities] = useState<any[]>([]);
  const [search, setSearch] = useState("");
  const [expandedId, setExpandedId] = useState<number | null>(null);

  const [developerFilter, setDeveloperFilter] =
    useState("All Developers");
  const [typeFilter, setTypeFilter] =
    useState("All Types");
  const [repoFilter, setRepoFilter] =
    useState("All Repositories");
    useEffect(() => {
  fetch("http://localhost:8080/api/v1/activities")
    .then((res) => res.json())
    .then((data) => {
      console.log("Backend data:",data)
      const formattedActivities = data.map((activity: any) => {
        const payload =
          typeof activity.payload === "string"
            ? JSON.parse(activity.payload)
            : activity.payload;

        return {
          id: activity.id,
          timestamp: new Date(
            activity.logged_at
          ).toLocaleString(),
          developer: payload.author,
          type: activity.type,
          repository: payload.repository,
          message:
            payload.commits?.[0]?.message ||
            "No message",
          weight: activity.weight,
        };
      });
      console.log("Formatted:",formattedActivities)
      setActivities(formattedActivities);
    })
    .catch((err) =>
      console.error(
        "Failed to fetch activities:",
        err
      )
    );
}, []);

  const filteredActivities = activities.filter((activity) => {
    const matchesSearch =
      activity.developer
        .toLowerCase()
        .includes(search.toLowerCase()) ||
      activity.message
        .toLowerCase()
        .includes(search.toLowerCase());

    const matchesDeveloper =
      developerFilter === "All Developers" ||
      activity.developer === developerFilter;

    const matchesType =
      typeFilter === "All Types" ||
      activity.type === typeFilter;

    const matchesRepo =
      repoFilter === "All Repositories" ||
      activity.repository === repoFilter;

    return (
      matchesSearch &&
      matchesDeveloper &&
      matchesType &&
      matchesRepo
    );
  });

  const totalEvents = filteredActivities.length;

  return (
    <div className="p-8 text-white">
      <h1 className="text-4xl font-bold">Activity</h1>

      <p className="text-slate-400 mt-1">
        All engineering events across repositories
      </p>
     <p className="text-slate-400">
    Total Events: {totalEvents}
    </p>
      

      {/* Search */}
      <input
        type="text"
        placeholder="Search by developer or commit message..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="w-full mt-8 p-3 rounded-lg border border-slate-600 bg-slate-900"
      />

      {/* Filters */}
      <div className="flex items-center gap-6 mt-6">
        {/* Developer */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-400">
            Developer:
          </span>

          <select
            value={developerFilter}
            onChange={(e) =>
              setDeveloperFilter(e.target.value)
            }
            className="bg-transparent outline-none"
          >
            <option>All Developers</option>

            {[...new Set(
              activities.map((a) => a.developer)
            )].map((developer) => (
              <option
                key={developer}
                value={developer}
              >
                {developer}
              </option>
            ))}
          </select>
        </div>

        {/* Type */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-400">
            Activity Type:
          </span>

          <select
            value={typeFilter}
            onChange={(e) =>
              setTypeFilter(e.target.value)
            }
            className="bg-transparent outline-none"
          >
            <option>All Types</option>

            {[...new Set(
              activities.map((a) => a.type)
            )].map((type) => (
              <option key={type} value={type}>
                {type}
              </option>
            ))}
          </select>
        </div>

        {/* Repository */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-400">
            Repository:
          </span>

          <select
            value={repoFilter}
            onChange={(e) =>
              setRepoFilter(e.target.value)
            }
            className="bg-transparent outline-none"
          >
            <option>All Repositories</option>

            {[...new Set(
              activities.map((a) => a.repository)
            )].map((repo) => (
              <option key={repo} value={repo}>
                {repo}
              </option>
            ))}
          </select>
        </div>
      </div>
<div className="flex justify-between items-center mt-6 mb-3">
  <h2 className="text-xl font-semibold">
    Activity Events
  </h2>

  <p className="text-slate-400">
    Total Events:{" "}
    <span className="text-white font-bold">
      {filteredActivities.length}
    </span>
  </p>
</div>
      {/* Table */}
      <div className="mt-6 rounded-xl border border-slate-700 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead className="bg-slate-800 text-slate-300">
              <tr>
                <th className="p-4">Timestamp</th>
                <th className="p-4">Developer</th>
                <th className="p-4">Type</th>
                <th className="p-4">Repository</th>
                <th className="p-4">Message</th>
                <th className="p-4">Weight</th>
                <th className="p-4"></th>
              </tr>
            </thead>

            <tbody>
              {filteredActivities.map((activity) => (
                <React.Fragment key={activity.id}>
                  <tr className="border-t border-slate-700 hover:bg-slate-800">
                    <td className="p-4">
                      {activity.timestamp}
                    </td>

                    <td className="p-4">
                      {activity.developer}
                    </td>

                    <td className="p-4">
                      <span
                        className={`px-2 py-1 rounded-md text-sm ${
                          activity.type === "Commit"
                            ? "bg-green-600"
                            : activity.type === "PR"
                            ? "bg-blue-600"
                            : "bg-red-600"
                        }`}
                      >
                        {activity.type}
                      </span>
                    </td>

                    <td className="p-4">
                      {activity.repository}
                    </td>

                    <td className="p-4">
                      {activity.message.length > 25
                        ? activity.message.substring(
                            0,
                            25
                          ) + "..."
                        : activity.message}
                    </td>

                    <td className="p-4">
                      {activity.weight}
                    </td>

                    <td className="p-4">
                      <button
                        onClick={() =>
                          setExpandedId(
                            expandedId === activity.id
                              ? null
                              : activity.id
                          )
                        }
                      >
                        {expandedId === activity.id
                          ? "▲"
                          : "▼"}
                      </button>
                    </td>
                  </tr>

                  {expandedId === activity.id && (
                    <tr className="bg-slate-900">
                      <td
                        colSpan={7}
                        className="p-4"
                      >
                        <div className="flex justify-between items-center gap-8">
                          <div>
                            <span className="font-bold">
                              Full Message:
                            </span>{" "}
                            {activity.message}
                          </div>

                          <div>
                            <span className="font-bold">
                              Repository:
                            </span>{" "}
                            {activity.repository}
                          </div>

                          <div>
                            <span className="font-bold">
                              Impact Weight:
                            </span>{" "}
                            {activity.weight}/10
                          </div>
                        </div>
                      </td>
                    </tr>
                  )}
                </React.Fragment>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}