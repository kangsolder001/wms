import { QRCodeCanvas as QRCode } from 'qrcode.react';
import type { LabelSize } from './LabelSizeSelector';
import '../styles/print-labels.css';

const qrSizes: Record<LabelSize, number> = { small: 60, medium: 80, large: 100 };

interface QRLabelProps {
  item: {
    sku: string;
    name: string;
    category: string;
  };
  size: LabelSize;
}

export default function QRLabel({ item, size }: QRLabelProps) {
  const qrContent = `${item.sku}#${item.name}#${item.category}`;

  return (
    <div className={`label-container label-${size}`}>
      <div className="label-content">
        <QRCode
          value={qrContent}
          size={qrSizes[size]}
          level="M"
        />
        <div className="label-info">
          <div>{item.sku}</div>
          <div>{item.name}</div>
          <div>{item.category}</div>
        </div>
      </div>
    </div>
  );
}
