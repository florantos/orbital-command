import clsx from "clsx";
import type React from "react";

import styles from "./Button.module.css";

type variant = "primary" | "secondary" | "danger" | "ghost";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: variant;
  mini?: boolean;
  loading?: boolean;
}

function Button({ variant = "primary", mini = false, loading = false, children, ...props }: ButtonProps) {
  return (
    <button className={clsx(styles[variant], { [styles.mini]: mini })} disabled={loading || props.disabled} {...props}>
      {loading && <span />}
      {children}
    </button>
  );
}

export { Button };
