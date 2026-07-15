# ConsultV2 Function - Detailed Documentation

## Overview

The `ConsultV2` function is the main function for performing promotion consultation version 2. This function processes promotion consultation requests and returns a list of valid promotions along with appropriate rewards based on outlet, salesman, product, slab, and strata criteria.

## Function Signature

```go
func (service *promotionServiceImpl) ConsultV2(req entity.ConsultPromoV2Req) (responses []entity.ConsultPromoResp, err error)
```

## Input Parameters

### ConsultPromoV2Req

```go
type ConsultPromoV2Req struct {
    CustID       string                     `json:"cust_id" validate:"required"`
    ParentCustID string                     `json:"parent_cust_id" validate:"required"`
    OrderDate    string                     `json:"order_date" validate:"required"`
    OutletID     int                        `json:"outlet_id" validate:"required"`
    SalesmanID   int                        `json:"salesman_id" validate:"required"`
    WhID         int                        `json:"wh_id" validate:"required"`
    Details      []ConsultPromoV2DetailsReq `json:"details" validate:"required"`
}
```

### ConsultPromoV2DetailsReq

```go
type ConsultPromoV2DetailsReq struct {
    ProID      int     `json:"pro_id"`
    Qty1       float64 `json:"qty1"`       // Quantity smallest unit
    Qty2       float64 `json:"qty2"`       // Quantity middle unit
    Qty3       float64 `json:"qty3"`       // Quantity largest unit
    ConvUnit2  float64 `json:"conv_unit2"` // Conversion unit 2
    ConvUnit3  float64 `json:"conv_unit3"` // Conversion unit 3
    SubTotal   int     `json:"sub_total"`   // Price subtotal
    SellPrice1 float64 `json:"sell_price1"`
    SellPrice2 float64 `json:"sell_price2"`
    SellPrice3 float64 `json:"sell_price3"`
    SellPrice4 float64 `json:"sell_price4"`
    SellPrice5 float64 `json:"sell_price5"`
}
```

## Output

### ConsultPromoResp

```go
type ConsultPromoResp struct {
    PromoID       string                          `json:"promo_id"`
    PromoDesc     string                          `json:"promo_desc"`
    SlabID        string                          `json:"slab_id"`
    SlabDesc      string                          `json:"slab_desc"`
    SlabReward    float64                         `json:"slab_reward"`
    Products      []int                           `json:"products"`
    RewardPrice   []ConsultPromoRewardPriceResp   `json:"reward_price"`
    RewardProduct []ConsultPromoRewardProductResp `json:"reward_product"`
}
```

## Detailed Process Flow

### Phase 1: Initial Quantity Conversion

**Purpose**: Convert quantity from each product detail to standard units.

**Process**:
1. Loop through each `detail` in `req.Details`
2. For each detail, create `CreateConversionBody` with:
   - `CustId`: from `req.CustID`
   - `ProductId`: from `detail.ProID`
   - `Qty1`, `Qty2`, `Qty3`: from detail
3. Call `service.Conversion()` to convert quantity
4. Update `detail.Qty1`, `Qty2`, `Qty3` with conversion results
5. If an error occurs, log the error and continue to the next detail

**Note**: This conversion is important to ensure all quantities are in consistent units before validation.

---

### Phase 2: Outlet & Salesman Validation

**Purpose**: Validate that the provided outlet, salesman, and warehouse are valid and exist in the system.

**Process**:

1. **Outlet Validation**:
   - Call `FindOutletByID(outletID, custID)`
   - If not found, return error: "Outlet ID: {outletID} not found"

2. **Salesman Validation**:
   - Call `FindSalesmanByID(salesmanID, custID)`
   - If not found, return error: "Salesman ID: {salesmanID} not found"

3. **Warehouse Validation**:
   - Call `FindWarehouseByID(whID, custID)`
   - If not found, return error: "Warehouse ID: {whID} not found"

**Note**: All these validations must succeed before proceeding to the next phase.

---

### Phase 3: Build Attribute Validation Criteria

**Purpose**: Prepare attribute validation criteria (in v2, no need to build a map like v1, directly use outlet criteria).

**Process**: 
- No specific process in this phase, only a comment that v2 uses outlet criteria directly.

---

### Phase 4: Find & Validate Promotions by Outlet Criteria

**Purpose**: Find active promotions that match the outlet and salesman criteria.

**Process**:
1. Call `FindActivePromotionsByOutletCriteria(req, outlet, salesman)`
2. If no promotions are found, return empty `responses[]`
3. Extract `promoIDs` from all found promotions

**Output**: Array `promotions[]` containing matching active promotions.

---

### Phase 5 & 6: Find and Validate Product Criteria

**Purpose**: Validate whether products in the request meet the product criteria of each promotion.

