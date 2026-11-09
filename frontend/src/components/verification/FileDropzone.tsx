import type { DragEvent } from "react";
import { UploadCloud } from "lucide-react";

type FileDropzoneProps = {
  file: File | null;
  isDragging: boolean;
  onDragStateChange: (isDragging: boolean) => void;
  onFileSelect: (file?: File) => void;
};

const formatFileSize = (size: number) => {
  if (size < 1024) return `${size} octets`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} Ko`;
  return `${(size / (1024 * 1024)).toFixed(1)} Mo`;
};

export function FileDropzone({
  file,
  isDragging,
  onDragStateChange,
  onFileSelect,
}: FileDropzoneProps) {
  const handleDrop = (event: DragEvent<HTMLLabelElement>) => {
    event.preventDefault();
    onDragStateChange(false);
    onFileSelect(event.dataTransfer.files[0]);
  };

  return (
    <label
      className={`dropzone${isDragging ? " is-dragging" : ""}`}
      htmlFor="s2ir-file"
      onDragEnter={(event) => {
        event.preventDefault();
        onDragStateChange(true);
      }}
      onDragOver={(event) => event.preventDefault()}
      onDragLeave={(event) => {
        if (event.currentTarget === event.target) onDragStateChange(false);
      }}
      onDrop={handleDrop}
    >
      <span className="dropzone-icon">
        <UploadCloud aria-hidden="true" />
      </span>
      <strong>{file ? file.name : "Choisir un fichier S2IR"}</strong>
      <span className="dropzone-help">
        {file
          ? `${formatFileSize(file.size)} · cliquez pour remplacer`
          : "ou glissez-déposez le fichier ici"}
      </span>
      <input
        id="s2ir-file"
        type="file"
        accept=".s2ir"
        onChange={(event) => onFileSelect(event.target.files?.[0])}
      />
    </label>
  );
}
