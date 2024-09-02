import { useSession } from "next-auth/react";
import LoginButton from "../components/LoginButton";
import axios from "axios";
import { useState } from "react";

export default function Home() {
  const { data: session } = useSession();
  const [apiResponse, setApiResponse] = useState("");

  const makeAuthenticatedRequest = async () => {
    if (session && session.accessToken) {
      try {
        const response = await axios.get("http://localhost:8080/api/v1/oauth/profile", {
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
          },
        });
        setApiResponse(JSON.stringify(response.data, null, 2));
      } catch (error) {
        setApiResponse("Error making authenticated request");
      }
    }
  };

  return (
    <div>
      <h1>Next.js OAuth Example</h1>
      <LoginButton />
      {session && (
        <div>
          <button onClick={makeAuthenticatedRequest}>
            Make Authenticated Request
          </button>
          <pre>{apiResponse}</pre>
        </div>
      )}
    </div>
  );
}