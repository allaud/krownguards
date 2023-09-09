package models

type Upgrade struct {
	GradeStep float64 // Value added by 1 upgrade step
	MaxGrade  int     // Max step of upgrade
	CurrGrade int     // Current step of upgrade
	Price     int     // Stone cost of upgrade
}

type KingUpgrades struct {
	HpRegGrade Upgrade //
	AtkGrade   Upgrade //
	MaxHpGrade Upgrade //
	//Treasury   int     // Amount of gold King earn
}
