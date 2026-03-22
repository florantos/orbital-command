import { useCallback, useEffect, useRef, useState } from "react";

import { API_BASE_URL } from "../../lib/api";
import type { Module, ModuleError } from "./module.types";

interface ReadAllModulesResponse {
  modules: Module[];
}

function useReadModules() {
  const [modules, setModules] = useState<Module[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const controllerRef = useRef<AbortController | null>(null);

  const load = useCallback(async (signal: AbortSignal) => {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE_URL}/modules`, {
        method: "GET",
        signal,
      });
      const data = (await res.json()) as ReadAllModulesResponse | ModuleError;

      if (!res.ok) {
        setError((data as ModuleError).error);
        return;
      }

      setModules((data as ReadAllModulesResponse).modules);
    } catch (err) {
      if (err instanceof DOMException && err.name === "AbortError") return;
      setError("Network error - please try again");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    void load(controller.signal);
    return () => {
      controller.abort();
    };
  }, [load]);

  const refetch = useCallback(() => {
    controllerRef.current?.abort();
    controllerRef.current = new AbortController();
    void load(controllerRef.current.signal);
    void load(controllerRef.current.signal);
  }, [load]);

  return { modules, loading, error, refetch };
}

export { useReadModules };
