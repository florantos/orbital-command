import { Tabs } from "radix-ui";

import { CrewPanel } from "../../features/crew";
import { ModulesPanel } from "../../features/modules";
import styles from "./AdminPage.module.css";

function AdminPage() {
  return (
    <Tabs.Root defaultValue="crew">
      <Tabs.List className={styles.adminTabs} aria-label="Manage your space station">
        <Tabs.Trigger className={styles.adminTab} value="crew">
          Crew Roster
        </Tabs.Trigger>
        <Tabs.Trigger className={styles.adminTab} value="modules">
          Station Modules
        </Tabs.Trigger>
      </Tabs.List>
      <Tabs.Content value="crew">
        <CrewPanel />
      </Tabs.Content>
      <Tabs.Content value="modules">
        <ModulesPanel />
      </Tabs.Content>
    </Tabs.Root>
  );
}

export { AdminPage };
