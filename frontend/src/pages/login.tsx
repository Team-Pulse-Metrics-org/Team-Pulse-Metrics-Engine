import type React from "react";
import Card from "../components/card";

function handleGitHubLogin() {
  const CLIENT_ID: string = import.meta.env.VITE_GITHUB_CLIENT_ID;
  const REDIRECT_URI: string = "http://localhost:5137/auth/callback";
  const SCOPE: string = "user,repo";

  const githubAuthUrl: string = `https://github.com/login/oauth/authorize?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&scope=${SCOPE}`;

  window.location.href = githubAuthUrl;
}

const Login: React.FC = () => {
  handleGitHubLogin();

  return <Card />;
};

export default Login;
