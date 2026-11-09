import { ShieldCheck } from "lucide-react";

type VerifyButtonProps = {
  isChecking: boolean;
  onClick: () => void;
};

export function VerifyButton({ isChecking, onClick }: VerifyButtonProps) {
  return (
    <button
      className="verify-button"
      type="button"
      onClick={onClick}
      disabled={isChecking}
    >
      <ShieldCheck aria-hidden="true" />
      {isChecking ? "Vérification en cours…" : "Vérifier le fichier"}
    </button>
  );
}
