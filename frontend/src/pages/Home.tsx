import { Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { LuSun, LuMoon } from "react-icons/lu";

export default function Home() {
  const [theme, setTheme] = useState(() => {
    if (typeof window !== "undefined") {
      return localStorage.getItem("app_theme") || "dark";
    }
    return "dark";
  });

  useEffect(() => {
    const root = window.document.documentElement;
    if (theme === "dark") {
      root.classList.add("dark");
    } else {
      root.classList.remove("dark");
    }
    localStorage.setItem("app_theme", theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === "dark" ? "light" : "dark"));
  };

  return (
    <div className="bg-slate-950 h-screen w-full overflow-y-auto text-slate-100 scroll-smooth">
      {/* Navbar */}
      <nav className="w-full flex justify-between items-center px-16 py-6 border-b border-slate-800">
        {/* Logo */}
        <div className="text-3xl font-bold">
          <span className="text-cyan-400">Team</span>
          <span className="text-violet-500">Pulse</span>
        </div>

        {/* Navigation Links */}
        <div className="flex items-center gap-8 text-lg">
          <a
            href="#features"
            className="hover:text-cyan-400 text-slate-300 transition duration-200 cursor-pointer"
          >
            Features
          </a>

          <a
            href="#how-it-works"
            className="hover:text-cyan-400 text-slate-300 transition duration-200 cursor-pointer"
          >
            How It Works
          </a>

          {/* Theme Toggle */}
          <button
            onClick={toggleTheme}
            className="p-3 rounded-xl border border-slate-800 text-slate-400 hover:bg-slate-900 hover:text-slate-100 transition cursor-pointer"
            aria-label="Toggle theme"
          >
            {theme === "dark" ? <LuSun size={20} /> : <LuMoon size={20} />}
          </button>

          <Link
            to="/login"
            className="bg-cyan-500 hover:bg-cyan-600 text-white px-6 py-3 rounded-xl font-semibold transition"
          >
            Login
          </Link>
        </div>
      </nav>

      {/* Hero Section */} 
      <section className="flex flex-col items-center justify-center text-center mt-20 px-6">
        <h1 className="text-6xl font-extrabold leading-tight max-w-5xl text-slate-50">
          Track Engineering Productivity
          <span className="text-cyan-400"> in Real Time</span>
        </h1>

        <p className="mt-6 text-xl text-slate-400 max-w-3xl">
          Monitor commits, pull requests, completed tasks and
          engineering blockers from a single intelligent platform.
        </p>

        {/* Mini Metrics Pulse Preview (Interactive/Relevant element) */}
        <div className="mt-12 bg-slate-900 border rounded-2xl p-6 max-w-2xl w-full mx-auto shadow-lg text-left grid grid-cols-3 gap-6 select-none">
          <div className="flex flex-col gap-1">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">Team Velocity</span>
            <span className="text-2xl font-bold text-slate-50">84.2%</span>
            <span className="text-[10px] text-emerald-500 flex items-center gap-1">
              ▲ +3.4% this week
            </span>
          </div>
          <div className="flex flex-col gap-1 border-l border-slate-800 pl-6">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">Active Blockers</span>
            <span className="text-2xl font-bold text-rose-500">2 Pending</span>
            <span className="text-[10px] text-slate-500">Avg resolution: 4.2h</span>
          </div>
          <div className="flex flex-col gap-1 border-l border-slate-800 pl-6">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">PR Cycle Time</span>
            <span className="text-2xl font-bold text-cyan-400">1.8 Days</span>
            <span className="text-[10px] text-emerald-500 flex items-center gap-1">
              ▼ -12% improvement
            </span>
          </div>
        </div>

        <div className="mt-10 mb-20">
          <Link
            to="/login"
            className="bg-cyan-500 hover:bg-cyan-600 text-white px-10 py-4 rounded-2xl font-semibold transition"
          >
            Get Started
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="px-10 py-24 border-t border-slate-800">
        <h2 className="text-4xl font-extrabold text-slate-50 text-center mb-4">
          Core Metrics Platform Features
        </h2>
        <p className="text-slate-400 text-center max-w-2xl mx-auto mb-16 text-lg">
          An automated system that turns git history and reviews into key indicators of pipeline health.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 max-w-6xl mx-auto">
          <div className="bg-slate-900 p-8 rounded-2xl border hover:-translate-y-3 hover:border-cyan-400 transition-all duration-300">
            <div className="text-5xl mb-4">💻</div>
            <h3 className="text-2xl font-semibold text-slate-50 mb-3">
              GitHub Ingestion
            </h3>
            <p className="text-slate-400 text-sm leading-relaxed">
              Track developer activity, daily commits, and workspace updates in real-time.
            </p>
          </div>

          <div className="bg-slate-900 p-8 rounded-2xl border hover:-translate-y-3 hover:border-cyan-400 transition-all duration-300">
            <div className="text-5xl mb-4">🔀</div>
            <h3 className="text-2xl font-semibold text-slate-50 mb-3">
              Code Quality
            </h3>
            <p className="text-slate-400 text-sm leading-relaxed">
              Review open pull requests, tracking cycle times and approval histories.
            </p>
          </div>

    <div className="bg-slate-900/70 backdrop-blur-lg p-8
rounded-2xl
border border-slate-700
shadow-2xl
hover:-translate-y-3
hover:shadow-cyan-500/20
hover:border-cyan-400
transition-all duration-300">
      <div className="text-5xl mb-4">✅</div>
      <h3 className="text-2xl font-semibold mb-3">
        Task Metrics
      </h3>
      <p className="text-slate-400">
        Measure completed tasks and delivery velocity.
      </p>
    </div>
 
    <div className="
bg-slate-900/70 backdrop-blur-lg p-8
rounded-2xl
border border-slate-700
shadow-2xl
hover:-translate-y-3
hover:shadow-cyan-500/20
hover:border-cyan-400
transition-all duration-300
">
      <div className="text-5xl mb-4">🚨</div>
      <h3 className="text-2xl font-semibold mb-3">
        Blocker Detection
      </h3>
      <p className="text-slate-400">
        Identify bottlenecks and active engineering blockers.
      </p>
    </div>
  </div>
</section>
      {/* How It Works Section */}
      <section id="how-it-works" className="py-24 border-t border-slate-800 bg-slate-950/40">
        <div className="max-w-6xl mx-auto px-6">
          <h2 className="text-4xl font-extrabold text-slate-50 text-center mb-4">
            How It Works
          </h2>
          <p className="text-slate-400 text-center max-w-2xl mx-auto mb-16 text-lg">
            Connect repositories and capture instant productivity dashboards in minutes.
          </p>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-10">
            {/* Step 1 */}
            <div className="flex flex-col items-center text-center p-8 bg-slate-900 border rounded-2xl relative">
              <div className="w-12 h-12 rounded-full bg-cyan-500/10 border border-cyan-500/30 flex items-center justify-center text-cyan-400 font-bold text-xl mb-6">
                1
              </div>
              <h3 className="text-xl font-bold text-slate-50 mb-3">Connect GitHub</h3>
              <p className="text-slate-400 text-sm leading-relaxed">
                Connect your organization repositories securely via standard GitHub OAuth authorization.
              </p>
            </div>

            {/* Step 2 */}
            <div className="flex flex-col items-center text-center p-8 bg-slate-900 border rounded-2xl relative">
              <div className="w-12 h-12 rounded-full bg-purple-500/10 border border-purple-500/30 flex items-center justify-center text-purple-400 font-bold text-xl mb-6">
                2
              </div>
              <h3 className="text-xl font-bold text-slate-50 mb-3">Gather Statistics</h3>
              <p className="text-slate-400 text-sm leading-relaxed">
                Our ingestion engine parses active branches to calculate commit frequencies and code turnaround.
              </p>
            </div>

            {/* Step 3 */}
            <div className="flex flex-col items-center text-center p-8 bg-slate-900 border rounded-2xl relative">
              <div className="w-12 h-12 rounded-full bg-indigo-500/10 border border-indigo-500/30 flex items-center justify-center text-indigo-400 font-bold text-xl mb-6">
                3
              </div>
              <h3 className="text-xl font-bold text-slate-50 mb-3">Drive Analytics</h3>
              <p className="text-slate-400 text-sm leading-relaxed">
                Get immediate access to engineering statistics, team blockers, and live pipeline analytics.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 border-t border-slate-800 text-center text-slate-500 text-xs">
        <p>© 2026 Team Pulse · Engineering Analytics Platform</p>
      </footer>
    </div>
  );
}