import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import { toast } from "sonner";
import { api } from "@/lib/api";

export default function AuthCallback() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [isProcessing, setIsProcessing] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get("code");
      const state = searchParams.get("state");
      const error = searchParams.get("error");
      const errorDescription = searchParams.get("error_description");

      // Handle OAuth errors
      if (error) {
        setError(errorDescription || "Authentication failed");
        setIsProcessing(false);
        toast.error(errorDescription || "Authentication failed");
        
        // Redirect to signup after showing error
        setTimeout(() => {
          navigate("/auth/signup");
        }, 3000);
        return;
      }

      // Validate we have a code
      if (!code) {
        setError("No authorization code received");
        setIsProcessing(false);
        toast.error("No authorization code received");
        
        setTimeout(() => {
          navigate("/auth/signup");
        }, 3000);
        return;
      }

      try {
        // Send authorization code to backend using API client
        await api.auth.createSession(code, state || undefined);

        // Success - cookie is automatically set by browser
        toast.success("Successfully authenticated!");
        
        // Redirect to dashboard
        navigate("/dashboard");
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Authentication failed";
        setError(errorMessage);
        setIsProcessing(false);
        toast.error(errorMessage);
        
        // Redirect to signup after showing error
        setTimeout(() => {
          navigate("/auth/signup");
        }, 3000);
      }
    };

    handleCallback();
  }, [navigate, searchParams]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        {isProcessing ? (
          <>
            <div className="mb-4">
              <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
            </div>
            <h2 className="text-xl font-semibold mb-2">Completing authentication...</h2>
            <p className="text-gray-600">Please wait while we set up your session.</p>
          </>
        ) : error ? (
          <>
            <div className="mb-4">
              <svg className="mx-auto h-12 w-12 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h2 className="text-xl font-semibold mb-2 text-red-600">Authentication Failed</h2>
            <p className="text-gray-600 mb-4">{error}</p>
            <p className="text-sm text-gray-500">Redirecting to sign up...</p>
          </>
        ) : null}
      </div>
    </div>
  );
}
