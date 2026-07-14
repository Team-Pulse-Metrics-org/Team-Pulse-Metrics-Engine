import type React from "react";
import { useState, useEffect } from "react";
import { FaGithub } from "react-icons/fa";
import { LuSun, LuMoon } from "react-icons/lu";

function handleGitHubLogin() {
  const CLIENT_ID: string = import.meta.env.VITE_GITHUB_CLIENT_ID;
  console.log(CLIENT_ID);

  const REDIRECT_URI: string = "http://localhost:5173/auth/callback";
  const SCOPE: string = "user,public_repo";

  const githubAuthUrl: string = `https://github.com/login/oauth/authorize?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&scope=${SCOPE}`;

  window.location.href = githubAuthUrl;
}

const Login: React.FC = () => {
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
    <div className="bg-slate-950 h-screen w-full overflow-y-auto flex flex-col items-center justify-center p-6 text-slate-100">
      {/* Outer Card Wrapper */}
      <div className="relative max-w-4xl w-full mx-auto h-[600px] rounded-3xl overflow-hidden shadow-2xl border border-slate-800 bg-slate-900 grid grid-cols-1 md:grid-cols-2">
        
        {/* Left Side: Graphic / Image (Theme-aware graphic) */}
        <div className={`relative hidden md:flex flex-col justify-between p-12 overflow-hidden select-none transition-colors duration-300
          ${theme === "dark" 
            ? "bg-black" 
            : "bg-gradient-to-br from-indigo-50 via-purple-50 to-cyan-50"
          }`}
        >
          {/* Background Image - Shown in both modes with theme-appropriate opacity */}
          <img
            src="rocket.jpg"
            alt="Rocket launch"
            className={`absolute inset-0 w-full h-full object-cover z-0 transition-all duration-300
              ${theme === "dark" ? "opacity-60 mix-blend-luminosity" : "opacity-25"}`}
          />

          {/* Theme-Aware Overlay Gradient */}
          <div className={`absolute inset-0 z-10 transition-all duration-300
            ${theme === "dark" 
              ? "bg-gradient-to-tr from-indigo-950/95 via-slate-950/70 to-indigo-950/30" 
              : "bg-gradient-to-tr from-indigo-50/90 via-white/85 to-purple-50/90"}`}
          />

          {/* Dot Grid Overlay */}
          <div className={`absolute inset-0 z-10 transition-all duration-300
            ${theme === "dark"
              ? "bg-[radial-gradient(rgba(255,255,255,0.05)_1px,transparent_1px)]"
              : "bg-[radial-gradient(rgba(99,102,241,0.08)_1px,transparent_1px)]"} 
            [background-size:16px_16px]`} 
          />

          {/* Left panel contents. Uses theme-aware text colors */}
          <div className="relative z-20 flex flex-col gap-2">
            <div className={`text-3xl font-extrabold tracking-tight flex items-center gap-2
              ${theme === "dark" ? "text-white" : "text-indigo-950"}`}
            >
              <span className="text-cyan-500">Team</span>
              <span className="text-purple-600">Pulse</span>
            </div>
            <span className={`text-xs font-semibold tracking-wider uppercase
              ${theme === "dark" ? "text-cyan-400/80" : "text-indigo-600/80"}`}
            >
              Engineering Metrics Engine
            </span>
          </div>

          <div className="relative z-20 mt-auto flex flex-col gap-4">
            <h2 className={`text-3xl font-bold leading-tight
              ${theme === "dark" ? "text-white" : "text-indigo-950"}`}
            >
              Launch engineering productivity to new heights.
            </h2>
            <p className={`text-sm leading-relaxed
              ${theme === "dark" ? "text-slate-300" : "text-slate-700"}`}
            >
              Track commits, analyze pull requests, monitor active tasks, and identify team blockers instantly.
            </p>
            
            {/* Visual Stats Indicators */}
            <div className={`flex gap-4 mt-4 pt-4 border-t text-xs
              ${theme === "dark" ? "border-white/10 text-slate-300" : "border-slate-200 text-slate-700"}`}
            >
              <div className={`flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg border
                ${theme === "dark" ? "bg-white/5 border-white/10" : "bg-white/70 border-slate-200"}`}
              >
                <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                Real-time Sync
              </div>
              <div className={`flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg border
                ${theme === "dark" ? "bg-white/5 border-white/10" : "bg-white/70 border-slate-200"}`}
              >
                <span>🔐</span> Secure GitHub Oauth
              </div>
            </div>
          </div>
        </div>

        {/* Right Side: Auth Form (Theme-sensitive) */}
        <div className="relative flex flex-col justify-between p-12 bg-transparent z-25">
          
          {/* Top-Right Theme Toggle */}
          <div className="absolute top-8 right-8">
            <button
              onClick={toggleTheme}
              className="p-3 rounded-xl border border-slate-800 text-slate-400 hover:bg-slate-900 hover:text-slate-100 transition cursor-pointer"
              aria-label="Toggle theme"
            >
              {theme === "dark" ? <LuSun size={18} /> : <LuMoon size={18} />}
            </button>
          </div>

          {/* Center Content */}
          <div className="my-auto flex flex-col gap-8 w-full max-w-sm mx-auto">
            <div className="flex flex-col gap-2">
              <h1 className="font-bold text-4xl tracking-tight" style={{ color: "var(--slate-50)" }}>Sign In</h1>
              <p className="text-sm" style={{ color: "var(--slate-400)" }}>
                Access your productivity metrics by authenticating with GitHub.
              </p>
            </div>

            <div className="flex flex-col gap-4">
              <button
                className="flex items-center justify-center gap-3 w-full bg-indigo-600 hover:bg-indigo-700 text-white py-3.5 px-6 transition-all rounded-xl font-semibold shadow-lg shadow-indigo-600/10 hover:shadow-indigo-600/20 active:scale-[0.98] group cursor-pointer"
                onClick={handleGitHubLogin}
                title="Sign in with GitHub"
              >
                <FaGithub
                  size={24}
                  className="transition-transform group-hover:scale-110"
                />
                <span>Continue with GitHub</span>
              </button>
              
              <div className="flex items-center gap-2 text-xs justify-center" style={{ color: "var(--slate-500)" }}>
                <span>🔒</span>
                <span>Authorized developer access only</span>
              </div>
            </div>
          </div>

          {/* Bottom Footer Notice */}
          <div className="text-center text-[10px] font-medium tracking-wide w-full" style={{ color: "var(--slate-500)" }}>
            By signing in, you authorize Team Pulse to synchronize public repository details and developer metrics.
          </div>
        </div>

      </div>

      {/* Footer Copy outside the main card */}
      <p className="mt-8 text-center text-xs font-medium tracking-wide" style={{ color: "var(--slate-500)" }}>
        © 2026 Team Pulse · Engineering Analytics Platform
      </p>
    </div>
  );
};

export default Login;
