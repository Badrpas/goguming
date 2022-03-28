package foight

func GetUnitFromEntity(entity *Entity) (*Unit, bool) {
	unit, ok := entity.Holder.(*Unit)
	if ok {
		return unit, true
	}

	player, ok := entity.Holder.(*Player)
	if ok {
		return player.Unit, true
	}

	return nil, false
}
