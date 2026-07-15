# Attachment Layout Evidence — SX-2045 Payment Recap

Task ID: `20260525-1300-sx-2045-payment-recap`
Tanggal: 2026-05-25 Asia/Jakarta

## Source inspected

- Workbook: `DownloadDepositPayment-250526-003.xlsx`
- Sheet: `Payment Deposit Report`
- Extract command used read-only `openpyxl` against local workbook.
- Inline screenshot from user confirms detail rows should not change and recap block is the changed area.

## Workbook structure

- `A1`: `Payment Deposit Report`
- `A2`: `Deposit Date`
- `B2`: `05-05-2026 - 08-05-2026`
- `A3`: `Collector`
- `B3`: `Jaka`
- Row `5`: detail header from `A` through `Q`:
  - `Deposit Date`
  - `Deposit Type`
  - `Deposit No`
  - `Collector`
  - `Document Date`
  - `Code`
  - `Business Name`
  - `Document No`
  - `Cash`
  - `Cheque / Giro`
  - `Transfer`
  - `Return`
  - `Credit / Debit`
  - `Discount`
  - `Payment Balance`
  - `Expense`
  - `Expense Name`
- Detail rows in sample: `6:15`.
- Recap starts after two blank rows from detail end: sample last detail row `15`, recap header row `18`.

## Recap layout extracted from workbook

### Account Receivable block

- Header: `B18 = Account Receivable`
- Labels: `A19:A26`
- Values: `B19:B26`
- Labels in order:
  1. `Total Cash`
  2. `Total Cheque / Giro`
  3. `Total Transfer`
  4. `Total Return`
  5. `Total Credit / Debit`
  6. `Total Discount`
  7. `Total Payment Balance`
  8. `Total Expense`

Sample values:

```text
B19 35080000
B20 0
B21 3195000
B22 0
B23 0
B24 1009000
B25 1000
B26 -16000
```

### Account Payable block

- Header: `E18 = Account Payable`
- Labels: `D19:D25`
- Values: `E19:E25`
- Labels in workbook order:
  1. `Total Cash`
  2. `Total Cheque / Giro`
  3. `Total Transfer`
  4. `Total Return`
  5. `Total Credit / Debit`
  6. `Total Discount`
  7. `Total Payment Balance`

Sample values:

```text
E19 20000000
E20 0
E21 0
E22 0
E23 0
E24 0
E25 0
D26/E26 blank in workbook
```

## Detail row evidence

- AP detail row remains one normal detail row in sample: row `15`.
- AP cash detail: `I15 = 20000000`.
- AR expense detail row is negative: row `14`, `P14 = -16000`, `Q14 = Makan CGR`.
- Screenshot blue arrow note `Tdk berubah` points at detail AP row, so implementation should not change detail row shape except expense sign fix.

## Source conflict / resolution

- Issue text requires `Total Expense` separated by deposit type.
- Original acceptance criteria says `Total Expense` must be separated and AP expense must not be wrongly counted.
- Workbook has AR `Total Expense` at `A26/B26` but no AP `Total Expense` label/value at `D26/E26`.

Plan resolution:

- Render AR `Total Expense` negative.
- Render AP `Total Expense` as `0` only if final BA/QA expects all required metric rows visible in both blocks.
- If strict workbook layout wins, leave `D26/E26` blank for AP. Implementation should make this a tiny renderer choice, not data model limitation.

Recommended implementation default:

- Data model must always carry `total_expense` per deposit type.
- Excel renderer should follow workbook cell layout unless BA confirms AP `Total Expense` row must appear.
- View `summary_by_deposit_type` should include AP `total_expense: 0` because API contract can represent required metric without Excel blank-cell ambiguity.
