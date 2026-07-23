import type React from "react";
import { useEffect, useRef, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";

interface LoginResponse {
  status: string;
  token: string;
  user: {
    id: string;
    email: string;
    role: string;
  };
}

const AuthCallback: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const exchangeTriggered = useRef(false);

  const [theme] = useState(() => {
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

  useEffect(() => {
    const code = searchParams.get("code");

    if (!code) {
      setIsLoading(false);
      setErrorMessage("GitHub temporary token is invalid");
      return;
    }

    if (!exchangeTriggered.current) {
      exchangeTriggered.current = true;

      fetch("http://localhost:8080/api/v1/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code }),
      })
        .then(async (response) => {
          if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(
              errorData.message ||
                `Server responded with status: ${response.status}`,
            );
          }
          return response.json() as Promise<LoginResponse>;
        })
        .then((data) => {
  localStorage.setItem("app_token", data.token);
  localStorage.setItem("role", data.user.role);
  localStorage.setItem("user_id", data.user.id);
  localStorage.setItem("email", data.user.email);
  navigate("/dashboard");
})
        .catch((error) => {
          console.error("Auth exchange failed", error);
          setErrorMessage(error.message || "Internal Server Error");
          setIsLoading(false);
        });
    }
  }, [searchParams, navigate]);

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
            src="/rocket.jpg"
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

        {/* Right Side: Loading / Error Status (Theme-sensitive) */}
        <div className="relative flex flex-col justify-between p-12 bg-transparent z-25">
          
          {/* Center Content */}
          <div className="my-auto flex flex-col gap-8 w-full max-w-sm mx-auto text-center animate-fade-in">
            {isLoading && !errorMessage && (
              <div className="space-y-6 flex flex-col items-center justify-center">
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-slate-800 border-t-indigo-500 mx-auto"></div>
                <div className="space-y-2">
                  <h1 className="text-2xl font-bold tracking-tight" style={{ color: "var(--slate-50)" }}>
                    Verifying Profile
                  </h1>
                  <p className="text-sm mt-1 animate-pulse" style={{ color: "var(--slate-400)" }}>
                    Syncing keys with server...
                  </p>
                </div>
              </div>
            )}

            {errorMessage && (
              <div className="space-y-4 text-center">
                <div className="text-red-500 text-5xl mb-2">⚠️</div>
                <h1 className="text-2xl font-bold tracking-tight" style={{ color: "var(--slate-50)" }}>
                  Authentication Error
                </h1>
                <p className="text-sm max-w-60 mx-auto leading-relaxed" style={{ color: "var(--slate-400)" }}>
                  {errorMessage}
                </p>
                <button
                  onClick={() => navigate("/login")}
                  className="mt-4 text-sm bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-2.5 rounded-xl shadow-lg transition-all active:scale-[0.98] cursor-pointer"
                >
                  Return to Login
                </button>
              </div>
            )}
          </div>

          {/* Bottom Footer Notice */}
          <div className="text-center text-[10px] font-medium tracking-wide w-full" style={{ color: "var(--slate-500)" }}>
            Secured authentication powered by GitHub OAuth protocol.
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

export default AuthCallback;
