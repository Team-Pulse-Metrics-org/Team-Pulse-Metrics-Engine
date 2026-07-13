import { Link } from "react-router-dom";
import { useState } from "react";
export default function Home() {
  const [showFeatures, setShowFeatures] = useState(false);
  const [showAbout, setShowAbout] = useState(false);

  return (
    //navigation
    <div className="min-h-screen bg-slate-950 text-white">
      {/* Navbar */}
      <nav className="w-full flex justify-between items-center px-16 py-6 border-b border-slate-800">
        {/* Logo */}
        <div className="text-3xl font-bold">
          <span className="text-cyan-400">Team</span>
          <span className="text-purple-500">Pulse</span>
        </div>

        {/* Navigation Links */}
        <div className="flex items-center gap-8 text-lg">
          <button
  onClick={() => {
    setShowFeatures(true);
    setShowAbout(false);
  }}
  className="hover:text-cyan-400 transition duration-200"
>

  Features
</button>

<button
  onClick={() => {
    setShowAbout(true);
    setShowFeatures(false);
  }}
  className="hover:text-cyan-400 transition duration-200"
>
  About
</button>
<Link
  to="/login"
  className="bg-cyan-500 hover:bg-cyan-600 px-6 py-3 rounded-xl font-semibold transition"
>
  Login
</Link>
        </div>
      </nav>
       {/*hero section lil details*/} 
<section className="flex flex-col items-center justify-center text-center mt-32 px-6">
  <h1 className="text-6xl font-extrabold leading-tight max-w-5xl">
    Track Engineering Productivity
    <span className="text-cyan-400"> in Real Time</span>
  </h1>

  <p className="mt-8 text-xl text-slate-400 max-w-3xl">
    Monitor commits, pull requests, completed tasks and
    engineering blockers from a single intelligent platform.
  </p>

  <div className="mt-10">
    <Link
  to="/login"
  className="bg-cyan-500 hover:bg-cyan-600 px-10 py-4 rounded-2xl font-semibold transition"
  >
  Get Started
  </Link>

  </div>
</section>
{/*features*/}
{showFeatures && (
<section className="mt-40 px-10 py-20">
 
  <h2 className="text-5xl font-bold text-center mb-16">
    Features
  </h2>

  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
    <div className="bg-slate-900/70 backdrop-blur-lg p-8
rounded-2xl
border border-slate-700
shadow-2xl
hover:-translate-y-3
hover:shadow-cyan-500/20
hover:border-cyan-400
transition-all duration-300">
      <div className="text-5xl mb-4">💻</div>
      <h3 className="text-2xl font-semibold mb-3">
        GitHub Tracking
      </h3>
      <p className="text-slate-400">
        Monitor commits and developer activity in real time.
      </p>
    </div>

    <div className="bg-slate-900/70 backdrop-blur-lg p-8
rounded-2xl
border border-slate-700
shadow-2xl
hover:-translate-y-3
hover:shadow-cyan-500/20
hover:border-cyan-400
transition-all duration-300">
      <div className="text-5xl mb-4">🔀</div>
      <h3 className="text-2xl font-semibold mb-3">
        Pull Requests
      </h3>
      <p className="text-slate-400">
        Track code reviews and merged pull requests.
      </p>
    </div>

    <div className="bg-slate-900/70 backdrop-blur-lg p-8
rounded-2xl
border border-slate-700
shadow-2xl
hover:-translate-y-3
hover:shadow-cyan-500/20
hover:border-cyan-400
transition-all duration-300">
      <div className="text-5xl mb-4">✅</div>
      <h3 className="text-2xl font-semibold mb-3">
        Task Metrics
      </h3>
      <p className="text-slate-400">
        Measure completed tasks and delivery velocity.
      </p>
    </div>
 
    <div className="
bg-slate-900/70 backdrop-blur-lg p-8
rounded-2xl
border border-slate-700
shadow-2xl
hover:-translate-y-3
hover:shadow-cyan-500/20
hover:border-cyan-400
transition-all duration-300
">
      <div className="text-5xl mb-4">🚨</div>
      <h3 className="text-2xl font-semibold mb-3">
        Blocker Detection
      </h3>
      <p className="text-slate-400">
        Identify bottlenecks and active engineering blockers.
      </p>
    </div>
  </div>
</section>
)}
{showAbout && (
  <section className="py-24 px-10 text-center">
  <h2 className="text-5xl font-bold mb-8">
    About Team Pulse
  </h2>

  <div className="max-w-4xl mx-auto bg-slate-900/70 backdrop-blur-lg border border-slate-800 rounded-3xl p-10 shadow-xl shadow-cyan-500/10">
    <p className="text-xl text-slate-300 leading-9">
      Team Pulse is an engineering analytics platform designed to
      provide real-time insights into team productivity and workflow
      health. By automatically tracking GitHub activities such as
      commits, pull requests, completed tasks, and blockers, Team
      Pulse helps engineering teams monitor progress, identify
      bottlenecks, and make data-driven decisions.
    </p>

    <p className="text-xl text-slate-300 leading-9 mt-6">
      Built using modern technologies and automated data collection,
      Team Pulse eliminates manual reporting and enables teams to
      focus more on development and collaboration.
    </p>
  </div>
 
</section>
)}
    </div>
     
  );
}