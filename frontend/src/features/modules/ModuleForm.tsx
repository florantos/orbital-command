import { Dialog, VisuallyHidden } from "radix-ui";
import { useState } from "react";

import { Button } from "../../components/Button/Button";
import styles from "./ModuleForm.module.css";
import type { Module } from "./module.types";
import { useCreateModule } from "./useCreateModule";

interface ModuleFormProps {
  onSuccess: (module: Module) => void;
  onCancel: () => void;
}

const ModuleForm = ({ onSuccess, onCancel }: ModuleFormProps) => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [nameError, setNameError] = useState<string | null>(null);
  const [descriptionError, setDescriptionError] = useState<string | null>(null);

  const { createModule, error, loading } = useCreateModule();

  function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    let valid = true;

    if (!name.trim()) {
      setNameError("Name is required");
      valid = false;
    } else if (name.trim().length > 100) {
      setNameError("Name must be 100 characters or less");
      valid = false;
    } else {
      setNameError(null);
    }

    if (!description.trim()) {
      setDescriptionError("Description is required");
      valid = false;
    } else {
      setDescriptionError(null);
    }

    if (!valid) return;

    void createModule(name, description).then((module) => {
      if (module) {
        setName("");
        setDescription("");
        setNameError(null);
        setDescriptionError(null);
        onSuccess(module);
      }
    });
  }
  return (
    <form onSubmit={handleSubmit}>
      <Dialog.Title className={styles.formTitle}>Register Module</Dialog.Title>
      <Dialog.Description>
        <VisuallyHidden.Root>A form for registering a new module in the orbital command system.</VisuallyHidden.Root>
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
      {nameError && <p className={styles.error}>{nameError}</p>}

      <label htmlFor="description" className={styles.label}>
        Description
      </label>
      <input
        id="description"
        className={styles.input}
        value={description}
        onChange={(e) => {
          setDescription(e.target.value);
        }}
      />
      {descriptionError && <p className={styles.error}>{descriptionError}</p>}

      {error && <p className={styles.errorButtons}>{error}</p>}
      <div className={styles.buttonRow}>
        <Button type="submit" loading={loading}>
          {loading && <span className={styles.spinner}></span>}
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

export { ModuleForm };