**Process**:

1. **Find Product Criterias**:
   - Call `FindProductCriteriasByPromoIDs(promoIDs, parentCustID)`
   - Group results by `PromoID` into map `productCriteriaByPromo`

2. **Validate Each Promotion**:
   
   For each promotion:
   
   a. **If there is no product criteria**:
      - All products in request are eligible
      - Add all `req.Details` to `validatedPromoProductGroups[promoID]`
      - Sum all `SubTotal` to `subTotalValidatedPromoProductGroups[promoID]`
      - Set `isPromoProductCriteriaValid = true`
   
   b. **If there is product criteria**:
      
      **Loop for each productCriteria (including mandatory)**:
      - Set `productFound = false`
      - Loop through each `detail` in `req.Details`
      - If `detail.ProID == productCriteria.ProID`:
        
        **Calculate buyValue**:
        - If `MinBuyType == RuleTypeValue`: `buyValue = detail.SubTotal`
        - If `MinBuyType == RuleTypeQuantity`:
          - `UomTypeMiddle`: `buyValue = (Qty3 * ConvUnit3) + Qty2`
          - `UomTypeSmallest`: `buyValue = (Qty3 * ConvUnit3 * ConvUnit2) + (Qty2 * ConvUnit2) + Qty1`
          - Default: `buyValue = Qty3`
        
        **Calculate minBuyValue**:
        - From `MinBuyValue` or `MinBuyQty`
        
        **Validation**:
        - If `buyValue >= minBuyValue`:
          - Set `productFound = true`
          - Set `isPromoProductCriteriaValid = true`
          - Add detail to `validatedPromoProductGroups[promoID][detail.ProID]`
          - Add `SubTotal` to `subTotalValidatedPromoProductGroups[promoID]`
      
      - If `!productFound`, set `isPromoProductCriteriaValid = false` and break loop
      
      **Loop for non-mandatory products**:
      - Loop through productCriteria where `Mandatory = false`
      - If product is found and meets `buyValue >= minBuyValue`:
        - Add to `validatedPromoProductGroups`
        - Add `SubTotal` to `subTotalValidatedPromoProductGroups`
   
   c. **Final Validation**:
      - If `isPromoProductCriteriaValid == true` AND `len(validatedPromoProductGroups[promoID]) > 0`:
        - Add `promoID` to `validatedPromoList`
      - Otherwise:
        - Remove from `validatedPromoProductGroups` and `subTotalValidatedPromoProductGroups`

3. **Check Result**:
   - If `len(validatedPromoList) == 0`, return empty `responses[]`

**Output**: 
- `validatedPromoList`: List of valid promoIDs
- `validatedPromoProductGroups`: Map of valid products per promotion
- `subTotalValidatedPromoProductGroups`: Total subtotal per promotion

---

### Phase 7: Validate Slab Rules

**Purpose**: Validate whether the total purchase meets the slab (range) criteria of the promotion.

**Process**:

1. **Find Slabs**:
   - Call `FindSlabsByPromoIDs(validatedPromoList, parentCustID)`
   - Group slabs into map `promoSlabMap`

2. **Validate Each Slab**:
   
   For each slab:
   
   a. **Calculate slabRuleValue**:
      
      **If `RuleType == RuleTypeQuantity`**:
      - Loop through all details in `validatedPromoProductGroups[slab.PromoID]`
      - Calculate `buyValue` based on `RuleUom`:
        - `UomTypeMiddle`: `(Qty3 * ConvUnit3) + Qty2`
        - `UomTypeSmallest`: `(Qty3 * ConvUnit3 * ConvUnit2) + (Qty2 * ConvUnit2) + Qty1`
        - Default: `Qty3`
      - Sum all `buyValue` to get `slabRuleValue`
      
      **If `RuleType == RuleTypeValue`**:
      - Sum all `SubTotal` from details in `validatedPromoProductGroups[slab.PromoID]`
      - Result becomes `slabRuleValue`
   
   b. **Check SlabMultiplied Flag**:
      - Check parent promotion for `SlabMultiplied` flag
      - If `true`, set `isMultiplied = true`
   
   c. **Validate Range**:
      - If `isMultiplied == true` OR (`slabRuleValue >= RangeFrom` AND `slabRuleValue <= RangeTo`):
        - Add slab to `validatedSlabs[slab.PromoID]`

**Output**: `validatedSlabs` map containing valid slabs for each promotion.

---

### Phase 7b: Validate Strata

**Purpose**: Validate whether the total purchase meets the strata (tier) criteria of the promotion.

**Process**:

1. **Find Stratas**:
   - Call `FindStratasByPromoIDs(validatedPromoList, parentCustID)`
   - Sort stratas by `PromoID`, then `Ordinal` (ASC)

