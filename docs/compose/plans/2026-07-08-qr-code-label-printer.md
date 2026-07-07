# QR Code Label Printer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use compose:subagent (recommended) or compose:execute to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add QR code label printing feature to Items page with single and batch print support.

**Architecture:** Client-side QR code generation using qrcode.react, with print via browser native print API and PDF export via html2canvas + jsPDF. Component-based architecture with reusable QRLabel component.

**Tech Stack:** React, TypeScript, Ant Design, qrcode.react, html2canvas, jsPDF

---

### Task 1: Install Dependencies

**Covers:** [S8]

**Files:**
- Modify: `frontend/package.json`

- [ ] **Step 1: Install qrcode.react, html2canvas, and jspdf**

Run: `cd frontend && npm install qrcode.react html2canvas jspdf`

Expected: Dependencies added to package.json and node_modules

- [ ] **Step 2: Verify installation**

Run: `cd frontend && npm list qrcode.react html2canvas jspdf`

Expected: All three packages listed with versions

- [ ] **Step 3: Commit**

```bash
git add frontend/package.json frontend/package-lock.json
git commit -m "feat: add QR code label printing dependencies"
```

---

### Task 2: Create LabelSizeSelector Component

**Covers:** [S4, S5]

**Files:**
- Create: `frontend/src/features/masterdata/components/LabelSizeSelector.tsx`

- [ ] **Step 1: Create LabelSizeSelector component**

```tsx
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
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/features/masterdata/components/LabelSizeSelector.tsx
git commit -m "feat: add LabelSizeSelector component"
```

---

### Task 3: Create QRLabel Component

**Covers:** [S3, S5, S7]

**Files:**
- Create: `frontend/src/features/masterdata/components/QRLabel.tsx`
- Create: `frontend/src/features/masterdata/styles/print-labels.css`

- [ ] **Step 1: Create print-labels.css**

```css
.label-container {
  border: 1px dashed #ccc;
  padding: 8px;
  display: inline-block;
  background: white;
}

.label-small {
  width: 50mm;
  height: 30mm;
}

.label-medium {
  width: 70mm;
  height: 50mm;
}

.label-large {
  width: 100mm;
  height: 70mm;
}

.label-content {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 100%;
}

.label-small .label-content {
  flex-direction: column;
  gap: 4px;
}

.label-info {
  flex: 1;
  overflow: hidden;
}

.label-sku {
  font-weight: bold;
  font-size: 12px;
  margin-bottom: 2px;
}

.label-name {
  font-size: 10px;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.label-category {
  font-size: 9px;
  color: #666;
}

.label-small .label-sku {
  font-size: 10px;
}

.label-small .label-name {
  font-size: 8px;
}

.label-small .label-category {
  font-size: 7px;
}

@media print {
  body * {
    visibility: hidden;
  }
  .print-area,
  .print-area * {
    visibility: visible;
  }
  .print-area {
    position: absolute;
    left: 0;
    top: 0;
  }
  .label-container {
    border: 1px dashed #ccc;
    page-break-inside: avoid;
  }
}
```

- [ ] **Step 2: Create QRLabel component**

```tsx
import QRCode from 'qrcode.react';
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
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/features/masterdata/components/QRLabel.tsx frontend/src/features/masterdata/styles/print-labels.css
git commit -m "feat: add QRLabel component with print styling"
```

---

### Task 4: Create QRLabelModal Component

**Covers:** [S5, S6, S7]

**Files:**
- Create: `frontend/src/features/masterdata/components/QRLabelModal.tsx`

- [ ] **Step 1: Create QRLabelModal component**

```tsx
import { useState, useRef } from 'react';
import { Modal, Button, Space, Typography } from 'antd';
import html2canvas from 'html2canvas';
import jsPDF from 'jspdf';
import QRLabel from './QRLabel';
import LabelSizeSelector from './LabelSizeSelector';
import type { LabelSize } from './LabelSizeSelector';

interface Item {
  sku: string;
  name: string;
  category: string;
}

interface QRLabelModalProps {
  items: Item[];
  isOpen: boolean;
  onClose: () => void;
}

export default function QRLabelModal({ items, isOpen, onClose }: QRLabelModalProps) {
  const [labelSize, setLabelSize] = useState<LabelSize>('small');
  const printAreaRef = useRef<HTMLDivElement>(null);

  const handlePrint = () => {
    window.print();
  };

  const handleDownloadPDF = async () => {
    if (!printAreaRef.current) return;
    
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
        </Space>
      </div>
      
      <div className="print-area" ref={printAreaRef}>
        <Space wrap>
          {items.map((item, index) => (
            <QRLabel key={index} item={item} size={labelSize} />
          ))}
        </Space>
      </div>
    </Modal>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/features/masterdata/components/QRLabelModal.tsx
git commit -m "feat: add QRLabelModal with print and PDF export"
```

---

### Task 5: Update ItemsPage with Print Functionality

**Covers:** [S6, S7]

**Files:**
- Modify: `frontend/src/features/masterdata/ItemsPage.tsx`

