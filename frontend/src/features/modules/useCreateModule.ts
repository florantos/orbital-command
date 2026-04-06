import { useState } from "react";

import { API_BASE_URL, type ApiError } from "../../lib/api";
import type { Module } from "./module.types";

function useCreateModule() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);

  const createModule = async (name: string, description: string): Promise<Module | null> => {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE_URL}/modules`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, description }),
      });

      const data = (await res.json()) as Module | ApiError;

      if (!res.ok) {
        const errData = data as ApiError;
        setError(errData);
        return null;
      }

      return data as Module;
    } catch {
      setError({ error: "Network error — please try again" });
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { createModule, loading, error };
}

export { useCreateModule };
