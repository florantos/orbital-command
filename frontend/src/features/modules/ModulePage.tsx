import { Dialog } from "radix-ui";
import { useState } from "react";

import { ModuleForm } from "./ModuleForm";
import styles from "./ModulePage.module.css";
import { useReadModules } from "./useReadModules";

function ModulePage() {
  const [open, setOpen] = useState(false);

  const { modules, error, loading, refetch } = useReadModules();

  return (
    <>
      <div className={styles.header}>
        <span className={styles.moduleCount}>{`${String(modules.length)} Modules`}</span>
        <Dialog.Root open={open} onOpenChange={setOpen}>
          <Dialog.Trigger asChild>
            <button className={styles.registerModuleButton}>+ Register Module</button>
          </Dialog.Trigger>
          <Dialog.Portal>
            <Dialog.Overlay className={styles.dialogOverlay} />
            <Dialog.Content className={styles.dialogContent}>
              <ModuleForm
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

      <div className={styles.moduleList}>
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
              Unable to retrieve module data from station core. Retrying...
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
        {modules.length === 0 && <>No Modules</>}
        {modules.length !== 0 &&
          modules.map((m) => {
            const healthStateClass = {
              operational: styles.healthOperational,
              degraded: styles.healthDegraded,
              critical: styles.healthCritical,
              unresponsive: styles.healthUnresponsive,
              offline: styles.healthOffline,
            }[m.healthState];

            const healthBorderClass = {
              operational: styles.healthBorderOperational,
              degraded: styles.healthBorderDegraded,
              critical: styles.healthBorderCritical,
              unresponsive: styles.healthBorderUnresponsive,
              offline: styles.healthBorderOffline,
            }[m.healthState];

            return (
              <div className={`${styles.card} ${healthBorderClass}`} key={m.id}>
                <div>
                  <span className={styles.cardName}>{m.name}</span>
                  <p className={styles.cardDescription}>{m.description}</p>
                </div>
                <div>
                  <span className={`${styles.cardHealthState} ${healthStateClass}`}>{m.healthState}</span>
                </div>
              </div>
            );
          })}
      </div>
    </>
  );
}

export { ModulePage };
