# D3 Transaction Wrapping

`importSecondarySales` wraps all operations in a single `WithinTransaction`:
- Lock → Delete detail → Delete header → Insert header → Insert detail
- Any failure → rollback entire transaction
