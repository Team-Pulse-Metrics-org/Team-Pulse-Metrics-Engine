import { useState } from "react";
import { NavLink } from "react-router-dom";

import { BsArrowLeftShort } from "react-icons/bs";
import { AiOutlineAntDesign } from "react-icons/ai";
import {LuLayoutDashboard,LuActivity,LuCircleUserRound,LuLogOut,} from "react-icons/lu";
import { BiBarChartAlt } from "react-icons/bi";

function Sidebar() {
  const [open, setOpen] = useState(true);

  const menus = [
    {
      title: "Dashboard",
      path: "/dashboard",
      icon: <LuLayoutDashboard />,
    },
    {
      title: "Metrics",
      path: "/metrics",
      icon: <BiBarChartAlt />,
    },
    {
      title: "Activity",
      path: "/activity",
      icon: <LuActivity />,
    },
  ];

  return (
    <aside
      className={`relative h-screen bg-slate-900 duration-300
      ${open ? "w-60" : "w-20"}
      flex flex-col p-5 pt-8`}
    >
      {/* Collapse Button */}
      <BsArrowLeftShort
        onClick={() => setOpen(!open)}
        className={`absolute -right-3 top-9
        text-3xl bg-white text-slate-900 rounded-full
        border border-slate-900 cursor-pointer
        duration-300
        ${!open && "rotate-180"}`}
      />

      {/* Logo */}
      <div
        className={`flex items-center ${
          open ? "justify-start" : "justify-center"
        }`}
      >
        <AiOutlineAntDesign
          className={`text-4xl bg-slate-400 rounded-md p-1 shrink-0 duration-300
          ${!open && "rotate-180"}
          ${open && "mr-2"}`}
        />

        <h1
          className={`text-2xl font-medium text-white whitespace-nowrap
          origin-left duration-300
          ${open ? "opacity-100" : "opacity-0 w-0 overflow-hidden"}`}
        >
          Team Pulse
        </h1>
      </div>

      {/* Navigation */}
      <nav className="mt-10 flex-1">
        {menus.map((menu) => (
          <NavLink
            key={menu.title}
            to={menu.path}
            className={({ isActive }) =>
              `flex items-center rounded-md p-3 mt-2 duration-200
              ${open ? "justify-start" : "justify-center"}
              ${
                isActive
                  ? "bg-slate-700 text-white"
                  : "text-slate-300 hover:bg-slate-800 hover:text-white"
              }`
            }
          >
            <span className={`text-2xl shrink-0 ${open && "mr-3"}`}>
              {menu.icon}
            </span>

            <span
              className={`whitespace-nowrap duration-300
              ${
                open
                  ? "opacity-100"
                  : "opacity-0 w-0 overflow-hidden"
              }`}
            >
              {menu.title}
            </span>
          </NavLink>
        ))}
      </nav>

      {/* Bottom Section */}
      <div className="border-t border-slate-700 pt-4">
        <NavLink
          to="/profile"
          className={({ isActive }) =>
            `flex items-center rounded-md p-3 duration-200
            ${open ? "justify-start" : "justify-center"}
            ${
              isActive
                ? "bg-slate-700 text-white"
                : "text-slate-300 hover:bg-slate-800 hover:text-white"
            }`
          }
        >
          <LuCircleUserRound
            className={`text-2xl shrink-0 ${open && "mr-3"}`}
          />

          <span
            className={`whitespace-nowrap duration-300
            ${
              open
                ? "opacity-100"
                : "opacity-0 w-0 overflow-hidden"
            }`}
          >
            Profile
          </span>
        </NavLink>

        <button
          className={`flex items-center w-full rounded-md p-3 mt-2
          duration-200
          ${open ? "justify-start" : "justify-center"}
          text-red-400 hover:bg-red-500 hover:text-white`}
        >
          <LuLogOut className={`text-2xl shrink-0 ${open && "mr-3"}`} />

          <span
            className={`whitespace-nowrap duration-300
            ${
              open
                ? "opacity-100"
                : "opacity-0 w-0 overflow-hidden"
            }`}
          >
            Sign Out
          </span>
        </button>
      </div>
    </aside>
  );
}

export default Sidebar;