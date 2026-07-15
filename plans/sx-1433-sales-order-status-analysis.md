# SX-1433 Sales Order Status Analysis

## Context

- Ticket: `SX-1433`
- Issue: Sales Order created from web remains `Need Review` and cannot be processed
- Expected by QA/Jira: `sls.order.data_status` should be stored as `2`
- Actual behavior found in backend:
  - request from web sends `data_status = 1`
  - backend maps request to order model directly
  - repository persists mapped model as-is
  - approval flow is designed around `Need Review -> Approval -> Processed`
  - credit limit scenario on the sample order is over-limit with outlet action `Restricted`

## Final Diagnosis Draft

### 1. What the code currently does

Create and update flow currently treat `data_status` as input-driven header state.

- [`OrderController.Create()`](../sales/controller/order_controller.go:81) validates then forwards request to service
- [`orderServiceImpl.Store()`](../sales/service/order_service.go:43) maps request directly into [`model.Order`](../sales/model/order.go)
- [`RepositoryOrderImpl.Store()`](../sales/repository/order_repository.go:26) inserts the model without normalizing status
- [`OrderController.Update()`](../sales/controller/order_controller.go:154) and [`orderServiceImpl.Update()`](../sales/service/order_service.go:142) repeat the same pattern for update

Approval flow treats `Need Review` as a required precondition.

- [`HierarchyApprovalService.RequestApproval()`](../sales/service/hierarchy_approval_service.go:429) rejects non-`NEED_REVIEW` orders
- [`OrderApprovalService.UpdateStatusDetail()`](../sales/service/order_approval_service.go:68) changes order status to `PROCESSED` only after approval finishes

Credit-limit validation computes review signals, but does not directly set order status.

- [`ValidateOrder()`](../sales/service/validate_order_service.go:51) sets validation flags and messages
- It does not assign `data_status`

### 2. Best classification of the issue

This is **not purely a simple mapping bug** and **not purely a business-rule mismatch**.

Most accurate classification:

1. **Immediate defect**: backend trusts request-level `data_status` too much and persists it directly for create/update web flow
2. **Underlying design issue**: `sls.order.data_status` is being used for two different concerns at once:
   - order lifecycle state
   - approval/review state
3. **Possible requirement mismatch**: QA expectation `data_status = 2` may conflict with the currently implemented approval design for over-limit restricted orders

So the strongest conclusion is:

> Primary classification: **design issue with a manifestation through unsafe mapping/persistence**

Secondary classification:

> There is also a **requirement mismatch risk** if business still expects approval to exist for restricted orders while QA expects header status to be `Processed`

## Decision Analysis

### Option A — Force `data_status = PROCESSED` during web create/update

#### Summary
Set backend status explicitly to `PROCESSED` for web-origin order flow, ignoring incoming `data_status = 1`.

#### Impact on approval flow
- High impact
- [`RequestApproval()`](../sales/service/hierarchy_approval_service.go:429) currently requires `NEED_REVIEW`
- If orders become `PROCESSED` immediately, approval flow may stop working or become logically bypassed
- Existing approval queues may no longer receive orders intended for review

#### Impact on order list/filter
- Any list/filter relying on `data_status = 1` for review candidates will lose visibility
- Any downstream list expecting `PROCESSED` may suddenly include over-limit restricted orders

#### Impact on downstream process
- Risk that downstream processes treat the order as fully processable even when approval should still gate execution
- Could unlock invoice, delivery, or follow-up process too early depending on downstream filters

#### Risks / side effects
- Bypasses current approval semantics
- Breaks approval request precondition
- Introduces mismatch between approval tables and order header state
- Hotfix may appear to solve UI symptom while silently corrupting workflow semantics

#### Implementation shape
- Changes likely in:
  - [`orderServiceImpl.Store()`](../sales/service/order_service.go:43)
  - [`orderServiceImpl.Update()`](../sales/service/order_service.go:142)
  - possibly response builders and approval entry conditions

#### Suitability
- **Hotfix only if business explicitly confirms approval must not block web orders anymore**
- Otherwise unsafe

### Option B — Keep `Need Review` for over-limit restricted orders

#### Summary
Preserve current approval-driven design and classify the ticket as expectation mismatch or upstream payload issue.

#### Impact on approval flow
- Low impact
- Fully aligned with current code path
- No structural change required

#### Impact on order list/filter
- Existing review queues continue to work
- Existing processed-only filters remain unchanged

#### Impact on downstream process
- Safest for operational consistency
- Orders remain blocked until approval finishes, which matches current backend design

#### Risks / side effects
- QA ticket remains unresolved if expected behavior is truly `Processed`
- Web client may still expose confusing state if UI/business now expects otherwise
- Might preserve a flawed design if the real intent is to separate approval state from lifecycle state

#### Implementation shape
- Little to no code change
- May require validation that web should not send `data_status = 1`, or documentation/business clarification

#### Suitability
- **Safest short-term containment**
- Best only if business confirms `Restricted` orders must remain review-gated

### Option C — Separate approval state from `sls.order.data_status`

#### Summary
Treat order lifecycle status and approval workflow status as separate concepts.

