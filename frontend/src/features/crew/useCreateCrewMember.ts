import { useState } from "react";

import { API_BASE_URL, type ApiError } from "../../lib/api";
import type { CrewMember } from "./crew.types";

function useCreateCrewMember() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);

  const createCrewMember = async (name: string, role: string, qualifications: string[]): Promise<CrewMember | null> => {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE_URL}/crew`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, role, qualifications }),
      });

      const data = (await res.json()) as CrewMember | ApiError;

      if (!res.ok) {
        const errData = data as ApiError;
        setError(errData);
        return null;
      }

      return data as CrewMember;
    } catch {
      setError({ error: "Network error — please try again" });
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { createCrewMember, loading, error };
}

export { useCreateCrewMember };
