import type React from "react";

type ClassProps = {
  className?: string;
  children?: React.ReactNode;
};

function Card({ className = "", children }: ClassProps) {
  return (
    <div className={`border bg-slate-900 rounded-2xl ${className}`}>
      {children}
    </div>
  );
}
export default Card;
