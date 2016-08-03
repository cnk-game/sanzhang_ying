package user

const (
	VipLevel_1_Diamond = 29
	VipLevel_2_Diamond = 49
	VipLevel_3_Diamond = 99
	VipLevel_4_Diamond = 299
	VipLevel_5_Diamond = 699
)

func GetVipLevel(diamond int) int {
	if diamond >= VipLevel_1_Diamond && diamond < VipLevel_2_Diamond {
		return 1
	} else if diamond >= VipLevel_2_Diamond && diamond < VipLevel_3_Diamond {
		return 2
	} else if diamond >= VipLevel_3_Diamond && diamond < VipLevel_4_Diamond {
		return 3
	} else if diamond >= VipLevel_4_Diamond && diamond < VipLevel_5_Diamond {
		return 4
	} else if diamond >= VipLevel_5_Diamond {
		return 5
	}

	return 0
}
