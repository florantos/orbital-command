import { Dialog, ToggleGroup, VisuallyHidden } from "radix-ui";
import { useState } from "react";

import { Button } from "../../components/Button/Button";
import styles from "./CrewMemberForm.module.css";
import type { CrewMember } from "./crew.types";
import { useCreateCrewMember } from "./useCreateCrewMember";

interface CrewMemberFormProps {
  onSuccess: (module: CrewMember) => void;
  onCancel: () => void;
}

const CrewMemberForm = ({ onSuccess, onCancel }: CrewMemberFormProps) => {
  const [name, setName] = useState("");
  const [role, setRole] = useState("");
  const [qualifications, setQualifications] = useState<string[]>([]);

  const [nameError, setNameError] = useState<string | null>(null);
  const [roleError, setRoleError] = useState<string | null>(null);
  const [qualificationsError, setQualificationsError] = useState<string | null>(null);

  const { createCrewMember, error, loading } = useCreateCrewMember();

  function handleRoleValueChange(value: string) {
    if (!value) return;
    setRole(value);
  }

  function handleQualificationsValueChange(value: string[]) {
    setQualifications(value);
  }

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    let valid = true;

    if (!name.trim()) {
      setNameError("name is required");
      valid = false;
    } else {
      setNameError(null);
    }

    if (!role.trim()) {
      setRoleError("role is required");
      valid = false;
    } else {
      setRoleError(null);
    }

    if (qualifications.length === 0) {
      setQualificationsError("at least one qualification is required");
      valid = false;
    } else {
      setQualificationsError(null);
    }

    if (!valid) return;

    const crewMember = await createCrewMember(name, role, qualifications);

    if (crewMember) {
      setName("");
      setRole("");
      setQualifications([]);

      setNameError(null);
      setRoleError(null);
      setQualificationsError(null);

      onSuccess(crewMember);
    }
  }
  return (
    <form onSubmit={handleSubmit}>
      <Dialog.Title className={styles.formTitle}>Register Crew Member</Dialog.Title>
      <Dialog.Description>
        <VisuallyHidden.Root>
          A form for registering a new crew member in the orbital command system.
        </VisuallyHidden.Root>
      </Dialog.Description>
      <label htmlFor="name" className={styles.label}>
        Name
      </label>
      <input
        id="name"
        className={styles.input}
        value={name}
        onChange={(e) => {
          setName(e.target.value);
        }}
      />
      {(nameError ?? error?.fields?.name) && <p className={styles.error}>{nameError ?? error?.fields?.name}</p>}

      <label htmlFor="role" className={styles.label}>
        Role
      </label>
      <ToggleGroup.Root
        id="role"
        className={styles.toggleGroup}
        type="single"
        aria-label="Crew Member Role"
        value={role}
        onValueChange={handleRoleValueChange}
      >
        <ToggleGroup.Item value="engineer" aria-label="engineer" className={styles.toggleItem}>
          <div>Engineer</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="technician" aria-label="technician" className={styles.toggleItem}>
          <div>Technician</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="specialist" aria-label="specialist" className={styles.toggleItem}>
          <div>Specialist</div>
        </ToggleGroup.Item>
      </ToggleGroup.Root>
      {(roleError ?? error?.fields?.role) && <p className={styles.error}>{roleError ?? error?.fields?.role}</p>}

      <label htmlFor="qualifications" className={styles.label}>
        Qualifications
      </label>
      <ToggleGroup.Root
        id="qualifications"
        className={styles.toggleGroup}
        type="multiple"
        aria-label="Crew Member Qualifications"
        value={qualifications}
        onValueChange={handleQualificationsValueChange}
      >
        <ToggleGroup.Item value="docking" aria-label="docking" className={styles.toggleItem}>
          <div>docking</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="navigation" aria-label="navigation" className={styles.toggleItem}>
          <div>navigation</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="hull-monitoring" aria-label="hull-monitoring" className={styles.toggleItem}>
          <div>hull-monitoring</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="water-recycling" aria-label="water-recycling" className={styles.toggleItem}>
          <div>water-recycling</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="power-distribution" aria-label="power-distribution" className={styles.toggleItem}>
          <div>power-distribution</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="thermal-regulation" aria-label="thermal-regulation" className={styles.toggleItem}>
          <div>thermal-regulation</div>
        </ToggleGroup.Item>
        <ToggleGroup.Item value="atmosphere-recycling" aria-label="atmosphere-recycling" className={styles.toggleItem}>
          <div>atmosphere-recycling</div>
        </ToggleGroup.Item>
      </ToggleGroup.Root>
      {(qualificationsError ?? error?.fields?.qualifications) && (
        <p className={styles.error}>{qualificationsError ?? error?.fields?.qualifications}</p>
      )}

      {error && <p className={styles.errorButtons}>{error.error}</p>}
      <div className={styles.buttonRow}>
        <Button type="submit" loading={loading}>
          {loading ? "Saving..." : "Register"}
        </Button>
        <Dialog.Close asChild>
          <Button type="button" variant="secondary" onClick={onCancel} disabled={loading}>
            Cancel
          </Button>
        </Dialog.Close>
      </div>
    </form>
  );
};

export { CrewMemberForm };
