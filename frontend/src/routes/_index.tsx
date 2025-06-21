import { Navigate } from "react-router";

import { useUser } from "@/hooks/use-user";

import { useAuthStore } from "@/stores/auth-store";

export default function Index() {
  const { isAuthenticated } = useAuthStore();
  const { isLoading } = useUser();

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  // Redirect authenticated users to dashboard
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  // Redirect non-authenticated users to signup
  return <Navigate to="/auth/signup" replace />;
}
