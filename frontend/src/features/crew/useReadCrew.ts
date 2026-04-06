import { useCallback, useEffect, useRef, useState } from "react";

import { API_BASE_URL, type ApiError } from "../../lib/api";
import type { CrewMember } from "./crew.types";

interface ReadAllCrewRespose {
  crew: CrewMember[];
}

function useReadCrew() {
  const [crew, setCrew] = useState<CrewMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const controllerRef = useRef<AbortController | null>(null);

  const load = useCallback(async (signal: AbortSignal) => {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${API_BASE_URL}/crew`, {
        method: "GET",
        signal,
      });
      const data = (await res.json()) as ReadAllCrewRespose | ApiError;

      if (!res.ok) {
        setError((data as ApiError).error);
        return;
      }

      setCrew((data as ReadAllCrewRespose).crew);
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
  }, [load]);

  return { crew, loading, error, refetch };
}

export { useReadCrew };
