import { Routes, Route, Navigate } from "react-router-dom";
import Dashboard from "../pages/dashboard";
import Metrics from "../pages/metrics";
import Activity from "../pages/activity";
import Layout from "../layouts/dashboardlayout";
import Profile from "../pages/profile";
import Login from "../pages/login";

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      {/* <Route path="/auth/callback" element={<AuthCallbackPage />} /> */}

      <Route element={<Layout />}>
        <Route index element={<Dashboard />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="profile" element={<Profile />} />
        <Route path="metrics" element={<Metrics />} />
        <Route path="activity" element={<Activity />} />
      </Route>

      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  );
}

export default AppRoutes;
