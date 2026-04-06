import { Dialog } from "radix-ui";
import { useState } from "react";

import { Button } from "../../components/Button/Button";
import { CrewMemberForm } from "./CrewMemberForm";
import styles from "./CrewPage.module.css";

function CrewPanel() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <div className={styles.header}>
        <span className={styles.count}></span>
        <Dialog.Root open={open} onOpenChange={setOpen}>
          <Dialog.Trigger asChild>
            <Button mini>+ Register Crew Member</Button>
          </Dialog.Trigger>
          <Dialog.Portal>
            <Dialog.Overlay className={styles.dialogOverlay} />
            <Dialog.Content className={styles.dialogContent}>
              <CrewMemberForm
                onSuccess={() => {
                  setOpen(false);
                }}
                onCancel={() => {
                  setOpen(false);
                }}
              />
            </Dialog.Content>
          </Dialog.Portal>
        </Dialog.Root>
      </div>
    </>
  );
}

export { CrewPanel };
