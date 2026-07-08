import { useState, useRef } from 'react';
import { Modal, Button, Space, Typography, InputNumber, message } from 'antd';
import html2canvas from 'html2canvas';
import jsPDF from 'jspdf';
import QRLabel from './QRLabel';
import LabelSizeSelector from './LabelSizeSelector';
import type { LabelSize } from './LabelSizeSelector';
import type { Item } from '../../../api/items';

interface QRLabelModalProps {
  items: Item[];
  isOpen: boolean;
  onClose: () => void;
}

export default function QRLabelModal({ items, isOpen, onClose }: QRLabelModalProps) {
  const [labelSize, setLabelSize] = useState<LabelSize>('small');
  const [quantity, setQuantity] = useState<number>(1);
  const printAreaRef = useRef<HTMLDivElement>(null);

  const handlePrint = () => {
    window.print();
  };

  const handleDownloadPDF = async () => {
    if (!printAreaRef.current) return;

    try {
      const canvas = await html2canvas(printAreaRef.current, {
        scale: 2,
        useCORS: true,
      });

      const imgData = canvas.toDataURL('image/png');
      const pdf = new jsPDF({
        orientation: 'landscape',
        unit: 'mm',
        format: 'a4',
      });

      const imgWidth = 297;
      const imgHeight = (canvas.height * imgWidth) / canvas.width;
      pdf.addImage(imgData, 'PNG', 0, 0, imgWidth, imgHeight);
      pdf.save('qr-labels.pdf');
    } catch {
      message.error('Failed to generate PDF. Please try again.');
    }
  };

  return (
    <Modal
      title="QR Code Label Print"
      open={isOpen}
      onCancel={onClose}
      width="90%"
      style={{ maxWidth: 800 }}
      footer={
        <Space>
          <Button onClick={onClose}>Cancel</Button>
          <Button onClick={handleDownloadPDF}>Download PDF</Button>
          <Button type="primary" onClick={handlePrint}>Print</Button>
        </Space>
      }
    >
      <div style={{ marginBottom: 16 }}>
        <Space>
          <Typography.Text>Label Size:</Typography.Text>
          <LabelSizeSelector value={labelSize} onChange={setLabelSize} />
          <Typography.Text>Quantity:</Typography.Text>
          <InputNumber
            min={1}
            max={100}
            value={quantity}
            onChange={(value) => setQuantity(value || 1)}
            style={{ width: 70 }}
          />
        </Space>
      </div>
      
      <div className="print-area" ref={printAreaRef}>
        {items.map((item, index) =>
          Array.from({ length: quantity }).map((_, qtyIndex) => (
            <QRLabel key={`${item.sku}-${index}-${qtyIndex}`} item={item} size={labelSize} />
          ))
        )}
      </div>
    </Modal>
  );
}
