interface LoginResponse {
  status: string;
  token: string;
  user: {
    id: string;
    email: string;
    role: string;
  };
}

function AuthCallback() {
  return <div>Hello</div>;
}

export default AuthCallback;
