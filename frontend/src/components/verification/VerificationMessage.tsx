import { AlertCircle, CheckCircle2, XCircle } from "lucide-react";
import type { VerificationState } from "./types";

type VerificationMessageProps = {
  state: VerificationState;
  message: string;
};

export function VerificationMessage({
  state,
  message,
}: VerificationMessageProps) {
  const StatusIcon =
    state === "valid"
      ? CheckCircle2
      : state === "invalid"
        ? XCircle
        : AlertCircle;

  return (
    <div className={`verification-message is-${state}`} role="status">
      <StatusIcon aria-hidden="true" />
      <span>{message}</span>
    </div>
  );
}
