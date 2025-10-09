import React from "react";

interface HeaderProps {
  totalCount: number;
  completedCount: number;
}

const Header: React.FC<HeaderProps> = ({ totalCount, completedCount }) => {
  const progress = totalCount > 0 ? (completedCount / totalCount) * 100 : 0;

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-600 text-white p-8 relative overflow-hidden">
      {/* Background animations */}
      <div className="absolute top-0 left-0 w-full h-full pointer-events-none">
        <div className="absolute top-4 left-8 w-32 h-32 bg-white/10 rounded-full blur-xl animate-pulse"></div>
        <div className="absolute bottom-4 right-12 w-24 h-24 bg-purple-300/20 rounded-full blur-lg animate-bounce"></div>
        <div className="absolute top-1/2 right-16 w-16 h-16 bg-pink-300/15 rounded-full blur-md animate-ping"></div>
      </div>

      {/* Content container */}
      <div className="relative z-10 max-w-2xl w-full bg-white/10 backdrop-blur-md rounded-3xl shadow-2xl p-8 space-y-6 text-center">
        {/* Header section */}
        <div>
          <h1 className="text-5xl font-extrabold text-white tracking-wide drop-shadow-md">
            <span className="bg-gradient-to-r from-yellow-300 via-pink-300 to-indigo-400 bg-clip-text text-transparent">
              ✨ My Todo List
            </span>
          </h1>

          <p className="text-lg text-white/80 mt-2">
            Stay organized and get things done
          </p>
        </div>

        {/* Stats section */}
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="bg-white/20 backdrop-blur-sm rounded-2xl p-4 border border-white/30 shadow-xl w-full sm:w-1/2">
            <div className="text-center">
              <div className="text-3xl font-bold">
                {completedCount}/{totalCount}
              </div>
              <div className="text-xs text-purple-100 uppercase tracking-widest">
                Tasks Done
              </div>
            </div>
          </div>

          <div className="bg-white/15 backdrop-blur-md rounded-2xl p-6 border border-white/20 shadow-xl w-full sm:w-1/2">
            <div className="flex items-center justify-between ">
              <div className="text-left">
                <div className="text-white font-semibold text-lg">
                  {completedCount} of {totalCount} completed
                </div>
                <div className="text-purple-100 text-sm">
                  {totalCount === 0
                    ? "No tasks yet"
                    : completedCount === totalCount
                    ? "All done! "
                    : `${totalCount - completedCount} remaining`}
                </div>
              </div>
              <div className="text-right">
                <div className="text-3xl font-black text-white mb-1">
                  {Math.round(progress)}%
                </div>
                <div className="text-xs text-purple-100 uppercase tracking-wider">
                  Complete
                </div>
              </div>
            </div>

            {/* Progress bar */}
            <div className="relative mt-4">
              <div className="w-full h-3 bg-white/20 rounded-full overflow-hidden shadow-inner">
                <div
                  className="h-full bg-gradient-to-r from-emerald-400 via-cyan-400 to-blue-400 rounded-full transition-all duration-700 ease-out shadow-lg relative"
                  style={{ width: `${progress}%` }}
                >
                  <div className="absolute inset-0 bg-white/30 rounded-full animate-pulse"></div>
                </div>
              </div>

              {/* Progress dot */}
              {progress > 0 && (
                <div
                  className="absolute top-1/2 transform -translate-y-1/2 w-4 h-4 bg-white rounded-full shadow-lg border-2 border-purple-300 transition-all duration-700 ease-out"
                  style={{
                    left: `${Math.min(progress, 100)}%`,
                    transform: "translate(-50%, -50%)",
                  }}
                >
                  <div className="absolute inset-0 bg-gradient-to-r from-emerald-400 to-cyan-400 rounded-full animate-ping opacity-75"></div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Motivation message */}
        <p className="text-purple-100 text-sm font-medium">
          {progress === 100
            ? " Amazing work! You've completed everything!"
            : progress >= 75
            ? "🔥 You're almost there! Keep going!"
            : progress >= 50
            ? "💪 Great progress! You're halfway done!"
            : progress > 0
            ? " Good start! Keep building momentum!"
            : "📝 Ready to tackle your tasks?"}
        </p>
      </div>
    </div>
  );
};

export default Header;
