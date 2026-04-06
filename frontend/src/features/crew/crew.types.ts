type Role = "engineer" | "tactician" | "specialist";

type Capability =
  | "docking"
  | "navigation"
  | "hull-monitoring"
  | "water-recycling"
  | "power-generation"
  | "power-distribution"
  | "thermal-regulation"
  | "atmosphere-recycling";

export interface CrewMember {
  id: string;
  name: string;
  role: Role;
  qualifications: Capability[];
}
