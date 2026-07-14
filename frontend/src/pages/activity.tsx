import React, { useEffect, useState } from "react";

//activity page component
export default function Activity() {
  const [activities, setActivities] = useState<any[]>([]);
  const [search, setSearch] = useState("");
  const [expandedId, setExpandedId] = useState<number | null>(null);

  //filter states
  const [developerFilter, setDeveloperFilter] = useState("All Developers");
  const [typeFilter, setTypeFilter] = useState("All Types");
  const [repoFilter, setRepoFilter] = useState("All Repositories");

  const [currentPage, setCurrentPage] = useState(1);
  const [sortOrder, setSortOrder] = useState("latest");
  const eventsPerPage = 10;

  // Reset to first page whenever filters or search change
  useEffect(() => {
    setCurrentPage(1);
  }, [search, developerFilter, typeFilter, repoFilter]);

  // Fetch activities from backend when component loads
  useEffect(() => {
    const token = localStorage.getItem("app_token");

    if (!token) {
      window.location.href = "/login";
      return;
    }

    fetch("http://localhost:8080/api/v1/activities", {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    })
      .then((res) => {
        if (res.status === 401) {
          localStorage.removeItem("app_token");
          window.location.href = "/login";
          throw new Error("Session expired. Please log in again.");
        }

        if (!res.ok) {
          throw new Error("failed to fetch activity records");
        }
        return res.json();
      })
      .then((data) => {
        console.log("Backend data:", data);
        if (!Array.isArray(data)) {
          console.error("Expected array but got:", data);
          return;
        }

        const formattedActivities = data.map((activity: any) => {
          const payload =
            typeof activity.payload === "string"
              ? JSON.parse(activity.payload)
              : activity.payload;

          return {
            id: activity.id,
            timestamp: activity.logged_at,
            displayTime: new Date(activity.logged_at).toLocaleString(),

            developer:
              activity.developer_name ||
              payload.developer ||
              payload.author ||
              payload.action_by ||
              payload.created_by ||
              payload.sender?.login ||
              payload.pull_request?.user?.login ||
              "Unknown",

            type: activity.type || "Unknown",

            repository: (
              payload.repository?.name ||
              payload.repository?.full_name ||
              payload.repository ||
              "Unknown"
            )
              .split("/")
              .pop(),

            message:
              payload.message ||
              payload.commits?.[0]?.message ||
              payload.pull_request?.title ||
              payload.title ||
              "No message",
          };
        });

        setActivities(formattedActivities);
      })
      .catch((err) => console.error("Failed to fetch activities:", err));
  }, []);
  //apply search and filter conditions
  const filteredActivities = activities.filter((activity) => {
    const developer = (activity.developer || "").toLowerCase();
    const message = (activity.message || "").toLowerCase();
    const searchText = search.toLowerCase();

    const matchesSearch =
      developer.includes(searchText) || message.includes(searchText);

    const matchesDeveloper =
      developerFilter === "All Developers" ||
      activity.developer === developerFilter;

    const matchesType =
      typeFilter === "All Types" || activity.type === typeFilter;

    const matchesRepo =
      repoFilter === "All Repositories" || activity.repository === repoFilter;

    return matchesSearch && matchesDeveloper && matchesType && matchesRepo;
  });
  const sortedFilteredActivities = [...filteredActivities].sort(
    (a: any, b: any) => {
      if (sortOrder === "latest") {
        return (
          new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
        );
      }

      return new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime();
    },
  );

  // Calculate pagination values
  //const totalEvents = filteredActivities.length;
  const indexOfLastEvent = currentPage * eventsPerPage;
  const indexOfFirstEvent = indexOfLastEvent - eventsPerPage;

  const currentActivities = sortedFilteredActivities.slice(
    indexOfFirstEvent,
    indexOfLastEvent,
  );

  const totalPages = Math.ceil(filteredActivities.length / eventsPerPage);
  const getVisiblePages = () => {
    if (totalPages <= 7) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }

    // Beginning pages
    if (currentPage <= 3) {
      return [1, 2, 3, "...", totalPages - 1, totalPages];
    }

    // Ending pages
    if (currentPage >= totalPages - 2) {
      return [1, 2, "...", totalPages - 2, totalPages - 1, totalPages];
    }

    // Middle pages
    return [
      1,
      "...",
      currentPage - 1,
      currentPage,
      currentPage + 1,
      "...",
      totalPages,
    ];
  };

  return (
    <div className="p-8 text-slate-100">
      <h1 className="text-4xl font-bold">Activity</h1>

      <p className="text-slate-400 mt-1">
        All engineering events across repositories
      </p>

      <p className="text-slate-400">
        Showing {(currentPage - 1) * eventsPerPage + 1} -
        {Math.min(currentPage * eventsPerPage, filteredActivities.length)} of{" "}
        {filteredActivities.length} events (Page {currentPage} of {totalPages})
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
          <span className="text-slate-400">Developer:</span>

          <select
            value={developerFilter}
            onChange={(e) => setDeveloperFilter(e.target.value)}
            className="bg-transparent outline-none"
          >
            <option value="All Developers" className="text-black">
              All Developers
            </option>

            {[...new Set(activities.map((a) => a.developer))].map(
              (developer) => (
                <option
                  key={developer}
                  value={developer}
                  className="text-black"
                >
                  {developer}
                </option>
              ),
            )}
          </select>
        </div>

        {/* Type */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-400">Activity Type:</span>
          <select
            value={typeFilter}
            onChange={(e) => setTypeFilter(e.target.value)}
            className="bg-transparent outline-none"
          >
            <option value="All Types" className="text-black">
              All Types
            </option>

            {[...new Set(activities.map((a) => a.type))].map((type) => {
              let displayType = type;
              if (type === "git_commit") displayType = "Git Commit";
              else if (type === "pull_request_closed")
                displayType = "PR Closed";
              else if (type === "open_issue") displayType = "Issue Opened";
              else if (type === "task_completed") displayType = "Issue Closed";
              return (
                <option key={type} value={type} className="text-black">
                  {displayType}
                </option>
              );
            })}
          </select>
        </div>

        {/* Repository */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-300">Repository:</span>

          <select
            value={repoFilter}
            onChange={(e) => setRepoFilter(e.target.value)}
            className="bg-transparent outline-none flex-1 min-w-0"
          >
            <option value="All Repositories" className="text-black">
              All Repositories
            </option>

            {[...new Set(activities.map((a) => a.repository))].map((repo) => (
              <option key={repo} value={repo} className="text-black">
                {repo}
              </option>
            ))}
          </select>
        </div>
        {/* Sort */}
        <div className="flex items-center gap-2 border border-slate-700 rounded-full px-4 py-2 bg-slate-900">
          <span className="text-slate-400">Sort:</span>

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
        <h2 className="text-xl font-semibold">Activity Events</h2>

        <p className="text-slate-400">
          Total Events:
          <span className="text-white font-bold ml-2">
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
                <th className="p-4">Time</th>
                <th className="p-4">Developer</th>
                <th className="p-4">Activity</th>
                <th className="p-4">Repository</th>
                <th className="p-4">Message</th>
                <th className="p-4 text-right">Details</th>
              </tr>
            </thead>

            <tbody>
              {currentActivities.map((activity) => (
                <React.Fragment key={activity.id}>
                  <tr className="border-t border-slate-700 hover:bg-slate-800">
                    {/* Time */}
                    <td className="p-4 text-sm text-slate-400">
                      {activity.displayTime}
                    </td>

                    {/* Developer */}
                    <td className="p-4">
                      <span className="text-sm font-semibold text-slate-100 block">
                        {activity.developer}
                      </span>
                    </td>

                    {/* Activity Type */}
                    <td className="p-4">
                      <span
                        className={`inline-block px-2 py-1 rounded text-xs font-bold uppercase ${
                          activity.type === "git_commit"
                            ? "bg-blue-600/20 text-blue-400 border border-blue-500/20"
                            : activity.type === "pull_request_closed"
                              ? "bg-green-600/20 text-green-400 border border-green-500/20"
                              : activity.type === "open_issue"
                                ? "bg-orange-600/20 text-orange-400 border border-orange-500/20"
                                : activity.type === "task_completed"
                                  ? "bg-rose-600/20 text-rose-400 border border-rose-500/20"
                                  : "bg-slate-700 text-slate-300"
                        }`}
                      >
                        {activity.type === "git_commit"
                          ? "Commit"
                          : activity.type === "pull_request_closed"
                            ? "PR Closed"
                            : activity.type === "open_issue"
                              ? "Issue Opened"
                              : activity.type === "task_completed"
                                ? "Issue Closed"
                                : activity.type}
                      </span>
                    </td>

                    {/* Repository */}
                    <td className="p-4 text-sm text-slate-300">
                      {activity.repository}
                    </td>

                    {/* Message */}
                    <td className="p-4 text-sm text-slate-300 max-w-[300px] truncate">
                      {activity.message.length > 40
                        ? activity.message.substring(0, 40) + "..."
                        : activity.message}
                    </td>

                    {/* Expand Button */}
                    <td className="p-4 text-right">
                      <button
                        onClick={() =>
                          setExpandedId(
                            expandedId === activity.id ? null : activity.id,
                          )
                        }
                        className="text-slate-400 hover:text-white transition-colors"
                      >
                        {expandedId === activity.id ? "▲" : "▼"}
                      </button>
                    </td>
                  </tr>

                  {/* Expanded Row */}
                  {expandedId === activity.id && (
                    <tr className="bg-slate-900">
                      <td colSpan={6} className="p-4">
                        <div className="flex flex-col gap-2 text-sm">
                          <div>
                            <span className="font-bold text-slate-200">
                              Full Message:
                            </span>{" "}
                            <span className="text-slate-400">
                              {activity.message}
                            </span>
                          </div>

                          <div>
                            <span className="font-bold text-slate-200">
                              Repository:
                            </span>{" "}
                            <span className="text-slate-400">
                              {activity.repository}
                            </span>
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

      {/* Pagination */}
      <div className="flex justify-center items-center gap-2 mt-4">
        <button
          disabled={currentPage === 1}
          onClick={() => setCurrentPage(currentPage - 1)}
          className="px-3 py-2 bg-slate-700 rounded disabled:opacity-50"
        >
          &lt;
        </button>

        {getVisiblePages().map((page, index) =>
          page === "..." ? (
            <span key={`dots-${index}`} className="px-2 text-gray-400">
              ...
            </span>
          ) : (
            <button
              key={page}
              onClick={() => setCurrentPage(Number(page))}
              className={`px-3 py-2 rounded ${
                currentPage === page
                  ? "bg-white text-black"
                  : "bg-slate-700 hover:bg-slate-600"
              }`}
            >
              {page}
            </button>
          ),
        )}

        <button
          disabled={currentPage === totalPages}
          onClick={() => setCurrentPage(currentPage + 1)}
          className="px-3 py-2 bg-slate-700 rounded disabled:opacity-50"
        >
          &gt;
        </button>
      </div>
    </div>
  );
}
