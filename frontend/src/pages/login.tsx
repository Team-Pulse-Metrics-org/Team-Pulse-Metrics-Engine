import type React from "react";
import Card from "../components/card";

import { FaGithub } from "react-icons/fa";

function handleGitHubLogin() {
  const CLIENT_ID: string = import.meta.env.VITE_GITHUB_CLIENT_ID;
<<<<<<< Updated upstream
=======
  console.log(CLIENT_ID)
>>>>>>> Stashed changes
  const REDIRECT_URI: string = "http://localhost:5173/auth/callback";
  const SCOPE: string = "user,public_repo";

  const githubAuthUrl: string = `https://github.com/login/oauth/authorize?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&scope=${SCOPE}`;

  window.location.href = githubAuthUrl;
}

const Login: React.FC = () => {
  return (
    <div className="text-white bg-slate-950 min-h-screen w-full flex flex-col gap-14 items-center justify-center">
      <h1 className="pb-10 text-transparent font-bold text-6xl bg-linear-to-r from-purple-500 via-pink-500 to-yellow-500 bg-clip-text font-inter">
        Team Pulse
      </h1>
      <div className=" relative max-w-6xl w-full mx-auto h-[550px] rounded-3xl overflow-hidden shadow-2xl flex items-center justify-end">
        <div>
          <img
            src="rocket.jpg"
            className="rounded-3xl absolute inset-0 w-full h-full object-cover z-0"
          />
        </div>

        <Card className="relative z-30 w-1/3 h-full px-8 py-12 flex flex-col gap-30 items-center justify-center border-none rounded-3xl backdrop-blur-md bg-transparent ">
          <h1 className="font-bold text-5xl font-inter p-10">Sign in</h1>
          <div className="flex flex-col gap-5">
            <p className="">Sign in with your Github Account</p>
            <button
              className="flex items-center justify-center gap-3 w-50 mx-auto bg-indigo-600 p-3 hover:bg-indigo-800 transition-all rounded-2xl group shadow-lg font-medium"
              onClick={handleGitHubLogin}
              title="Sign in with GitHub"
            >
              <FaGithub
                size={35}
                className="transition-transform group-hover:scale-110"
              />
            </button>
          </div>
        </Card>
      </div>
      <p className="text-center text-xs text-slate-500 font-medium tracking-wide">
        © 2026 Team Pulse · Engineering Analytics Platform
      </p>
    </div>
  );
};

export default Login;
