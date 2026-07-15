# Promotion v2 Requirements Document

## Overview

This document outlines the enhancement of the existing promotion module in Scylla X. The enhancement aims to improve the functionality of the promotion module to meet user requirements, including:

1. **Promotion Types**
2. **Claim Realization**
3. **Promotion Criteria**
4. **Promotion Coverage**

## Reference Documents

- Promotion Calculation Flow
- New Promotion Status Flow
- Figma Flow Step Creation Promo
- Figma Epic: Promotion

## User Stories

### US-001: List Promotions
**As a promotion user, I can view promotion status**

**Workflow:**
1. User navigates to Promo & Discount → Promotion page
2. System displays promotions related to the distributor and promotions created by that distributor
3. If promotions are created by Principal/HO users and not yet approved, the Promotion ID will not appear to covered distributors
4. Promotions will not appear for distributors not registered as promotion participants

**Fields:**
- **Promotion ID**: Displays the ID of promotions applicable to the distributor or created by the distributor
- **Description**: Displays the description of promotions applicable to the distributor or created by the distributor
- **Effective Date**: Contains information about promotion validity dates. When promotions haven't entered the validity date range, they won't affect order calculations

### US-002: Initial Promotion Setting
**As a user, I can create initial promotion settings to determine promotion type and claim type**

**Workflow:**
1. User navigates to Promo & Discount → Promotion, then clicks + Add Promotion to create a new promotion
2. System displays Initial Promotion Setting page, then user selects promotion type

**Fields:**
- **Promotion Type**: Radio button dropdown containing promotion type options (Slab/Strata). User can only select one promotion type for one Promotion ID
- **Promotion Creation Type**: Radio button dropdown containing creation type options (New/Replacement). User can only select one type. If user selects Replacement type, the Existing Promotion ID field becomes enabled
- **Existing Promotion ID**: Dropdown containing promotion IDs available in the Promotion list menu. This field is only enabled if user selects replacement promotion type
- **Promotion ID**: Promotion ID input with maximum 50 alphanumeric characters with special characters
- **Promotion Description**: Promotion description input with maximum 100 alphanumeric characters with special characters
- **Budget Reference Toggle**: If budget reference is active (Yes), then Budget Reference and Budget Amount fields become enabled. If inactive (No), these fields become disabled
- **Budget Reference**: Field that becomes active if user enables Budget Reference toggle, containing radio buttons with Unlimited and Limited options
- **Budget Amount**: Field that becomes active if user selects Limited as budget reference, allowing numeric input with thousand separator format
- **Budget Control Level**: Optional dropdown radio button containing Region, Area, Distributor, Salesman
- **Execution Level**: Optional dropdown radio button containing Region, Area, Distributor, Salesman
- **Effective Date**: Field containing start-end date range for the created promotion
- **Claimable**: Radio button dropdown containing Yes/No options
- **Claim Type**: Radio button dropdown containing Full/Partial options
- **Claim Realization (%)**: Field that becomes enabled if user selects partial claim type, allowing numeric percentage input
- **Claim Will Start After**: Field that becomes enabled if user selects Claim Type Yes, allowing numeric input for days
- **Maximal Invoice Per Outlet**: Optional setting regarding limits on number of invoices eligible for promotion within one Promotion ID per outlet code
- **Maximal Total Reward Per Outlet**: Optional setting regarding limits on reward amount within one outlet, can be selected based on Amount/Qty

### US-003: Slab Promotion Type
**As a user, I can create promotions with Slab type so that within one Promotion ID I can create multiple promotion criteria**

**Workflow:**
1. If user selects Promotion Type Slab, system displays promotion criteria page in step 2

**Fields:**
- **Multiplied**: Dropdown field containing Yes/No radio buttons. After adding Slab, the dropdown becomes disabled
- **Slab Description**: Field allowing user to describe the slab with 50 alphanumeric character limit
- **Slab Rules**: Section containing basic promotion calculation requirements. User can only select one rule type: Quantity or Value
- **Rewards**: Section containing promotion reward options
- **List of Slabs**: Summary of slabs already added to the Promotion ID

### US-004: Strata Promotion Type
**As a user, I can create promotions with Strata type so users can see detailed contribution calculations for each promotion stratum**

**Workflow:**
1. If user selects Promotion Type Strata, system displays promotion criteria page

**Fields:**
- **Sequentially Calculated**: Dropdown field containing Yes/No radio buttons. After adding one stratum, this field becomes disabled
- **Claimable**: Dropdown field containing Yes/No radio buttons
- **Strata Description**: Field allowing user to describe the stratum with 50 alphanumeric character limit
- **Strata Rules**: Section containing basic promotion calculation requirements. User can only select one rule type: Quantity or Value
- **Strata Reward**: Section containing promotion reward options
- **List of Strata**: Summary of strata already added to the Promotion ID

### US-005: Promotion Criteria
**As a user, I want to set up product criteria requirements for promotions**

**Fields:**
- **Minimum SKU**: Field for determining the number of product type combinations that must exist in promotion criteria to be eligible for promotion
- **Principal Filter**: Dropdown for filtering product data based on supplier master
- **Category Filter**: Dropdown for filtering product data based on Product Category
- **Brand Filter**: Dropdown for filtering product data based on Product Brand
- **Product**: Dropdown containing product master according to selected filter data
- **Mandatory**: Radio button dropdown containing Yes/No options
- **Minimum Buy Setup**: Section for determining minimum purchase requirements from selected promotion criteria products

### US-006: Reward Product Setup
**If user selects product reward type, after completing step 3 (promotion product criteria) and clicking next, system displays Reward Product Setup page**

