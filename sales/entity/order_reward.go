package entity

type CreateOrderRewardBody struct {
	CustId        string `json:"cust_id"`
	RoNo          string `json:"ro_no"`
	OrderRewardId int    `json:"order_reward_id"`
	ReffId        string `json:"reff_id"`
	SlabId        int    `json:"slab_id"`
	SlabDesc      string `json:"slab_desc"`
	RewardTypeId  int    `json:"reward_type_id"`
}

type OrderRewardResponse struct {
	CustId         string `json:"cust_id"`
	RoNo           string `json:"ro_no"`
	OrderRewardId  int    `json:"order_reward_id"`
	ReffId         string `json:"reff_id"`
	ReffName       string `json:"reff_name"`
	SlabDesc       string `json:"slab_desc"`
	RewardTypeId   int    `json:"reward_type_id"`
	RewardTypeName string `json:"reward_type_name"`
}

var rewardTypeName = map[int]string{
	1: "Promo",
	2: "Discount",
}

func (reward OrderRewardResponse) GenerateRewardTypeName() string {
	return rewardTypeName[int(reward.RewardTypeId)]
}
