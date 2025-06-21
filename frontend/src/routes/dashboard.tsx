import { useNavigate } from "react-router";

import { useAuthStore } from "@/stores/auth-store";

import { api } from "@/lib/api";

import { ProtectedRoute } from "@/components/auth/protected-route";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function Dashboard() {
  const { user } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await api.auth.logout();
      // Clear the auth store and redirect
      useAuthStore.getState().clearUser();
      navigate("/auth/signup");
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  return (
    <ProtectedRoute>
      <div className="container mx-auto p-6 max-w-4xl">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <Button onClick={handleLogout} variant="outline">
            Logout
          </Button>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Welcome, {user?.name || "User"}!</CardTitle>
            <CardDescription>You are successfully authenticated via {user?.provider}.</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p>
                <strong>Email:</strong> {user?.email}
              </p>
              <p>
                <strong>ID:</strong> {user?.id}
              </p>
              <p>
                <strong>Provider:</strong> {user?.provider}
              </p>
              {user?.picture && (
                <div className="mt-4">
                  <img src={user.picture} alt="Profile" className="w-20 h-20 rounded-full" />
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </ProtectedRoute>
  );
}
