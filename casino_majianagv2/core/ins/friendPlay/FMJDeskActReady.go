package friendPlay

func (d *FMJDesk) Ready(userId uint32) error {
	d.SkeletonMJDesk.Ready(userId)
	//如果人数还是不够，就需要在计划增加机器人
	d.begin()
	return nil
}
