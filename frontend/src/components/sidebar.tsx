import { useState } from "react";
import { BsArrowLeftShort } from "react-icons/bs";
import { AiOutlineAntDesign } from "react-icons/ai";



function Sidebar(){
    const [Open,setOpen]=useState(true);
    return(
        <div className='flex flex-col h-full'>
            <div className={`bg-slate-800 h-screen p-5 ${Open ? "w-60" : "w-20"} pt-8 duration-300 relative`}>
                <BsArrowLeftShort className={`bg-white text-slate-800 text-3xl 
                    rounded-full absolute -right-3 top-9 border
                    border-slate-800 cursor-pointer ${!Open && "rotate-180"}`}
                    onClick={()=>setOpen(!Open)}
                />
                <div className="inline-flex">
                    <AiOutlineAntDesign className={`bg-slate-300 rounded-md shrink-0 duration-500 
                        ${!Open && "rotate-180"}
                        text-4xl block float-left mr-2`}
                    />
                    <h1 className={`${!Open && "scale-0"} duration-300 text-2xl text-white origin-left font-medium`}>Team Pulse</h1>
                </div>
            </div>

        </div>
    );
}

export default Sidebar;