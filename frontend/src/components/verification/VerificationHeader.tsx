import { FileCheck2 } from "lucide-react";

export function VerificationHeader() {
  return (
    <header className="verification-header">
      <div className="eyebrow">
        <FileCheck2 aria-hidden="true" />
        Import sécurisé
      </div>
      <h1 id="page-title">Vérifier un fichier avant l’import</h1>
      <p>
        Sélectionnez un fichier S2IR. Nous vérifions le format, les totaux, les
        doublons et le HASH S2IR.
      </p>
    </header>
  );
}
