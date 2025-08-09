import { FaGoogle } from "react-icons/fa";
import { FaMicrosoft } from "react-icons/fa";

import { Button } from "@/components/ui/button";

export default function Signup() {
  const onMicrosoftSignUp = () => {
    // Initiate OAuth2 flow with Hydra
    window.location.href =
      `${import.meta.env.VITE_HYDRA_PUBLIC_URL}/oauth2/auth?` +
      new URLSearchParams({
        client_id: "microsoft",
        response_type: "code",
        redirect_uri: window.location.origin,
        scope: "offline_access offline openid",
        state: Math.random().toString(36).substring(2, 15),
      });
  };

  return (
    <div className="flex flex-1 flex-row items-center justify-center mx-auto h-screen gap-3 p-7 max-w-7xl">
      <div className="flex flex-1 flex-col w-full h-full justify-center items-center gap-6">
        <div className="flex flex-col justify-start gap-4 w-[380px]">
          <img src="/images/conformitea.svg" alt="Conformitea Logo" className="w-10 h-10" />
          <div className="gap-2 flex flex-col items-start">
            <span className="header-md">Create an account</span>
            <span className="text-sm text-muted-foreground">
              Let's get you started. Choose one of the authentication methods below.
            </span>
          </div>
        </div>
        <div className="flex flex-col gap-4 w-[380px]">
          <Button className="w-full" variant="outline" size="lg">
            <FaGoogle className="mr-1" />
            Sign up with Google
          </Button>
          <Button className="w-full" variant="outline" size="lg" onClick={onMicrosoftSignUp}>
            <FaMicrosoft className="mr-1" />
            Sign up with Microsoft
          </Button>
        </div>
        <div></div>
      </div>
      <div className="flex flex-1 w-full h-full bg-[url(/images/auth/background.jpg)] bg-cover bg-center opacity-40 rounded-lg"></div>
    </div>
  );
}