Possible target model:
- `sls.order.data_status` tracks lifecycle such as `Processed`
- approval state comes from approval tables or a dedicated field such as `approval_status`

#### Impact on approval flow
- High but correct architectural impact
- Approval entry no longer depends on `data_status = NEED_REVIEW`
- [`RequestApproval()`](../sales/service/hierarchy_approval_service.go:429) must be redesigned to use explicit approval criteria
- approval completion may update approval state only, not header lifecycle status

#### Impact on order list/filter
- Review lists should query approval tables or explicit approval status
- Order lists become more semantically correct
- Need audit of all filters using `data_status` as proxy for approval state

#### Impact on downstream process
- Cleaner contract for downstream consumers
- Downstream processes can rely on lifecycle status and separately check approval state if needed
- Reduces ambiguity and future regression risk

#### Risks / side effects
- Broader refactor surface
- Requires schema/contract audit across API, UI, reports, batch jobs, and operational process
- Must carefully migrate any existing code using `data_status = 1` as business meaning for pending approval

#### Implementation shape
- Likely touches:
  - create/update service logic
  - approval request preconditions
  - approval completion logic
  - list queries for review items
  - response DTOs and status-name generation
  - possibly DB schema or at least response contract rules

#### Suitability
- **Best long-term solution**
- Not ideal as a rushed hotfix unless scope is accepted

## Code Impact Map if `data_status` behavior changes

### Create flow
- [`OrderController.Create()`](../sales/controller/order_controller.go:81)
- [`orderServiceImpl.Store()`](../sales/service/order_service.go:43)
- [`RepositoryOrderImpl.Store()`](../sales/repository/order_repository.go:26)
- [`entity.CreateOrderBody`](../sales/entity/order.go:59)
- [`model.Order`](../sales/model/order.go)

### Update flow
- [`OrderController.Update()`](../sales/controller/order_controller.go:154)
- [`orderServiceImpl.Update()`](../sales/service/order_service.go:142)
- [`RepositoryOrderImpl.Update()`](../sales/repository/order_repository.go:34)
- [`entity.UpdateOrderBody`](../sales/entity/order.go:109)

### Approval request flow
- [`HierarchyApprovalService.RequestApproval()`](../sales/service/hierarchy_approval_service.go:429)
- [`OrderApprovalRequestRepository.FindApprovalProcessedByRoNo()`](../sales/repository/order_approval_request_repository.go:78)

### Approval completion flow
- [`OrderApprovalService.UpdateStatusDetail()`](../sales/service/order_approval_service.go:68)
- [`OrderApprovalRepository.UpdateStatusOrder()`](../sales/repository/order_approval_repository.go:591)

### Query/list “Need Review” and status labels
- [`GenerateDataStatusName()`](../sales/entity/order.go:383)
- approval review queries in [`RepositoryOrderApprovalImpl.FindNeedReview()`](../sales/repository/order_approval_repository.go:46)
- any order list query filtering on header `data_status`

### Validation and source-related logic to re-check
- [`ValidateOrder()`](../sales/service/validate_order_service.go:51)
- [`MapDataSourceToSource()`](../sales/service/order_service.go:28)
- source-specific branches in [`orderServiceImpl.Update()`](../sales/service/order_service.go:175)

## Recommendation Draft

### Short-term safest option

**Option B** is the safest operationally **unless** business explicitly confirms that web orders must become `Processed` immediately even when over limit and restricted.

Reason:
- It preserves current approval behavior
- It avoids accidentally bypassing approval
- It avoids breaking review queues with a narrow header-status hotfix

However, if the ticket must be closed with visible behavior change, short-term implementation should be:

- only proceed with **Option A** after explicit confirmation that approval should no longer depend on header `Need Review` for this case
- otherwise do not force `Processed`

### Long-term most correct solution

**Option C** is the most correct architectural solution.

Reason:
- Current model overloads one field for two concerns
- Lifecycle state and approval state should not share one code path
- This design is the most robust against recurring defects like this ticket

## Questions That Must Be Confirmed Before Coding

1. For over-credit-limit orders with outlet action `Restricted`, should the order:
   - remain blocked in approval state, or
   - still be considered `Processed` at header level?

2. Is `Need Review` intended to mean:
   - lifecycle status of the order, or
   - approval workflow status only?

3. If `sls.order.data_status` must be `2`, where should pending approval be represented?
   - approval tables only
   - new dedicated status field
   - derived response field only

4. Which downstream processes currently use `data_status` as gate?
   - order list screens
   - pick/pack/delivery flow
   - invoice generation
   - approval dashboards
   - reports/export jobs

5. Is web behavior supposed to differ from mobile behavior, or should both follow the same lifecycle + approval contract?

## Proposed Execution Plan

1. Confirm business meaning of `Need Review` vs `Processed`
2. Identify all consumers of `sls.order.data_status`
3. Decide target option A, B, or C with PM/BA/QA
4. If option A chosen, define approval-flow mitigation before coding
5. If option C chosen, design approval-state separation contract first
6. Switch to implementation mode only after decision approval
