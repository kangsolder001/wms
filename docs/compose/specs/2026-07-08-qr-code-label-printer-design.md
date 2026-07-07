# QR Code Label Printer Design

## [S1] Problem

Items in the warehouse need QR code labels for easy identification and scanning. Each item should have a printable label with a QR code containing item information in a human-readable format.

## [S2] Solution Overview

Add a QR code label printing feature to the Items page that allows users to:
- Print QR code labels for single or multiple items
- Choose from three label sizes (Small, Medium, Large)
- Print directly via browser or download as PDF

## [S3] QR Code Content Format

QR code content format: `SKU#name#category` (delimiter: #)

Example: `ELC-001#Laptop ASUS 14"#Electronics`

## [S4] Label Sizes

| Size | Dimensions | Use Case |
|------|------------|----------|
| Small | 50mm x 30mm | Small items, shelf labels |
| Medium | 70mm x 50mm | Standard items |
| Large | 100mm x 70mm | Large items, easy reading |

User can select size via dropdown in the print modal.

## [S5] Component Architecture

```
ItemsPage.tsx
  └─ QRLabelModal.tsx (modal preview + print)
       ├─ QRLabel.tsx (single label component)
       │    ├─ QRCode (qrcode.react)
       │    └─ LabelInfo (SKU, Name, Category)
       └─ LabelSizeSelector.tsx (dropdown)
```

### QRLabel.tsx

Reusable label component that renders a single QR code label.

**Props:**
- `item: { sku: string; name: string; category: string }`
- `size: 'small' | 'medium' | 'large'`

**Behavior:**
- Renders QR code with content `SKU#name#category`
- Renders text info below/beside QR code
- Applies CSS class for print styling
- Includes cutting guide border

### QRLabelModal.tsx

Modal component for previewing and printing labels.

**Props:**
- `items: Array<{ sku: string; name: string; category: string }>`
- `isOpen: boolean`
- `onClose: () => void`

**Behavior:**
- Shows label size selector dropdown
- Previews all labels in selected size
- Print button triggers `window.print()`
- Download PDF button uses html2canvas + jsPDF

### LabelSizeSelector.tsx

Dropdown component for selecting label size.

**Props:**
- `value: 'small' | 'medium' | 'large'`
- `onChange: (size: string) => void`

## [S6] Data Flow

### Single Item Print

1. User clicks "Print Label" button in items table row
2. QRLabelModal opens with single item data
3. User selects label size
4. User clicks "Print" or "Download PDF"

### Batch Print

1. User selects multiple items via checkboxes in table
2. User clicks "Print Selected" button above table
3. QRLabelModal opens with all selected items
4. User selects label size
5. User clicks "Print" or "Download PDF"

### Print Process

1. `window.print()` triggers browser print dialog
2. CSS `@media print` hides all elements except `.print-area`
3. Labels render at selected size
4. User confirms print in browser dialog

### PDF Export

1. `html2canvas` captures label area as image
2. `jsPDF` creates PDF document
3. Labels are arranged on PDF pages
4. Browser downloads the PDF file

## [S7] UI Design

### Items Table Modifications

- Add "Print Label" button in Actions column
- Add checkbox selection for batch operations
- Add "Print Selected" button above table (visible when items selected)

### QRLabelModal Layout

```
┌─────────────────────────────────────┐
│ QR Code Label Print           [X]   │
├─────────────────────────────────────┤
│ Label Size: [Small ▼]              │
│                                     │
│ ┌─────────────────────┐            │
│ │  ┌───────────────┐  │            │
│ │  │  QR CODE      │  │            │
│ │  │  SKU#name#cat │  │            │
│ │  └───────────────┘  │            │
│ │  ELC-001            │            │
│ │  Laptop ASUS 14"    │            │
│ │  Electronics        │            │
│ └─────────────────────┘            │
│                                     │
│ [Print] [Download PDF] [Cancel]     │
└─────────────────────────────────────┘
```

### QRLabel Component Layout

- QR code on left (or top for small size)
- Text info on right (or bottom for small size)
- Thin border for cutting guide

## [S8] Implementation Details

### Dependencies

```bash
npm install qrcode.react html2canvas jspdf
```

### Files to Create

1. `frontend/src/features/masterdata/components/QRLabel.tsx`
2. `frontend/src/features/masterdata/components/QRLabelModal.tsx`
3. `frontend/src/features/masterdata/components/LabelSizeSelector.tsx`
4. `frontend/src/features/masterdata/styles/print-labels.css`

### Files to Modify

1. `frontend/src/features/masterdata/ItemsPage.tsx`
   - Add print label button in Actions column
   - Add checkbox selection for batch print
   - Add "Print Selected" button above table
   - Import QRLabelModal component

### Print CSS

```css
@media print {
  body * { visibility: hidden; }
  .print-area, .print-area * { visibility: visible; }
  .print-area { position: absolute; left: 0; top: 0; }
}

.label-container {
  border: 1px dashed #ccc;
  padding: 8px;
  display: inline-block;
}

.label-small { width: 50mm; height: 30mm; }
.label-medium { width: 70mm; height: 50mm; }
.label-large { width: 100mm; height: 70mm; }
```

### Label Size CSS Classes

- `.label-small` - 50mm x 30mm
- `.label-medium` - 70mm x 50mm
- `.label-large` - 100mm x 70mm

## [S9] Testing

### Manual Testing

1. Single item print:
   - Navigate to Items page
   - Click "Print Label" button on any item
   - Verify modal opens with correct item data
   - Select different sizes, verify preview updates
   - Click "Print", verify browser print dialog opens
   - Click "Download PDF", verify PDF downloads

2. Batch print:
   - Select multiple items via checkboxes
   - Click "Print Selected" button
   - Verify modal opens with all selected items
   - Print/Download, verify all labels included

3. Print styling:
   - Verify only labels are visible in print preview
   - Verify labels render at correct sizes
   - Verify QR code is scannable

### Acceptance Criteria

- [ ] QR code contains `SKU#name#category` format
- [ ] User can select label size (Small/Medium/Large)
- [ ] Single item print works via table button
- [ ] Batch print works via checkbox selection
- [ ] Print via browser works correctly
- [ ] PDF download works correctly
- [ ] Labels render at correct dimensions
- [ ] QR codes are scannable
