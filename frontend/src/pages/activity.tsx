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

const [currentPage, setCurrentPage] = useState(1);
const [sortOrder, setSortOrder] = useState("latest");
const eventsPerPage = 10;
useEffect(() => {
  setCurrentPage(1);
}, [search, developerFilter, typeFilter, repoFilter]);

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
console.log(activity.type,payload);
 console.log("Payload:", payload);
console.log("Developer field:", payload.author);
console.log("Developer field:", payload.developer);
console.log("Developer field:", payload.action_by);
console.log("Developer field:", payload.created_by);
  return {
    id: activity.id,
    timestamp: activity.logged_at,
    displayTime: new Date(
    activity.logged_at
    ).toLocaleString(),

   developer:
  payload.developer ||
  payload.author ||
  payload.action_by ||
  payload.created_by ||
  payload.sender?.login ||
  payload.pull_request?.user?.login ||
  "Unknown",

    type: activity.type || "Unknown",

    repository:
      payload.repository?.name ||
      payload.repository ||
      "Unknown",

    message:
      payload.commits?.[0]?.message ||
      payload.pull_request?.title ||
      "No message",

  
  };
});

const sortedActivities = formattedActivities.sort(
  (a :any, b:any) =>
    new Date(b.timestamp).getTime() -
    new Date(a.timestamp).getTime()
);

setActivities(sortedActivities);
console.log("Formatted:", sortedActivities);
setActivities(sortedActivities);
})
.catch((err) =>
  console.error(
    "Failed to fetch activities:",
    err
  )
);
}, []);
  const filteredActivities = activities.filter((activity) => {
  const developer = (activity.developer || "").toLowerCase();
  const message = (activity.message || "").toLowerCase();
  const searchText = search.toLowerCase();

  const matchesSearch =
    developer.includes(searchText) ||
    message.includes(searchText);

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
const sortedFilteredActivities = [...filteredActivities].sort(
  (a: any, b: any) => {
    if (sortOrder === "latest") {
      return (
        new Date(b.timestamp).getTime() -
      new Date(a.timestamp).getTime()
      );
    }

    return (
      new Date(a.timestamp).getTime() -
      new Date(b.timestamp).getTime()
    );
  }
);
  const totalEvents = filteredActivities.length;
  const indexOfLastEvent = currentPage * eventsPerPage;
const indexOfFirstEvent = indexOfLastEvent - eventsPerPage;

const currentActivities = sortedFilteredActivities.slice(
  indexOfFirstEvent,
  indexOfLastEvent
);

const totalPages = Math.ceil(
  filteredActivities.length / eventsPerPage
);

  return (
    <div className="p-8 text-white">
      <h1 className="text-4xl font-bold">Activity</h1>

      <p className="text-slate-400 mt-1">
        All engineering events across repositories
      </p>
     <p className="text-slate-400">
    Total Events: {totalEvents}
    </p>
      

<p className="text-slate-400">
  Showing {(currentPage - 1) * eventsPerPage + 1} -
  {Math.min(
    currentPage * eventsPerPage,
    filteredActivities.length
  )}{" "}
  of {filteredActivities.length} events
  (Page {currentPage} of {totalPages})
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
      <div className="grid grid-cols-4 gap-3 mt-4">
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
             <option value="All Developers" className="text-black">
    All Developers
  </option>

            {[...new Set(
              activities.map((a) => a.developer)
            )].map((developer) => (
              <option
                key={developer}
                value={developer}
                className="text-black"
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
  onChange={(e) => setTypeFilter(e.target.value)}
  className="bg-transparent outline-none"
>
  <option value="All Types" className="text-black">
    All Types
  </option>

  {[...new Set(activities.map((a) => a.type))].map((type) => (
    <option
      key={type}
      value={type}
      className="text-black"
    >
      {type}
    </option>
  ))}
</select>
         
        </div>

        {/* Repository */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-300">
            Repository:
          </span>

         <select
  value={repoFilter}
  onChange={(e) => setRepoFilter(e.target.value)}
  className="bg-transparent outline-none"
>
  <option value="All Repositories" className="text-black">
    All Repositories
  </option>

  {[...new Set(activities.map((a) => a.repository))].map((repo) => (
    <option
      key={repo}
      value={repo}
      className="text-black"
    >
      {repo}
    </option>
  ))}
</select>
        </div>
        {/* Sort */}
<div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
  <span className="text-slate-400">
    Sort:
  </span>

  <select
    value={sortOrder}
    onChange={(e) => setSortOrder(e.target.value)}
    className="bg-transparent outline-none"
  >
    <option value="latest" className="text-black">
      Latest First
    </option>

    <option value="oldest" className="text-black">
      Oldest First
    </option>
  </select>
</div>
      </div>
<div className="flex justify-between items-center mt-5 mb-2">
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
                
                <th className="p-4"></th>
              </tr>
            </thead>

            <tbody>
              {currentActivities.map((activity) => (
                <React.Fragment key={activity.id}>
                  <tr className="border-t border-slate-700 hover:bg-slate-800">
                    <td className="p-4">
                      {activity.displayTime}
                    </td>

                    <td className="p-4">
                      {activity.developer}
                    </td>

                    <td className="p-4">
                  <span
              className={`px-2 py-1 rounded-md text-sm text-white ${
              activity.type === "git_commit"
              ? "bg-blue-600"
        : activity.type === "pull_request_closed"
        ? "bg-green-600"
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
      <div className="flex justify-center gap-4 mt-6">
  <button
    disabled={currentPage === 1}
    onClick={() => setCurrentPage(currentPage - 1)}
    className="px-4 py-2 bg-slate-700 rounded disabled:opacity-50"
  >
    Previous
  </button>

  <span>
    Page {currentPage} of {totalPages}
  </span>

  <button
    disabled={currentPage === totalPages}
    onClick={() => setCurrentPage(currentPage + 1)}
    className="px-4 py-2 bg-slate-700 rounded disabled:opacity-50"
  >
    Next
  </button>
</div>
    </div>
  );
}