**Fields:**
- **Principal Filter**: Dropdown for filtering product data based on supplier master
- **Category Filter**: Dropdown for filtering product data based on Product Category
- **Brand Filter**: Dropdown for filtering product data based on Product Brand
- **Product**: Multiple select dropdown for choosing products that will become bonus products
- **List of Reward**: Summary of products selected as reward products

### US-007: Coverage
**If I login as a Principal level user, system displays Coverage page (step 6 if reward promo is product/step 5 if reward promo is other than product)**

**Coverage Types:**
- **National**: User can directly click next without clicking view. If user selects National coverage, all distributors with Active status in Distributor master are eligible
- **By Distributor**: If user selects by distributor coverage type, Region and Area fields become enabled and user can select distributor data based on Region and/or Area filters

### US-008: Outlet Criteria
**This step functions to determine which outlets are eligible for promotion. User can determine by directly selecting outlets or by Outlet Attribute**

**Fields:**
- **Select By Outlet**: Dropdown containing list of outlets with active status in that distributor
- **Outlet Class**: Contains data from setup parameter → Classification in setup parameter web
- **Outlet Group**: Contains data from setup parameter → Outlet Group in setup parameter web
- **Outlet Type**: Contains data from setup parameter → Classification in setup parameter web
- **Sales Team**: Contains data from master → salesman → Sales Team

### US-009: Summary
**After user completes all promotion setup, system displays Preview page containing summary of promotions set up in steps 1-7**

**Sections:**
- **Initial Promotion Setting**: Displays fields selected in Step 1 Initial Promotion Setting
- **Promotion Criteria**: Displays number of Slabs/Strata in the Promotion ID
- **Product Criteria**: Displays Minimum Buy SKU, Product Criteria, Minimum Buy in the Promotion ID
- **Reward Product Setup**: Displays list of products selected as reward products
- **Outlet Criteria Setup**: Displays list of outlet attributes and salesman attributes selected as promotion criteria
- **Coverage**: Displays list of distributors eligible for the promotion

### US-010: View Detail
**User navigates to Promotion list page, then clicks action button and selects view detail**

**Fields:**
- **Promotion Type**: Displays promotion type selected during creation
- **Promotion Creation Type**: Displays creation type selected during creation
- **Promotion ID**: Displays Promotion ID input during creation
- **Promotion Description**: Displays promotion description input during creation
- **Budget Reference**: Displays budget setting selection during creation
- **Budget**: Displays budget input
- **Budget Realization**: Displays total promotion value related to the Promotion ID from generated invoices
- **Remaining Budget**: Displays remaining budget that can be used for the Promotion ID
- **Budget Control Level**: Displays budget control level selected during creation
- **Budget Execution Level**: Displays budget execution level selected during creation
- **Effective Date**: Displays promotion validity date determined during creation
- **Claimable**: Displays claimable promotion selection during creation
- **Claim Type**: Displays claim type selection during creation
- **Claim Realization**: Displays realization percentage during creation
- **Claim Will Start After**: Displays claim will start after input during creation
- **Maximal Invoice Per Outlet**: Displays maximum invoice number input during creation
- **Maximal Total Reward Per Outlet**: Displays maximum reward amount input during creation

### US-011: Edit Promotion
**Edit promotion can only be done by users with the same level. If promotion is created by Principal level user, only Principal level users can edit**

**Status and Access:**
- **Draft**: Can be edited by same level users
- **Submit**: Edit button becomes disabled
- **Approved**: Can change status from Approved to Active, change effective date, add notes
- **Rejected**: Edit button becomes disabled
- **Active**: Can change status from Active to Inactive, change effective date, add notes
- **Inactive**: Can change status from Inactive to Closed, effective date disabled, add notes
- **Closed**: Edit button becomes disabled

### US-012: Duplicate Promotion
**Duplicate promotion is an action to duplicate already created promotions. When user performs duplicate, system creates a copy of selected promotion with sequence number [-001] added to the end of Promotion ID**

**Workflow:**
1. User navigates to promotion list page, then clicks duplicate action
2. System displays toast "promotion successfully duplicate" and creates promotion with sequence [-001] added to the end of Promotion ID
3. Duplicated promotion will have Draft status

## UI/UX Design

Figma Design: https://www.figma.com/design/V4qP9jfOg5oiPXEw2R0BRH/Scylla-X---2025?node-id=9137-105919&t=xVi8qtABgtZn8erB-1

## Technical Implementation

The promotion system is implemented using Go with the following key components:

- **Entity Layer**: Defines data structures and validation rules
- **Controller Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Manages data persistence
- **Model Layer**: Database models and relationships

### Key Features

1. **Flexible Promotion Types**: Support for both Slab and Strata promotion types
2. **Budget Management**: Unlimited and limited budget options with realization tracking
3. **Multi-level Access Control**: Different access levels for Principal and Distributor users
4. **Comprehensive Criteria**: Product, outlet, and coverage criteria support
5. **Status Management**: Complete promotion lifecycle management
6. **Real-time Tracking**: Budget realization and remaining budget calculations

### Database Schema

The promotion system uses PostgreSQL with the following main tables:
- `promo.promotions` - Main promotion table
- `promo.promotion_criteria` - Product criteria
- `promo.promotion_reward_products` - Reward products
- `promo.promotion_coverage_distributors` - Coverage settings
- `promo.promotion_outlet_criteria` - Outlet criteria

### API Endpoints

- `GET /v2/promotions` - List promotions with filtering
- `POST /v2/promotions` - Create new promotion
- `GET /v2/promotions/{id}` - Get promotion details
- `PUT /v2/promotions/{id}` - Update promotion
- `DELETE /v2/promotions/{id}` - Delete promotion
- `POST /v2/promotions/{id}/duplicate` - Duplicate promotion
