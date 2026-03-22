export interface Module {
  id: string;
  name: string;
  description: string;
  healthState: string;
}

export interface ModuleError {
  error: string;
}
