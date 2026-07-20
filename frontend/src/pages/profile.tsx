import { useEffect, useState } from "react";
import Card from "../components/card";
import { FaGithub, FaEnvelope, FaUser, FaUsers, FaInfoCircle, FaExclamationCircle } from "react-icons/fa";

interface ProfileData {
  github_id: number;
  username: string;
  name: string;
  email: string;
  avatar_url: string;
  followers: number;
  following: number;
}

export default function Profile() {
  const [profile, setProfile] = useState<ProfileData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchProfile = () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem("app_token");

    if (!token) {
      window.location.href = "/login";
      return;
    }

    fetch("http://localhost:8080/api/v1/profile", {
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
          throw new Error("Failed to fetch profile details");
        }
        return res.json();
      })
      .then((data) => {
        setProfile(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching profile:", err);
        setError(err.message || "Failed to load profile");
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchProfile();
  }, []);

  if (loading) {
    return (
      <div className="bg-slate-950 min-h-screen p-8 text-slate-100 flex flex-col justify-center items-center">
        <div className="relative w-16 h-16">
          <div className="absolute inset-0 rounded-full border-4 border-slate-800"></div>
          <div className="absolute inset-0 rounded-full border-4 border-indigo-500 border-t-transparent animate-spin"></div>
        </div>
        <p className="mt-4 text-slate-400 font-medium animate-pulse">Loading profile data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-slate-950 min-h-screen p-8 text-slate-100 flex flex-col justify-center items-center">
        <div className="max-w-md w-full bg-slate-900 border border-red-500/20 rounded-2xl p-6 text-center shadow-xl">
          <div className="inline-flex p-3 bg-red-500/10 rounded-full text-red-400 mb-4">
            <FaExclamationCircle className="w-8 h-8" />
          </div>
          <h2 className="text-xl font-bold text-red-200 mb-2">Error Loading Profile</h2>
          <p className="text-slate-400 text-sm mb-6">{error}</p>
          <button
            onClick={fetchProfile}
            className="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-500 active:bg-indigo-700 text-white font-medium rounded-lg transition duration-200"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-slate-950 min-h-screen p-8 text-slate-100">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-extrabold tracking-tight bg-gradient-to-r from-slate-50 via-slate-100 to-slate-400 bg-clip-text text-transparent">
          Profile Settings
        </h1>
        <p className="text-slate-400 text-sm mt-1">Manage and view your linked GitHub account details</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Main Card (Avatar, Username, Stats) */}
        <div className="lg:col-span-1 flex flex-col gap-6">
          <Card className="p-6 flex flex-col items-center text-center relative overflow-hidden group">
            {/* Ambient Background Glow */}
            <div className="absolute -top-24 -left-24 w-48 h-48 bg-indigo-500/10 rounded-full blur-3xl group-hover:bg-indigo-500/15 transition-all duration-500"></div>
            
            {/* Avatar with Glow on Hover */}
            <div className="relative mb-4 mt-2">
              <div className="absolute inset-0 rounded-full bg-indigo-500/25 blur-sm opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
              <img
                src={profile?.avatar_url}
                alt={profile?.name || "Avatar"}
                className="w-28 h-28 rounded-full border-2 border-slate-800 object-cover relative z-10"
              />
            </div>

            {/* User Info */}
            <h2 className="text-xl font-bold text-slate-100">{profile?.name}</h2>
            <p className="text-slate-400 text-sm font-medium">@{profile?.username}</p>

            {/* External Link */}
            <a
              href={`https://github.com/${profile?.username}`}
              target="_blank"
              rel="noreferrer"
              className="mt-4 flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 active:bg-slate-700/80 rounded-xl text-sm font-medium border border-slate-700/50 hover:border-slate-600 transition duration-200 cursor-pointer"
            >
              <FaGithub className="w-4 h-4" />
              View GitHub Profile
            </a>

            {/* Stats Separator */}
            <div className="w-full border-t border-slate-800/80 my-5"></div>

            {/* Stats Row */}
            <div className="flex justify-around w-full">
              <div className="text-center">
                <span className="block text-2xl font-extrabold text-slate-100">{profile?.followers}</span>
                <span className="text-xs text-slate-500 font-semibold uppercase tracking-wider flex items-center gap-1 mt-1 justify-center">
                  <FaUsers className="w-3.5 h-3.5" />
                  Followers
                </span>
              </div>
              <div className="w-px bg-slate-800/80 self-stretch"></div>
              <div className="text-center">
                <span className="block text-2xl font-extrabold text-slate-100">{profile?.following}</span>
                <span className="text-xs text-slate-500 font-semibold uppercase tracking-wider flex items-center gap-1 mt-1 justify-center">
                  <FaUsers className="w-3.5 h-3.5" />
                  Following
                </span>
              </div>
            </div>
          </Card>
        </div>

        {/* Right Column - Profile Details */}
        <div className="lg:col-span-2 flex flex-col gap-6">
          <Card className="p-6">
            <h3 className="text-lg font-bold text-slate-200 mb-6 flex items-center gap-2 border-b border-slate-800 pb-3">
              <FaUser className="w-5 h-5 text-indigo-400" />
              Linked Account Information
            </h3>

            <div className="space-y-5">
              {/* Username Info */}
              <div className="flex flex-col sm:flex-row sm:items-center justify-between border-b border-slate-800/40 pb-4">
                <span className="text-slate-400 text-sm font-medium">GitHub ID</span>
                <span className="text-slate-200 font-mono text-sm mt-1 sm:mt-0 bg-slate-950 px-3 py-1.5 rounded-lg border border-slate-800/60">
                  {profile?.github_id}
                </span>
              </div>

              {/* Full Name */}
              <div className="flex flex-col sm:flex-row sm:items-center justify-between border-b border-slate-800/40 pb-4">
                <span className="text-slate-400 text-sm font-medium">Full Name</span>
                <span className="text-slate-200 text-sm mt-1 sm:mt-0 font-semibold">
                  {profile?.name || "N/A"}
                </span>
              </div>

              {/* Email */}
              <div className="flex flex-col sm:flex-row sm:items-center justify-between border-b border-slate-800/40 pb-4">
                <span className="text-slate-400 text-sm font-medium flex items-center gap-1.5">
                  <FaEnvelope className="w-4 h-4 text-indigo-400" />
                  Primary Email
                </span>
                <span className="text-slate-200 text-sm mt-1 sm:mt-0 font-medium">
                  {profile?.email || "No email available"}
                </span>
              </div>

              {/* Status */}
              <div className="flex flex-col sm:flex-row sm:items-center justify-between border-b border-slate-800/40 pb-4">
                <span className="text-slate-400 text-sm font-medium">Authentication Source</span>
                <span className="inline-flex items-center gap-1.5 text-xs font-semibold bg-emerald-500/10 text-emerald-400 px-3 py-1.5 rounded-full border border-emerald-500/20 mt-1 sm:mt-0">
                  <FaGithub className="w-3.5 h-3.5" />
                  GitHub OAuth
                </span>
              </div>
            </div>
          </Card>

          {/* Quick Notice Card */}
          <div className="bg-indigo-950/20 border border-indigo-500/10 rounded-2xl p-5 flex items-start gap-4 shadow-sm">
            <div className="p-2 bg-indigo-500/10 rounded-xl text-indigo-400 mt-0.5">
              <FaInfoCircle className="w-5 h-5" />
            </div>
            <div>
              <h4 className="font-bold text-indigo-200 text-sm mb-1">Looking for profile changes?</h4>
              <p className="text-slate-400 text-xs leading-relaxed">
                This account is currently linked using GitHub OAuth. To update your name, email, avatar, or basic details, please update them on your GitHub account settings. Changes will be synced automatically the next time you log in.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}