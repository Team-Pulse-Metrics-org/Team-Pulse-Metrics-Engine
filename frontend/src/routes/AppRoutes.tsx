import { Routes, Route, Navigate } from "react-router-dom";
import Dashboard from "../pages/dashboard";
import Metrics from "../pages/metrics";
import Activity from "../pages/activity";
import Teams from "../pages/Teams";
import Layout from "../layouts/dashboardlayout";
import Profile from "../pages/profile";
import Login from "../pages/login";
import AuthCallback from "../pages/AuthCallback";
import Home from "../pages/Home";
import AdminPage from "../pages/Admin";
function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/" element={<Login />} />
      <Route path="/login" element={<Login />} />
      <Route path="/auth/callback" element={<AuthCallback />} />
      
      <Route element={<Layout />}>
        {/* <Route index element={<Dashboard />} /> */}
        
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="profile" element={<Profile />} />
        <Route path="metrics" element={<Metrics />} />
        <Route path="activity" element={<Activity />} />
        <Route path="teams" element={<Teams />} />
        <Route path="admin" element={<AdminPage />} />

      </Route>

      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  );
}

export default AppRoutes;
