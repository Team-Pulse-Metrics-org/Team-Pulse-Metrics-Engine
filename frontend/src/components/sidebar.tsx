import { useState } from "react";
import { BsArrowLeftShort } from "react-icons/bs";
import { AiOutlineAntDesign } from "react-icons/ai";
import { LuLayoutDashboard,LuActivity,LuCircleUserRound   } from "react-icons/lu";
import { NavLink } from "react-router-dom";
import { BiBarChartAlt } from "react-icons/bi";
import { LuLogOut } from "react-icons/lu";



function Sidebar(){
    const [Open,setOpen]=useState(true);
    const menus = [
    {
        title: "Dashboard",
        path: "/dashboard",
        icon: <LuLayoutDashboard/>
    },
    {
        title: "Metrics",
        path: "/metrics",
        icon:<BiBarChartAlt/>
    },
    {
        title: "Activity",
        path: "/activity",
        icon:<LuActivity/>
    },
];

    return(
        <div className='flex flex-col h-full'>
            <div className={`bg-slate-800 h-screen flex flex-col p-5 ${Open ? "w-60" : "w-20"} pt-8 duration-300 relative`}>
                <BsArrowLeftShort className={`bg-white text-slate-800 text-3xl 
                    rounded-full absolute -right-3 top-9 border
                    border-slate-800 cursor-pointer ${!Open && "rotate-180"}`}
                    onClick={()=>setOpen(!Open)}
                />
                <div className="inline-flex">
                    <AiOutlineAntDesign className={`bg-slate-300 rounded-md shrink-0 duration-500 
                        ${!Open && "rotate-180"} ${Open ? "justify-start" : "justify-center"}
                        text-4xl block float-left mr-2`}
                    />
                    <h1 className={`${!Open && "scale-0"} duration-300 text-2xl text-white origin-left font-medium`}>Team Pulse</h1>
                </div>
                <div className="flex flex-col flex-1">
                    <nav className="mt-10">
                        {menus.map((menu)=>(
                            <NavLink key={menu.title} to={menu.path}
                                className={({isActive})=>
                                    `flex items-center ${Open ? "justify-start" : "justify-center"} p-3 mt-2 rounded-md text-white hover:bg-slate-700 ${
                                        isActive ? "bg-slate-500" : ""
                                    }`
                                }
                            >
                                <span className={`text-2xl cursor-pointer shrink-0 ${Open && "mr-1.5"}`}>{menu.icon}</span>
                                <span className={`duration-300 ${!Open && "scale-0 hidden"}`}>
                                    {menu.title}
                                </span>
                            </NavLink>
                        ))}
                    </nav>
                    <div className="mt-auto">
                        <NavLink to="/profile" 
                            className={({isActive})=>`flex items-center p-3 text-white rounded-md hover:bg-slate-700
                                        ${isActive ? "bg-slate-500" : ""}
                                        ${Open ? "justify-start" : "justify-center"}
                                `}
                        >
                            <LuCircleUserRound className={`text-2xl cursor-pointer shrink-0 ${Open && "mr-1.5"}`}/>
                            {Open && <span className="text-white ml-2">Profile</span>}
                        </NavLink>
                        <button className={`flex text-red-400 items-center rounded-md p-3 mt-2
                                            w-full hover:bg-red-500 hover:text-white
                                            ${Open ? "justify-start" : "justify-center"}
                            `}>
                            <LuLogOut className={`text-2xl cursor-pointer shrink-0 ${Open && "mr-1.5"}`}/>
                            {Open && <span className="text-white ml-2">Sign Out</span>}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Sidebar;