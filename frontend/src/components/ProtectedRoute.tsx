import { Navigate } from "react-router-dom";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";
interface ProtectedRouteProps {
  children: ReactNode;
  allowedRoles: string[];
}

const ProtectedRoute = ({
  children,
  allowedRoles,
}: ProtectedRouteProps) => {
  // Change this later to get the role from your login/JWT
  const role = localStorage.getItem("role");
  console.log("Current role:", role);
console.log("Allowed roles:", allowedRoles);
  const [redirect, setRedirect] = useState(false);
  useEffect(() => {
    if (role && !allowedRoles.includes(role)) {
      const timer = setTimeout(() => {
        setRedirect(true);
      }, 2000); // 2 seconds

      return () => clearTimeout(timer);
    }
  }, [role, allowedRoles]);
  if (!role) {
    return <Navigate to="/login" replace />;
  }

  if (!allowedRoles.includes(role)) {
    if (redirect){
    return <Navigate to="/dashboard" replace />;
  }
return (
      <div className="flex h-screen items-center justify-center">
        <div className="rounded-lg border p-6 text-center shadow">
          <h2 className="text-xl font-semibold text-red-600">
            Access Denied
          </h2>
          <p className="mt-2">
            You don't have permission to access this page.
          </p>
          <p className="mt-1 text-sm text-gray-500">
            Redirecting to Dashboard...
          </p>
        </div>
      </div>
    );
  }
  return <>{children} </> ;
};

export default ProtectedRoute;