import { CheckCircle2, Hash, ShieldCheck, UploadCloud, X } from "lucide-react";

type VerifiedFileResultProps = {
  fileName: string;
  hash: string;
  isImportReady: boolean;
  period: string;
  payments: string;
  amount: string;
  onClear: () => void;
  onImport: () => void;
  onVerifyAgain: () => void;
};

export function VerifiedFileResult({
  fileName,
  hash,
  isImportReady,
  period,
  payments,
  amount,
  onClear,
  onImport,
  onVerifyAgain,
}: VerifiedFileResultProps) {
  return (
    <section
      className="verified-result"
      aria-label="Résultat de la vérification"
    >
      <div className="hash-summary">
        <span className="hash-summary-icon">
          <Hash aria-hidden="true" />
        </span>
        <div className="hash-summary-copy">
          <strong>{hash}</strong>
          <span>HASH S2IR généré · SHA-256 · correspond au hash attendu</span>
        </div>
        <button
          className="clear-result-button"
          type="button"
          aria-label="Masquer le résultat de vérification"
          onClick={onClear}
        >
          <X aria-hidden="true" />
        </button>
      </div>

      <div className="verified-actions">
        <button className="verify-button" type="button" onClick={onVerifyAgain}>
          <ShieldCheck aria-hidden="true" />
          Vérifier le fichier
        </button>
        <button
          className="import-verified-button"
          type="button"
          onClick={onImport}
        >
          <UploadCloud aria-hidden="true" />
          Générer le fichier vérifié
        </button>
      </div>

      {isImportReady ? (
        <div className="import-ready-message" role="status">
          <CheckCircle2 aria-hidden="true" />
          Fichier vérifié prêt à être généré.
        </div>
      ) : (
        <div className="verification-success" role="status">
          <CheckCircle2 aria-hidden="true" />
          <strong>Fichier vérifié. Vous pouvez le générer.</strong>
        </div>
      )}

      <div className="verification-details">
        <div className="details-heading">
          <span className="details-check">
            <CheckCircle2 aria-hidden="true" />
          </span>
          <div>
            <strong>Fichier vérifié</strong>
            <p>
              Le format S2IR, les totaux, les doublons et le HASH S2IR sont
              valides.
            </p>
          </div>
          <span className="valid-badge">VALIDE</span>
        </div>

        <div className="details-grid">
          <div>
            <span>PÉRIODE</span>
            <strong>{period}</strong>
          </div>
          <div>
            <span>PAIEMENTS</span>
            <strong>{payments}</strong>
          </div>
          <div>
            <span>MONTANT</span>
            <strong>{amount}</strong>
          </div>
          <div>
            <span>HASH</span>
            <strong>{hash.slice(0, 10)}…</strong>
          </div>
        </div>
        <p className="verified-file-name">{fileName}</p>
      </div>
    </section>
  );
}
