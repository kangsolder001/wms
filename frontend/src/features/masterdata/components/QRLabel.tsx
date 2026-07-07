import { QRCodeCanvas as QRCode } from 'qrcode.react';
import type { LabelSize } from './LabelSizeSelector';
import '../styles/print-labels.css';

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
          size={size === 'small' ? 60 : size === 'medium' ? 80 : 100}
          level="M"
        />
        <div className="label-info">
          <div className="label-sku">{item.sku}</div>
          <div className="label-name">{item.name}</div>
          <div className="label-category">{item.category}</div>
        </div>
      </div>
    </div>
  );
}
