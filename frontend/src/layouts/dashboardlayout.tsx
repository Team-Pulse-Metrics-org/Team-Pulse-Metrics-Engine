import Sidebar from "../components/sidebar";
import { Outlet } from "react-router-dom";

function Layout(){
    return(
        <div className="flex">
            <Sidebar/>

            <main className="flex-1 p-6">
                <Outlet/>
            </main>
        </div>
    );
}

export default Layout;