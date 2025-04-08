// LoadingFallback.tsx
const LoadingFallback = () => {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-80 z-50">
      <span className="text-6xl font-bold animate-pulse">
        Loading...
      </span>
    </div>
  );
};

export default LoadingFallback;