2. **Validate Each Strata**:
   
   For each `promoID` in `validatedPromoProductGroups`:
   
   - Loop through stratas (already sorted by ordinal)
   - If `strata.PromoID == promoID`:
     
     a. **Calculate ruleValue**:
        - Call `calculateStrataRuleValue(strata, details)`
        - This function calculates based on `RuleType`:
          - `RuleTypeQuantity`: Sum converted quantity based on `RuleUom`
          - `RuleTypeValue`: Sum `SubTotal`
     
     b. **Validate Range**:
        - If `ruleValue >= strata.RangeFrom` AND `ruleValue <= strata.RangeTo`:
          - Set `validatedStrata[promoID] = strata`
          - Break loop (take first matching strata because already sorted)

**Output**: `validatedStrata` map containing valid strata for each promotion.

**Note**: Strata is selected based on the lowest ordinal that meets the criteria.

---

### Phase 8: Calculate Rewards

**Purpose**: Calculate rewards for each valid promotion and build the response.

**Process**:

1. **Combine Validated Promotions**:
   - Combine `promoID` from `validatedSlabs` and `validatedStrata`
   - Result: `combinedPromoIDs` map

2. **Process Each Promotion**:
   
   For each `promoID` in `combinedPromoIDs`:
   
   a. **Determine Reward Type**:
      - Priority: Use reward from `strata` if exists
      - Fallback: Use reward from `slab` if exists
      - If neither exists, skip this promotion
      - Extract: `rewardType`, `rewardValue`, `rewardUom`
   
   b. **Find Promotion Details**:
      - Find promotion object from `promotions[]` by `promoID`
   
   c. **Create Base Response**:
      ```go
      response.PromoID = promoID
      response.PromoDesc = promo.PromoDesc
      if hasSlab {
          response.SlabID = slab.ID
          response.SlabDesc = *slab.Description
      }
      ```
   
   d. **Calculate SlabReward**:
      - If `rewardType == RewardTypePercentage`:
        - `SlabReward = Round((subTotal * rewardValue) / 100)`
      - If `rewardType == RewardTypeFixedValue`:
        - `SlabReward = rewardValue`
      - Default: `SlabReward = 0`
   
   e. **Process Products for Price Rewards**:
      
      For each `proID` in `validatedPromoProductGroups[promoID]`:
      
      - Add `proID` to `response.Products`
      
      - If `response.SlabReward > 0`:
        - Create `rewardPrice`:
          - `ProID = proID`
          - `SubTotal = detail.SubTotal`
          - `Reward = slabReward` (temporary)
          - `slabReward -= reward`
          - If `slabReward <= 0`, adjust `reward += slabReward`
          - `Total = SubTotal - Reward`
        - Add `rewardPrice` to `response.RewardPrice`
      
      - If `len(response.RewardProduct) > 0` OR `len(response.RewardPrice) > 0`:
        - Append `response` to `responses[]`
   
   f. **Handle Product Rewards** (if `rewardType == RewardTypeProduct`):
      
      - **Create Reward Context**:
        - Create `rewardCtx` with `PromoID` and `RewardUom` from strata or slab
      
      - **Get Reward Products from Stock**:
        - Call `GetAllRewardProductFromStockV2(req, rewardCtx)`
        - Get `rewards[]` (reward products available in stock)
      
      - **Calculate Multiplied Value** (if `SlabMultiplied == true`):
        - Calculate `slabRuleValue` (same as in Phase 7)
        - `multipliedValue = slabRuleValue / slab.RangeTo`
        - Default: `multipliedValue = 1`
      
      - **Calculate Total Quantity Reward**:
        - If `rewardValue > 0`: `totalQtyReward = rewardValue * multipliedValue`
        - Else if `multipliedValue > 0`: `totalQtyReward = multipliedValue`
      
      - **Distribute Reward Products**:
        
        If `totalQtyReward > 0`:
        
        Loop through each `reward` in `rewards[]`:
        
        - Calculate `qtyReward = min(totalQtyReward, reward.QtyStock)`
        - `totalQtyReward -= qtyReward`
        
        - **Convert Quantity**:
          - Create `rewardProductConversion` based on `RewardUom`:
            - `UomTypeSmallest`: set `Qty1 = qtyReward`
            - `UomTypeMiddle`: set `Qty2 = qtyReward`
            - Default: set `Qty3 = qtyReward`
          - Call `Conversion(rewardProductConversion, custID, parentCustID)`
          - Get conversion results `Qty1`, `Qty2`, `Qty3`
        
        - **Create Reward Product**:
          - `ProID = reward.ProID`
          - `Qty1`, `Qty2`, `Qty3` from conversion results
          - Add to `response.RewardProduct`
        
        - If `totalQtyReward <= 0`, break loop
      
      - If `len(response.RewardProduct) > 0` OR `len(response.RewardPrice) > 0`:
        - Append `response` to `responses[]`

