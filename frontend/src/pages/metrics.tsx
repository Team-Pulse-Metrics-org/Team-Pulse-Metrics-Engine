import { useEffect, useState } from "react";
import { MetricBarChart } from "../components/metricBarChart";
import { AlertCircle, Calendar, Users } from "lucide-react";

interface MetricCoordinate {
  label: string;
  value: number;
}

interface ChartTimeline {
  weekly: MetricCoordinate[];
  monthly: MetricCoordinate[];
}

interface UnifiedMetricsResponse {
  commits: ChartTimeline;
  velocity_score: ChartTimeline;
  tasks_resolved: ChartTimeline;
  open_issues: ChartTimeline;
}

type TimeFrame = "weekly" | "monthly";

function Metrics() {
  const [metricsData, setMetricsData] = useState<UnifiedMetricsResponse | null>(
    null,
  );
  const [Loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [timeframe, setTimeframe] = useState<TimeFrame>("weekly");

  useEffect(() => {
    const token = localStorage.getItem("app_token");

    if (!token) {
      window.location.href = "/login";
      return;
    }

    fetch("http://localhost:8080/api/v1/metrics", {
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
          throw new Error("Session expired. Please login again.");
        }
        if (!res.ok) {
          throw new Error("failed too fetch team metrics snapshot.");
        }
        return res.json();
      })
      .then((data: UnifiedMetricsResponse) => {
        setMetricsData(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Metrics engine load error:", err);
        setError(
          err.message || "An unexpected error occurred loading analytics.",
        );
        setLoading(false);
      });
  }, []);

  if (Loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-950 text-white w-full">
        <div className="flex flex-col items-center gap-4">
          <div className="h-12 w-12 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
          <p className="text-slate-400 animate-pulse font-medium">
            Loading team metrics stream...
          </p>
        </div>
      </div>
    );
  }

  if (error || !metricsData) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-950 text-white p-6 w-full">
        <div className="max-w-md text-center border border-rose-500/20 bg-rose-950/20 rounded-2xl p-8 backdrop-blur-md">
          <AlertCircle className="h-12 w-12 text-rose-500 mx-auto mb-4" />
          <h2 className="text-xl font-bold text-rose-400 mb-2">
            Error Loading Metrics
          </h2>
          <p className="text-slate-400 mb-6">
            {error || "No data stream available right now."}
          </p>
          <button
            onClick={() => window.location.reload()}
            className="px-6 py-2 bg-rose-600 hover:bg-rose-500 text-white rounded-lg transition-colors font-medium shadow-lg"
          >
            Retry Connection
          </button>
        </div>
      </div>
    );
  }

  const activePeriodLabel = timeframe === "weekly" ? "Weekly" : "Monthly";

  return (
    <div className="bg-slate-950 min-h-screen p-8 text-slate-100 w-full">
      {/* Primary Header Stack */}
      <div className="flex flex-col gap-6 mb-8">
        <div>
          <h1 className="text-4xl font-bold tracking-tight bg-linear-to-r from-slate-50 via-slate-100 to-slate-400 bg-clip-text text-transparent">
            Metrics
          </h1>
          <p className="text-slate-400 mt-1">
            Engineering performance analytics
          </p>
        </div>

        {/* Filter Controls Row */}
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 bg-slate-900 border border-slate-800 rounded-xl px-4 py-2 text-sm text-slate-300">
              <Users className="h-4 w-4 text-slate-400" />
              <select className="bg-transparent focus:outline-none cursor-pointer pr-2">
                <option value="all">Entire Team</option>
              </select>
            </div>
          </div>

          {/* Timeframe Toggle Switch */}
          <div className="bg-slate-900/60 border border-slate-800 p-1 rounded-xl flex items-center gap-1 text-xs font-medium">
            <button
              onClick={() => setTimeframe("weekly")}
              className={`px-3 py-1.5 rounded-lg transition-all ${
                timeframe === "weekly"
                  ? "bg-blue-600 text-white shadow-md"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              Weekly
            </button>
            <button
              onClick={() => setTimeframe("monthly")}
              className={`px-3 py-1.5 rounded-lg transition-all ${
                timeframe === "monthly"
                  ? "bg-blue-600 text-white shadow-md"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              Monthly
            </button>
          </div>
        </div>
      </div>
      {/* Grid */}
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-8 w-full">
        {/* 1. Commits Card */}
        <MetricBarChart
          title={`${activePeriodLabel} Team Commits`}
          subtitle="Total volume of codebase contributions"
          data={metricsData.commits[timeframe]}
          valueLabel="commits"
          color="#3b82f6" // Blue
        />

        {/* 2. Velocity Score Card */}
        <MetricBarChart
          title={`${activePeriodLabel} Velocity Score`}
          subtitle="Averaged baseline task sizing execution velocity"
          data={metricsData.velocity_score[timeframe]}
          valueLabel="pts"
          color="#06b6d4" // Cyan
          yAxisMax={100} // Keep scoring range inside fixed ceilings
        />

        {/* 3. Tasks Resolved Card */}
        <MetricBarChart
          title={`${activePeriodLabel} Tasks Completed`}
          subtitle="Total user tickets moved directly to a resolved status state"
          data={metricsData.tasks_resolved[timeframe]}
          valueLabel="tasks"
          color="#10b981" // Emerald
        />

        {/* 4. Open Issues Card */}
        <MetricBarChart
          title={`${activePeriodLabel} Outstanding Open Issues`}
          subtitle="Active total backlog volume pending assignment or fix"
          data={metricsData.open_issues[timeframe]}
          valueLabel="issues"
          color="#ea580c" // Orange
        />
      </div>
    </div>
  );
}

export default Metrics;
