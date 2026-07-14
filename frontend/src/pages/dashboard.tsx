import { useEffect, useState } from "react";
import Card from "../components/card";
import {
  GitCommit,
  TrendingUp,
  CheckCircle,
  AlertCircle,
  Users,
  Activity,
  Calendar,
  Layers,
  ArrowRight,
} from "lucide-react";

interface DashboardData {
  stats: {
    total_commits: number;
    velocity_score: number;
    tasks_resolved: number;
    active_blockers: number;
  };
  commit_trend: Array<{
    week: string;
    commits: number;
  }>;
  activity_breakdown: {
    git_commits: number;
    pull_requests_closed: number;
    tasks_resolved: number;
    active_blockers: number;
  };
  top_contributors: Array<{
    user_id: string;
    name: string;
    commits: number;
  }>;
  recent_activity: Array<{
    timestamp: string;
    developer: string;
    type: string;
    repository: string;
    message: string;
    payload: any;
  }>;
}

function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hoveredPoint, setHoveredPoint] = useState<number | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("app_token");

    if (!token) {
      window.location.href = "/login";
      return;
    }

    fetch("http://localhost:8080/api/v1/dashboard", {
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
          throw new Error("Failed to fetch dashboard data");
        }
        return res.json();
      })
      .then((json: DashboardData) => {
        setData(json);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Dashboard error:", err);
        setError(err.message || "An error occurred");
        setLoading(false);
      });
  }, []);

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center bg-slate-950 text-white">
        <div className="flex flex-col items-center gap-4">
          <div className="h-12 w-12 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
          <p className="text-slate-400 animate-pulse font-medium">
            Loading Dashboard metrics...
          </p>
        </div>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex h-full items-center justify-center bg-slate-950 text-white p-6">
        <div className="max-w-md text-center border border-rose-500/20 bg-rose-950/20 rounded-2xl p-8 backdrop-blur-md">
          <AlertCircle className="h-12 w-12 text-rose-500 mx-auto mb-4" />
          <h2 className="text-xl font-bold text-rose-400 mb-2">
            Error Loading Dashboard
          </h2>
          <p className="text-slate-400 mb-6">
            {error || "No dashboard data available."}
          </p>
          <button
            onClick={() => window.location.reload()}
            className="px-6 py-2 bg-rose-600 hover:bg-rose-500 text-white rounded-lg transition-colors font-medium shadow-lg shadow-rose-900/30"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  // Circular donut chart calculations
  const totalBreakdown =
    (data.activity_breakdown.git_commits || 0) +
      (data.activity_breakdown.pull_requests_closed || 0) +
      (data.activity_breakdown.tasks_resolved || 0) +
      (data.activity_breakdown.active_blockers || 0) || 1;

  const breakdownCategories = [
    {
      name: "Git Commits",
      value: data.activity_breakdown.git_commits || 0,
      color: "from-blue-600 to-blue-500",
      hexColor: "#2563eb",
      icon: GitCommit,
    },
    {
      name: "PRs Closed",
      value: data.activity_breakdown.pull_requests_closed || 0,
      color: "from-green-600 to-green-500",
      hexColor: "#16a34a",
      icon: Layers,
    },
    {
      name: "Tasks Resolved",
      value: data.activity_breakdown.tasks_resolved || 0,
      color: "from-pink-600 to-pink-500",
      hexColor: "#db2777",
      icon: CheckCircle,
    },
    {
      name: "Active Blockers",
      value: data.activity_breakdown.active_blockers || 0,
      color: "from-orange-600 to-orange-500",
      hexColor: "#ea580c",
      icon: AlertCircle,
    },
  ];

  const donutRadius = 50;
  const donutCircumference = 2 * Math.PI * donutRadius;
  let accumulatedPercent = 0;

  // Commit trend calculations
  const maxCommitsInTrend = Math.max(
    ...data.commit_trend.map((t) => t.commits),
    1,
  );

  const chartWidth = 600;
  const chartHeight = 240;
  const paddingX = 40;
  const paddingY = 30;

  const points = data
    ? data.commit_trend.map((item, i) => {
        const x =
          paddingX +
          i *
            ((chartWidth - 2 * paddingX) /
              Math.max(data.commit_trend.length - 1, 1));
        const y =
          chartHeight -
          paddingY -
          (item.commits / maxCommitsInTrend) * (chartHeight - 2 * paddingY);
        return { x, y, week: item.week, commits: item.commits };
      })
    : [];

  const linePath = points
    .map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`)
    .join(" ");
  const areaPath =
    points.length > 0
      ? `${linePath} L ${points[points.length - 1].x} ${chartHeight - paddingY} L ${points[0].x} ${chartHeight - paddingY} Z`
      : "";

  return (
    <div className="bg-slate-950 min-h-screen p-8 text-white">
      {/* Header */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-4xl font-bold tracking-tight bg-gradient-to-r from-white via-slate-100 to-slate-400 bg-clip-text text-transparent">
            Dashboard
          </h1>
          <p className="text-slate-400 mt-1">
            Real-time engineering operations & performance analytics.
          </p>
        </div>
        <div className="flex items-center gap-2 bg-slate-900 border border-slate-800 rounded-xl px-4 py-2 text-sm text-slate-400">
          <Activity className="h-4 w-4 text-emerald-500 animate-pulse" />
          <span>Real-time Sync Active</span>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {/* Total Commits */}
        <Card className="p-6 border-slate-800 hover:border-blue-500/50 hover:bg-slate-900/60 transition-all duration-300 relative overflow-hidden group">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-blue-500/10 border border-blue-500/20 rounded-xl text-blue-400">
              <GitCommit className="h-6 w-6" />
            </div>
            <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-blue-500/10 text-blue-400">
              Commits
            </span>
          </div>
          <p className="text-slate-400 text-sm font-medium">Total Commits</p>
          <h3 className="text-3xl font-bold mt-1 text-white tracking-tight">
            {data.stats.total_commits.toLocaleString()}
          </h3>
        </Card>

        {/* Velocity Score */}
        <Card className="p-6 border-slate-800 hover:border-violet-500/50 hover:bg-slate-900/60 transition-all duration-300 relative overflow-hidden group">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-violet-500/10 border border-violet-500/20 rounded-xl text-violet-400">
              <TrendingUp className="h-6 w-6" />
            </div>
            <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-violet-500/10 text-violet-400">
              Score
            </span>
          </div>
          <p className="text-slate-400 text-sm font-medium">Velocity Score</p>
          <h3 className="text-3xl font-bold mt-1 text-white tracking-tight">
            {data.stats.velocity_score}
            <span className="text-lg text-slate-500 font-normal">/100</span>
          </h3>
        </Card>

        {/* Tasks Resolved */}
        <Card className="p-6 border-slate-800 hover:border-pink-500/50 hover:bg-slate-900/60 transition-all duration-300 relative overflow-hidden group">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-pink-500/10 border border-pink-500/20 rounded-xl text-pink-400">
              <CheckCircle className="h-6 w-6" />
            </div>
            <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-pink-500/10 text-pink-400">
              Completed
            </span>
          </div>
          <p className="text-slate-400 text-sm font-medium">Tasks Resolved</p>
          <h3 className="text-3xl font-bold mt-1 text-white tracking-tight">
            {data.stats.tasks_resolved.toLocaleString()}
          </h3>
        </Card>

        {/* Active Blockers */}
        <Card className="p-6 border-slate-800 hover:border-orange-500/50 hover:bg-slate-900/60 transition-all duration-300 relative overflow-hidden group">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-orange-500/10 border border-orange-500/20 rounded-xl text-orange-400">
              <AlertCircle className="h-6 w-6" />
            </div>
            <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-orange-500/10 text-orange-400">
              Issues
            </span>
          </div>
          <p className="text-slate-400 text-sm font-medium">Active Blockers</p>
          <h3 className="text-3xl font-bold mt-1 text-white tracking-tight">
            {data.stats.active_blockers.toLocaleString()}
          </h3>
        </Card>
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
        {/* Commit Trend Chart */}
        <Card className="lg:col-span-2 p-6 border-slate-800 flex flex-col justify-between">
          <div>
            <div className="flex justify-between items-center mb-6">
              <div className="flex items-center gap-2">
                <Calendar className="h-5 w-5 text-indigo-400" />
                <h2 className="text-lg font-semibold">Commit Trend</h2>
              </div>
              <span className="text-xs text-slate-500">
                Weekly Commits Grouping
              </span>
            </div>

            {/* Custom SVG Line Chart */}
            <div className="relative w-full h-60">
              {data.commit_trend.length === 0 ? (
                <div className="w-full flex items-center justify-center text-slate-500 text-sm h-full">
                  No commit trend data found
                </div>
              ) : (
                <>
                  <svg
                    className="w-full h-full"
                    viewBox={`0 0 ${chartWidth} ${chartHeight}`}
                    preserveAspectRatio="none"
                  >
                    <defs>
                      <linearGradient id="lineGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="#3b82f6" />
                        <stop offset="100%" stopColor="#8b5cf6" />
                      </linearGradient>
                      <linearGradient id="areaGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop
                          offset="0%"
                          stopColor="#3b82f6"
                          stopOpacity="0.3"
                        />
                        <stop
                          offset="100%"
                          stopColor="#3b82f6"
                          stopOpacity="0.0"
                        />
                      </linearGradient>
                    </defs>

                    {/* Grid lines (horizontal) */}
                    {[0, 0.25, 0.5, 0.75, 1].map((ratio, idx) => {
                      const y = paddingY + ratio * (chartHeight - 2 * paddingY);
                      return (
                        <line
                          key={idx}
                          x1={paddingX}
                          y1={y}
                          x2={chartWidth - paddingX}
                          y2={y}
                          stroke="#334155"
                          strokeWidth="1"
                          strokeDasharray="4 4"
                        />
                      );
                    })}

                    {/* Area path */}
                    {areaPath && <path d={areaPath} fill="url(#areaGrad)" />}

                    {/* Line path */}
                    {linePath && (
                      <path
                        d={linePath}
                        fill="none"
                        stroke="url(#lineGrad)"
                        strokeWidth="3"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      />
                    )}

                    {/* Interactive dots and hover trigger circles */}
                    {points.map((p, idx) => (
                      <g key={idx}>
                        {/* Dash vertical guide line on hover */}
                        {hoveredPoint === idx && (
                          <line
                            x1={p.x}
                            y1={paddingY}
                            x2={p.x}
                            y2={chartHeight - paddingY}
                            stroke="#64748b"
                            strokeWidth="1.5"
                            strokeDasharray="3 3"
                          />
                        )}

                        {/* Point circle */}
                        <circle
                          cx={p.x}
                          cy={p.y}
                          r={hoveredPoint === idx ? "6" : "4"}
                          fill={hoveredPoint === idx ? "#60a5fa" : "#3b82f6"}
                          stroke="#0f172a"
                          strokeWidth="2"
                          className="transition-all duration-150"
                        />

                        {/* Large invisible circle for hover trigger */}
                        <circle
                          cx={p.x}
                          cy={p.y}
                          r="16"
                          fill="transparent"
                          className="cursor-pointer"
                          onMouseEnter={() => setHoveredPoint(idx)}
                          onMouseLeave={() => setHoveredPoint(null)}
                        />
                      </g>
                    ))}
                  </svg>

                  {/* Dynamic Floating Tooltip */}
                  {hoveredPoint !== null && points[hoveredPoint] && (
                    <div
                      className="absolute bg-slate-900 border border-slate-700 text-xs px-3 py-2 rounded-xl shadow-2xl z-20 pointer-events-none transform -translate-x-1/2 -translate-y-full transition-all duration-150"
                      style={{
                        left: `${(points[hoveredPoint].x / chartWidth) * 100}%`,
                        top: `${(points[hoveredPoint].y / chartHeight) * 100 - 4}%`,
                      }}
                    >
                      <div className="font-bold text-white text-center">
                        {points[hoveredPoint].commits} commits
                      </div>
                      <div className="text-[10px] text-slate-400 text-center">
                        Week of {points[hoveredPoint].week}
                      </div>
                    </div>
                  )}

                  {/* X-Axis labels at the bottom */}
                  <div className="flex justify-between px-8 mt-2">
                    {points.map((p, idx) => (
                      <span
                        key={idx}
                        className="text-[10px] text-slate-500 font-semibold tracking-wider"
                      >
                        {p.week}
                      </span>
                    ))}
                  </div>
                </>
              )}
            </div>
          </div>
        </Card>

        {/* Activity Breakdown Donut Chart */}
        <Card className="p-6 border-slate-800 flex flex-col justify-between">
          <div>
            <div className="flex items-center gap-2 mb-6">
              <Layers className="h-5 w-5 text-pink-400" />
              <h2 className="text-lg font-semibold">Activity Breakdown</h2>
            </div>

            {/* Circular Donut Chart */}
            <div className="flex flex-col items-center justify-center my-4">
              <div className="relative h-44 w-44">
                <svg className="w-full h-full" viewBox="0 0 160 160">
                  <circle
                    cx="80"
                    cy="80"
                    r={donutRadius}
                    fill="transparent"
                    stroke="#1e293b"
                    strokeWidth="12"
                  />
                  {breakdownCategories.map((cat) => {
                    const percent = (cat.value / totalBreakdown) * 100;
                    const strokeDasharray = `${(percent / 100) * donutCircumference} ${donutCircumference}`;
                    const strokeDashoffset = -(
                      (accumulatedPercent / 100) *
                      donutCircumference
                    );
                    accumulatedPercent += percent;
                    return (
                      <circle
                        key={cat.name}
                        cx="80"
                        cy="80"
                        r={donutRadius}
                        fill="transparent"
                        stroke={cat.hexColor}
                        strokeWidth="12"
                        strokeDasharray={strokeDasharray}
                        strokeDashoffset={strokeDashoffset}
                        transform="rotate(-90 80 80)"
                        className="transition-all duration-700 ease-out hover:stroke-[14px] cursor-pointer"
                      />
                    );
                  })}
                </svg>
                <div className="absolute inset-0 flex flex-col items-center justify-center">
                  <span className="text-[10px] text-slate-500 uppercase tracking-widest">
                    Total Events
                  </span>
                  <span className="text-2xl font-bold">
                    {totalBreakdown.toLocaleString()}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Categories Legend */}
          <div className="grid grid-cols-2 gap-3 mt-4 pt-4 border-t border-slate-800/60">
            {breakdownCategories.map((cat) => (
              <div key={cat.name} className="flex items-center gap-2">
                <span
                  className={`h-2.5 w-2.5 rounded-full bg-gradient-to-r ${cat.color}`}
                />
                <div className="flex flex-col">
                  <span className="text-xs text-slate-400 font-medium">
                    {cat.name}
                  </span>
                  <span className="text-[11px] text-slate-500 font-bold">
                    {cat.value} (
                    {Math.round((cat.value / totalBreakdown) * 100)}%)
                  </span>
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>

      {/* Grid Bottom: Top Contributors and Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Top Contributors */}
        <Card className="p-6 border-slate-800 flex flex-col justify-between h-full">
          <div>
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center gap-2">
                <Users className="h-5 w-5 text-emerald-400" />
                <h2 className="text-lg font-semibold">Top Contributors</h2>
              </div>
              <span className="text-xs text-slate-500">Commits Count</span>
            </div>

            <div className="space-y-4">
              {data.top_contributors.length === 0 ? (
                <div className="text-slate-500 text-sm py-8 text-center">
                  No contributor records found
                </div>
              ) : (
                data.top_contributors.map((contrib, index) => (
                  <div
                    key={contrib.user_id}
                    className="flex justify-between items-center p-3 rounded-xl hover:bg-slate-900 border border-transparent hover:border-slate-800 transition-all duration-300"
                  >
                    <div className="flex items-center gap-3">
                      <div className="flex items-center justify-center h-8 w-8 rounded-full bg-slate-800 font-bold text-sm text-slate-300 border border-slate-700">
                        #{index + 1}
                      </div>
                      <div>
                        <p className="font-semibold text-sm text-white">
                          {contrib.name}
                        </p>
                        <p className="text-[11px] text-slate-500">Developer</p>
                      </div>
                    </div>
                    <span className="text-xs font-bold px-3 py-1 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20">
                      {contrib.commits} commits
                    </span>
                  </div>
                ))
              )}
            </div>
          </div>
        </Card>

        {/* Recent Activity */}
        <Card className="lg:col-span-2 p-6 border-slate-800">
          <div className="flex justify-between items-center mb-6">
            <div className="flex items-center gap-2">
              <Activity className="h-5 w-5 text-teal-400" />
              <h2 className="text-lg font-semibold">Recent Activity</h2>
            </div>
            <a
              href="/activity"
              className="text-xs text-blue-400 hover:text-blue-300 flex items-center gap-1 group font-medium"
            >
              View Full Stream
              <ArrowRight className="h-3 w-3 group-hover:translate-x-1 transition-transform" />
            </a>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b border-slate-800 text-slate-400 text-xs font-bold">
                  <th className="pb-3 pr-4">Developer</th>
                  <th className="pb-3 pr-4">Activity</th>
                  <th className="pb-3 pr-4">Repository</th>
                  <th className="pb-3">Details</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/50">
                {data.recent_activity.length === 0 ? (
                  <tr>
                    <td
                      colSpan={4}
                      className="py-8 text-slate-500 text-sm text-center"
                    >
                      No recent activities recorded
                    </td>
                  </tr>
                ) : (
                  data.recent_activity.map((activity, idx) => (
                    <tr
                      key={idx}
                      className="hover:bg-slate-900/40 transition-colors"
                    >
                      <td className="py-3.5 pr-4">
                        <span className="text-sm font-semibold text-white block">
                          {activity.developer}
                        </span>
                        <span className="text-[10px] text-slate-500">
                          {new Date(activity.timestamp).toLocaleString()}
                        </span>
                      </td>
                      <td className="py-3.5 pr-4">
                        <span
                          className={`inline-block px-2 py-0.5 rounded text-[10px] font-bold text-white uppercase ${
                            activity.type === "git_commit"
                              ? "bg-blue-600/20 text-blue-400 border border-blue-500/20"
                              : activity.type === "pull_request_closed"
                                ? "bg-green-600/20 text-green-400 border border-green-500/20"
                                : activity.type === "open_issue"
                                  ? "bg-orange-600/20 text-orange-400 border border-orange-500/20"
                                  : activity.type === "task_completed"
                                    ? "bg-rose-600/20 text-rose-400 border border-rose-500/20"
                                    : "bg-slate-700"
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
                      <td className="py-3.5 pr-4 text-xs font-semibold text-slate-400">
                        {activity.repository}
                      </td>
                      <td className="py-3.5 text-xs text-slate-300 max-w-[240px] truncate">
                        {activity.message}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      </div>
    </div>
  );
}

export default Dashboard;
