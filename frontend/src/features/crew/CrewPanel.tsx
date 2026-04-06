import { Dialog } from "radix-ui";
import { useState } from "react";

import { Button } from "../../components/Button/Button";
import { CrewMemberForm } from "./CrewMemberForm";
import styles from "./CrewPanel.module.css";
import { useReadCrew } from "./useReadCrew";

function CrewPanel() {
  const [open, setOpen] = useState(false);

  const { crew, error, loading, refetch } = useReadCrew();

  return (
    <>
      <div className={styles.header}>
        <span className={styles.count}>{`${String(crew.length)} Crew Members`}</span>
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
                  refetch();
                }}
                onCancel={() => {
                  setOpen(false);
                }}
              />
            </Dialog.Content>
          </Dialog.Portal>
        </Dialog.Root>
      </div>
      <div className={styles.list}>
        {loading && (
          <>
            <div className={styles.moduleList}>
              <div className={styles.cardLoading}>
                <div className={styles.cardLoadingSkelName} />
                <div className={styles.cardLoadingSkelDescription} />
              </div>
            </div>
            <div className={styles.moduleList}>
              <div className={styles.cardLoading}>
                <div className={styles.cardLoadingSkelName} />
                <div className={styles.cardLoadingSkelDescription} />
              </div>
            </div>
            <div className={styles.moduleList}>
              <div className={styles.cardLoading}>
                <div className={styles.cardLoadingSkelName} />
                <div className={styles.cardLoadingSkelDescription} />
              </div>
            </div>
            <div className={styles.moduleList}>
              <div className={styles.cardLoading}>
                <div className={styles.cardLoadingSkelName} />
                <div className={styles.cardLoadingSkelDescription} />
              </div>
            </div>
          </>
        )}
        {error && (
          <div className={styles.cardError}>
            <div className={styles.cardErrorTitle}>Connection Error</div>
            <div className={styles.cardErrorDescription}>
              Unable to retrieve crew data from station core. Retrying...
            </div>
            <button
              type="button"
              className={styles.cardErrorButton}
              onClick={() => {
                refetch();
              }}
            >
              Retry
            </button>
          </div>
        )}
        {!loading && !error && crew.length === 0 && <>No Crew Members</>}
        {!loading &&
          !error &&
          crew.length !== 0 &&
          crew.map((cm) => {
            return (
              <div className={styles.card} key={cm.id}>
                <div>
                  <span className={styles.cardName}>{cm.name}</span>
                  <span className={styles.cardRole}>{cm.role}</span>
                </div>
                <div className={styles.cardQualifications}>
                  {cm.qualifications.map((q) => {
                    return (
                      <span className={styles.cardQualification} key={q}>
                        {q}
                      </span>
                    );
                  })}
                </div>
              </div>
            );
          })}
      </div>
    </>
  );
}

export { CrewPanel };
