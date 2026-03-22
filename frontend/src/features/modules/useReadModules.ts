import { useState } from "react";

import { API_BASE_URL } from "../../lib/api";
import type { Module, ModuleError } from "./module.types";

interface ReadAllModulesResponse {
  modules: Module[];
}
function useReadModules() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const readModules = async (): Promise<ReadAllModulesResponse | null> => {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE_URL}/modules`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = (await res.json()) as ReadAllModulesResponse | ModuleError;

      if (!res.ok) {
        const errData = data as ModuleError;
        setError(errData.error);
        return null;
      }
      return data as ReadAllModulesResponse;
    } catch {
      setError("Network error - please try again");
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { readModules, loading, error };
}

export { useReadModules };
