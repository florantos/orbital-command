import { useState } from "react";

import { API_BASE_URL } from "../../lib/api";
import type { Module } from "./module.types";

interface CreateModuleError {
  error: string;
}

function useCreateModule() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

      const data = (await res.json()) as Module | CreateModuleError;

      if (!res.ok) {
        const errData = data as CreateModuleError;
        setError(errData.error);
        return null;
      }

      return data as Module;
    } catch {
      setError("Network error — please try again");
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { createModule, loading, error };
}

export { useCreateModule };
