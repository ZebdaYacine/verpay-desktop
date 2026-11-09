import { FileUp } from "lucide-react";

export function FileImportMode() {
  return (
    <div className="file-mode" aria-label="Mode d’import">
      <FileUp aria-hidden="true" />
      Importer un fichier
    </div>
  );
}