3. **Return Results**:
   - Return `responses[]` containing all valid promotions with their rewards

---

## Helper Functions

### calculateStrataRuleValue

Helper function to calculate the rule value from strata based on rule type and UOM.

```go
func calculateStrataRuleValue(
    strata model.PromotionV2Strata,
    details map[int]*entity.ConsultPromoV2DetailsReq,
) float64
```

**Process**:
- If `RuleType == RuleTypeQuantity`:
  - Loop through all details
  - Convert quantity based on `RuleUom`:
    - `UomTypeSmallest`: `(Qty3 * ConvUnit3 * ConvUnit2) + (Qty2 * ConvUnit2) + Qty1`
    - `UomTypeMiddle`: `(Qty3 * ConvUnit3) + Qty2`
    - Default: `Qty3`
  - Sum all conversion results
- If `RuleType == RuleTypeValue`:
  - Sum all `SubTotal` from details

---

## Key Data Structures

### Internal Maps and Lists

1. **validatedPromoProductGroups**: `map[string]map[int]*ConsultPromoV2DetailsReq`
   - Key: `promoID`
   - Value: Map of valid products (key: `proID`, value: detail)

2. **subTotalValidatedPromoProductGroups**: `map[string]int64`
   - Key: `promoID`
   - Value: Total subtotal of valid products

3. **validatedPromoList**: `[]string`
   - List of valid promoIDs after product validation

4. **validatedSlabs**: `map[string]model.PromotionV2Slabs`
   - Key: `promoID`
   - Value: Valid slab

5. **validatedStrata**: `map[string]model.PromotionV2Strata`
   - Key: `promoID`
   - Value: Valid strata

6. **combinedPromoIDs**: `map[string]bool`
   - Combined promoIDs from validatedSlabs and validatedStrata

---

## Error Handling

1. **Phase 1 (Conversion)**: 
   - Error is logged, process continues to next detail
   - Does not stop the overall process

2. **Phase 2 (Validation)**:
   - If outlet/salesman/warehouse not found, return error immediately
   - Process is stopped

3. **Phase 4-7**:
   - If error occurs on repository calls, return error
   - If no valid promotions, return empty array (not an error)

---

## Business Rules

1. **Product Criteria**:
   - If there is no product criteria, all products are eligible
   - Mandatory products must exist and meet minimum buy requirement
   - Non-mandatory products can be added if they meet minimum buy requirement

2. **Slab Validation**:
   - Slab is valid if `slabRuleValue` is within range `[RangeFrom, RangeTo]`
   - Or if `SlabMultiplied == true` (no need to check range)

3. **Strata Validation**:
   - Strata is selected based on the lowest ordinal that meets the criteria
   - Only one strata per promotion is selected

4. **Reward Calculation**:
   - Priority: Strata reward > Slab reward
   - Price reward: Can be percentage or fixed value
   - Product reward: Can be multiplied if `SlabMultiplied == true`
   - Quantity reward is converted based on specified UOM

5. **Reward Distribution**:
   - Price reward is distributed to all eligible products
   - Product reward is taken from stock, with priority based on order in `rewards[]`
   - If stock is insufficient, take according to available stock

---

## Performance Considerations

1. **Database Queries**:
   - Phase 4-7 performs multiple queries to the database
   - Consider batch queries if possible

2. **Loop Optimization**:
   - Some nested loops can be optimized by using maps for lookup
   - Example: `productCriteriaByPromo` map for grouping

3. **Memory Usage**:
   - Several maps and lists are stored in memory
   - For requests with many promotions and products, consider memory usage

---

## Testing Scenarios

1. **No Promotions Found**: Return empty array
2. **No Valid Products**: Return empty array
3. **No Valid Slabs/Strata**: Promotions do not enter the response
4. **Multiple Promotions**: All valid promotions must enter the response
5. **Product Rewards with Stock**: Product rewards must match available stock
6. **SlabMultiplied**: Reward must be multiplied according to calculation
7. **Strata Selection**: Must select strata with the lowest valid ordinal

---

## Related Files

- **Service**: `/sales/service/promotion_service.go` (line 2258-2832)
- **Repository**: `/sales/repository/promotionV2_repository.go`
- **Entity**: `/sales/entity/promotionV2.go`
- **Model**: `/sales/model/promotionV2.go`

---

## Notes

- This function is version 2 of `ConsultPromotion`, with a more modular structure
- Main difference from v1: no need to build attribute validation map, directly use outlet criteria
- All quantities must be converted first before validation
- Rewards can be price (discount) or product (free item)
