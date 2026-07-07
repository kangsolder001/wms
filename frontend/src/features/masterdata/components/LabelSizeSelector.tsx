import { Select } from 'antd';

export type LabelSize = 'small' | 'medium' | 'large';

interface LabelSizeSelectorProps {
  value: LabelSize;
  onChange: (size: LabelSize) => void;
}

const sizeOptions = [
  { value: 'small', label: 'Small (5x3cm)' },
  { value: 'medium', label: 'Medium (7x5cm)' },
  { value: 'large', label: 'Large (10x7cm)' },
];

export default function LabelSizeSelector({ value, onChange }: LabelSizeSelectorProps) {
  return (
    <Select
      value={value}
      onChange={onChange}
      options={sizeOptions}
      style={{ width: 150 }}
    />
  );
}