- [ ] **Step 1: Update ItemsPage to add print functionality**

```tsx
import { useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Typography, message, Checkbox } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, PrinterOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { itemApi } from '../../api/items';
import type { Item, CreateItemRequest } from '../../api/items';
import type { ColumnsType } from 'antd/es/table';
import QRLabelModal from './components/QRLabelModal';

export default function ItemsPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Item | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [printModalItems, setPrintModalItems] = useState<Item[]>([]);
  const [isPrintModalOpen, setIsPrintModalOpen] = useState(false);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemApi.list().then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateItemRequest) => itemApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      message.success('Item created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create item'),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => itemApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      message.success('Item deleted');
    },
  });

  const handlePrintSingle = (item: Item) => {
    setPrintModalItems([item]);
    setIsPrintModalOpen(true);
  };

  const handlePrintSelected = () => {
    const selectedItems = data?.data?.filter((item: Item) => 
      selectedRowKeys.includes(item.id)
    ) || [];
    setPrintModalItems(selectedItems);
    setIsPrintModalOpen(true);
    setSelectedRowKeys([]);
  };

  const columns: ColumnsType<Item> = [
    { title: 'SKU', dataIndex: 'sku', key: 'sku' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Category', dataIndex: 'category', key: 'category', responsive: ['md'] },
    { title: 'UoM', dataIndex: 'unit_of_measure', key: 'unit_of_measure', responsive: ['md'] },
    { title: 'Weight', dataIndex: 'weight', key: 'weight', responsive: ['lg'] },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 150,
      render: (_: any, record: Item) => (
        <Space>
          <Button icon={<PrinterOutlined />} onClick={() => handlePrintSingle(record)} />
          <Button icon={<EditOutlined />} onClick={() => { setEditingItem(record); setIsModalOpen(true); }} />
          <Button icon={<DeleteOutlined />} danger onClick={() => deleteMutation.mutate(record.id)} />
        </Space>
      ),
    },
  ];

  const handleSubmit = (values: CreateItemRequest) => {
    createMutation.mutate(values);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Items</Typography.Title>
        <Space>
          {selectedRowKeys.length > 0 && (
            <Button icon={<PrinterOutlined />} onClick={handlePrintSelected}>
              Print Selected ({selectedRowKeys.length})
            </Button>
          )}
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingItem(null); setIsModalOpen(true); }}>
            Add Item
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={data?.data || []}
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
        scroll={{ x: 'max-content }}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
      />

      <Modal
        title={editingItem ? 'Edit Item' : 'Add Item'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingItem(null); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={editingItem || {}}>
          <Form.Item name="sku" label="SKU" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="category" label="Category">
            <Input />
          </Form.Item>
          <Form.Item name="unit_of_measure" label="Unit of Measure" rules={[{ required: true }]}>
            <Select options={[{ value: 'pcs', label: 'Pieces' }, { value: 'kg', label: 'Kilogram' }, { value: 'box', label: 'Box' }]} />
          </Form.Item>
          <Form.Item name="weight" label="Weight">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>

      <QRLabelModal
        items={printModalItems}
        isOpen={isPrintModalOpen}
        onClose={() => setIsPrintModalOpen(false)}
      />
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/features/masterdata/ItemsPage.tsx
git commit -m "feat: add QR code print functionality to Items page"
```

---

### Task 6: Build and Test

**Covers:** [S9]

**Files:**
- None (testing only)

- [ ] **Step 1: Build frontend**

Run: `cd frontend && npm run build`

Expected: Build succeeds without errors

- [ ] **Step 2: Start development server**

Run: `cd frontend && npm run dev`

Expected: Server starts on http://localhost:5173

- [ ] **Step 3: Manual testing**

1. Navigate to Items page
2. Click "Print Label" button on any item
3. Verify modal opens with correct item data
4. Select different sizes, verify preview updates
5. Click "Print", verify browser print dialog opens
6. Click "Download PDF", verify PDF downloads
7. Select multiple items via checkboxes
8. Click "Print Selected" button
9. Verify modal opens with all selected items
10. Print/Download, verify all labels included

- [ ] **Step 4: Verify QR codes are scannable**

1. Print a label
2. Scan QR code with phone camera
3. Verify content matches `SKU#name#category` format

- [ ] **Step 5: Final commit**

```bash
git add .
git commit -m "feat: QR code label printer complete"
```

---

## Self-Review

**1. Spec coverage:**
- [S1] Problem - Covered by overall feature
- [S2] Solution Overview - Covered by architecture
- [S3] QR Code Content Format - Covered by Task 3
- [S4] Label Sizes - Covered by Task 2
- [S5] Component Architecture - Covered by Tasks 2-4
- [S6] Data Flow - Covered by Tasks 4-5
- [S7] UI Design - Covered by Tasks 3-5
- [S8] Implementation Details - Covered by Tasks 1-5
- [S9] Testing - Covered by Task 6

**2. Placeholder scan:** No TBD/TODO found

**3. Type consistency:** All types match across components
