import React from "react";

interface HeaderProps {
  totalCount: number;
  completedCount: number;
}

const Header: React.FC<HeaderProps> = ({ totalCount, completedCount }) => {
  const progress = totalCount > 0 ? (completedCount / totalCount) * 100 : 0;

  return (
    <header className="rounded-2xl bg-gradient-to-r from-slate-900 via-blue-900 to-cyan-800 p-6 sm:p-8 text-white shadow-xl border border-white/10">
      <div className="flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-sm uppercase tracking-[0.2em] text-cyan-200">Task Dashboard</p>
          <h1 className="mt-2 text-3xl sm:text-4xl font-bold">My Todo List</h1>
          <p className="mt-2 text-cyan-100/90">Stay focused and ship one task at a time.</p>
        </div>

        <div className="rounded-xl bg-white/10 px-4 py-3 backdrop-blur-md border border-white/15">
          <p className="text-xs uppercase tracking-wide text-cyan-100/80">Progress</p>
          <p className="mt-1 text-2xl font-semibold">
            {completedCount}/{totalCount}
          </p>
        </div>
      </div>

      <div className="mt-6">
        <div className="flex items-center justify-between text-sm text-cyan-100/90">
          <span>{completedCount} completed</span>
          <span>{Math.round(progress)}%</span>
        </div>
        <div className="mt-2 h-2.5 rounded-full bg-black/20 overflow-hidden">
          <div
            className="h-full rounded-full bg-gradient-to-r from-emerald-300 via-cyan-300 to-sky-200 transition-all duration-500"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>
    </header>
  );
};

export default Header;
