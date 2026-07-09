import type React from "react";
import { useEffect, useRef, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import Card from "../components/card";

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

  useEffect(() => {
    const code = searchParams.get("code");

    if (!code) {
      setIsLoading(false);
      setErrorMessage("github temporary token invalid");
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
    <div className="text-white bg-slate-950 min-h-screen w-full flex flex-col gap-14 items-center justify-center">
      <h1 className="pb-10 text-transparent font-bold text-6xl bg-linear-to-r from-purple-500 via-pink-500 to-yellow-500 bg-clip-text font-inter animate-pulse">
        Team Pulse
      </h1>

      {/* Outer Box: Filled with the image, card pinned to the right */}
      <div className=" relative max-w-6xl w-full mx-auto h-[550px] rounded-3xl overflow-hidden shadow-2xl flex items-center justify-end">
        <img
          src="/rocket.jpg"
          alt="Rocket Nasa background"
          className="absolute inset-0 w-full h-full object-cover z-0"
        />

        <Card className="relative z-30 w-1/3 h-full px-8 py-12 flex flex-col gap-30 items-center justify-center border-none rounded-3xl backdrop-blur-md bg-transparent">
          {isLoading && !errorMessage && (
            <div className="space-y-4 text-center">
              <div className="h-10 w-10 animate-spin rounded-full border-4 border-slate-800 border-t-indigo-500 mx-auto"></div>
              <div>
                <h1 className="text-xl font-bold tracking-wide mt-5">
                  Verifying Profile
                </h1>
                <p className="text-xs text-slate-400 mt-1 animate-pulse">
                  Syncing Keys with Go server...
                </p>
              </div>
            </div>
          )}

          {errorMessage && (
            <div className="space-y-4 text-center">
              <div className="text-red-500 text-3xl font-bold">Error</div>
              <div className="text-lg font-bold text-red-500">
                Authentication Error
              </div>
              <p className="text-xs text-slate-400 max-w-60 mx-auto">
                {errorMessage}
              </p>
              <button
                onClick={() => navigate("/login")}
                className="mt-2 text-xs bg-slate-800 border border-slate-700 hover:bg-slate-700 px-4 py-2 rounded-lg font-medium transition-all"
              >
                Return to Login
              </button>
            </div>
          )}
        </Card>
      </div>

      <p className="text-center text-xs text-slate-500 font-medium tracking-wide">
        © 2026 Team Pulse · Engineering Analytics Platform
      </p>
    </div>
  );
};

export default AuthCallback;
