import Sidebar from "../components/sidebar";
import { Outlet } from "react-router-dom";

function Layout(){
    return(
        <div className="flex h-screen">
            <Sidebar/>

            <main className="flex-1 overflow-y-auto bg-slate-950">
                <Outlet/>
            </main>
        </div>
    );
}

export default Layout;