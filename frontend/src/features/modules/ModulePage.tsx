import { Dialog } from "radix-ui";
import { useState } from "react";

import { ModuleForm } from "./ModuleForm";
import styles from "./ModulePage.module.css";
import type { Module } from "./module.types";

function ModulePage() {
  const [open, setOpen] = useState(false);
  const [createdModule, setCreatedModule] = useState<Module | null>(null);
  return (
    <>
      <div className={styles.header}>
        <span className={styles.moduleCount}>0 Modules</span>
        <Dialog.Root open={open} onOpenChange={setOpen}>
          <Dialog.Trigger asChild>
            <button className={styles.registerModuleButton}>+ Register Module</button>
          </Dialog.Trigger>
          <Dialog.Portal>
            <Dialog.Overlay className={styles.dialogOverlay} />
            <Dialog.Content className={styles.dialogContent}>
              <ModuleForm
                onSuccess={(module) => {
                  setOpen(false);
                  setCreatedModule(module);
                }}
                onCancel={() => {
                  setOpen(false);
                }}
              />
            </Dialog.Content>
          </Dialog.Portal>
        </Dialog.Root>
      </div>
      <div>{createdModule?.name}</div>
    </>
  );
}

export { ModulePage };
