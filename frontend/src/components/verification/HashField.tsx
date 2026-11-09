import { Hash } from "lucide-react";

type HashFieldProps = {
  value: string;
  onChange: (value: string) => void;
  onEdit: () => void;
};

export function HashField({ value, onChange, onEdit }: HashFieldProps) {
  return (
    <div className="field-group">
      <label htmlFor="expected-hash">HASH S2IR ATTENDU (OPTIONNEL)</label>
      <div className="hash-input-wrap">
        <Hash aria-hidden="true" />
        <input
          id="expected-hash"
          type="text"
          inputMode="text"
          autoComplete="off"
          spellCheck="false"
          maxLength={64}
          value={value}
          onChange={(event) => {
            onChange(event.target.value.replace(/\s/g, ""));
            onEdit();
          }}
          placeholder="Collez le HASH S2IR de 64 caractères (optionnel)"
        />
      </div>
      <p className="field-help">
        Laissez vide pour vérifier le fichier sans comparer de hash.
      </p>
    </div>
  );
}
