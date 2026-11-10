import { useState } from "react";
import {
  FileVerificationService,
  type FileVerificationResult,
} from "../bindings/changeme";

import {
  FileDropzone,
  FileImportMode,
  HashField,
  VerificationHeader,
  VerificationMessage,
  VerifiedFileResult,
  VerifyButton,
  type VerificationState,
} from "@/components/verification";

const fileToBase64 = async (file: File): Promise<string> => {
  const bytes = new Uint8Array(await file.arrayBuffer());
  const chunkSize = 0x8000;

  let binary = "";

  for (let index = 0; index < bytes.length; index += chunkSize) {
    binary += String.fromCharCode(
      ...bytes.subarray(index, Math.min(index + chunkSize, bytes.length)),
    );
  }

  return btoa(binary);
};

function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [expectedHash, setExpectedHash] = useState("");
  const [btnStatus, setBtnStatus] = useState("disabled");

  const [result, setResult] = useState<FileVerificationResult | null>(null);

  const [isDragging, setIsDragging] = useState(false);

  const [verificationState, setVerificationState] =
    useState<VerificationState>("idle");

  const [message, setMessage] = useState("");
  const [calculatedHash, setCalculatedHash] = useState("");
  const [isImportReady, setIsImportReady] = useState(false);

  // ---------------------------------------------------
  // RESET VERIFICATION
  // ---------------------------------------------------

  const resetVerification = () => {
    setVerificationState("idle");
    setMessage("");
    setCalculatedHash("");
    setResult(null);
    setIsImportReady(false);
  };

  // ---------------------------------------------------
  // SELECT FILE
  // ---------------------------------------------------

  const selectFile = (file?: File) => {
    if (!file || verificationState === "checking") return;

    setSelectedFile(file);
    resetVerification();
    setBtnStatus("disabled");
  };

  // ---------------------------------------------------
  // VERIFY FILE
  // ---------------------------------------------------

  const verifyFile = async () => {
    if (!selectedFile || verificationState === "checking") {
      if (!selectedFile) {
        setVerificationState("invalid");
        setMessage("Sélectionnez un fichier S2IR.");
      }
      return;
    }

    const file = selectedFile;
    const normalizedHash = expectedHash.trim().toLowerCase();

    setVerificationState("checking");
    setMessage("");
    setResult(null);
    setCalculatedHash("");
    setIsImportReady(false);

    try {
      // Convert file to Base64
      const encodedFile = await fileToBase64(file);

      // Call Wails Go backend
      const response = await FileVerificationService.VerifyFile(
        file.name,
        encodedFile,
        normalizedHash,
      );

      const verificationResult = response as FileVerificationResult;

      setResult(verificationResult);

      setCalculatedHash(verificationResult.hash ?? "");

      if (!verificationResult.valid) {
        setVerificationState("invalid");

        setMessage(
          verificationResult.message || "Le fichier S2IR n'est pas valide.",
        );

        return;
      }

      setVerificationState("valid");
      setBtnStatus("enabled");

      setMessage(
        normalizedHash
          ? "Fichier et empreinte SHA-256 vérifiés avec succès."
          : "Fichier S2IR vérifié avec succès.",
      );
    } catch (error: unknown) {
      console.error("File verification error:", error);

      setVerificationState("invalid");

      setMessage(
        error instanceof Error
          ? error.message
          : typeof error === "string"
            ? error
            : "Impossible de vérifier ce fichier.",
      );
    }
  };

  const GenerateOutPutFile = async () => {
    const file = selectedFile!;
    const normalizedHash = expectedHash.trim().toLowerCase();

    setVerificationState("checking");
    setMessage("");
    setResult(null);
    setCalculatedHash("");
    setIsImportReady(false);

    try {
      // Convert file to Base64
      const encodedFile = await fileToBase64(file);

      // Call Wails Go backend
      const response = await FileVerificationService.GenerateOutPutFile(
        file.name,
        encodedFile,
        normalizedHash,
      );

      alert(response);

      setMessage(
        normalizedHash
          ? "Fichier et empreinte SHA-256 vérifiés avec succès."
          : "Fichier S2IR vérifié avec succès.",
      );
    } catch (error: unknown) {
      console.error("File verification error:", error);

      setMessage(
        error instanceof Error
          ? error.message
          : typeof error === "string"
            ? error
            : "Impossible de vérifier ce fichier.",
      );
    }
  };
  // ---------------------------------------------------
  // UI
  // ---------------------------------------------------

  return (
    <main className="verification-page">
      <section className="verification-card" aria-labelledby="page-title">
        <VerificationHeader />

        <div className="verification-content">
          <FileImportMode />

          <FileDropzone
            file={selectedFile}
            isDragging={isDragging}
            onDragStateChange={setIsDragging}
            onFileSelect={selectFile}
          />

          <HashField
            value={expectedHash}
            onChange={setExpectedHash}
            onEdit={resetVerification}
          />

          {verificationState === "valid" && result ? (
            <VerifiedFileResult
              fileName={result.fileName || selectedFile?.name || ""}
              hash={calculatedHash}
              isImportReady={isImportReady}
              period={result.period}
              payments={String(result.paymentCount)}
              amount={result.totalAmount}
              onClear={resetVerification}
              onImport={() => GenerateOutPutFile()}
              onVerifyAgain={verifyFile}
              status={btnStatus}
            />
          ) : (
            <>
              <VerifyButton
                isChecking={verificationState === "checking"}
                onClick={verifyFile}
              />

              {message && (
                <VerificationMessage
                  state={verificationState}
                  message={message}
                />
              )}
            </>
          )}
        </div>
      </section>

      <p className="privacy-note">
        Rien n'est enregistré avant la validation du fichier, du HASH S2IR et
        des doublons.
      </p>
    </main>
  );
}

export default App;
