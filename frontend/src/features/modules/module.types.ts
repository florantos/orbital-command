export type HealthState = "operational" | "degraded" | "critical" | "unresponsive" | "offline";

export interface Module {
  id: string;
  name: string;
  description: string;
  healthState: HealthState;
}

export interface ModuleError {
  error: string;
}
