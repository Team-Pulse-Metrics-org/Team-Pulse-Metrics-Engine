import {Routes, Route } from 'react-router-dom';
import Dashboard from '../pages/dashboard';
import Metrics from '../pages/metrics';
import Activity from '../pages/activity';
import Layout from "../layouts/dashboardlayout";

function AppRoutes(){
    return(
        <Routes>
            <Route element={<Layout/>}>
                <Route index element={<Dashboard/>}/>
                <Route path='dashboard' element={<Dashboard/>}/>
                <Route path='metrics' element={<Metrics/>}/>
                <Route path='activity' element={<Activity/>}/>
            </Route>
        </Routes>
    )
}

export default AppRoutes;