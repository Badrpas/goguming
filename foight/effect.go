package foight

type Effect struct {
	*Entity

	Target *Player

	Duration  int64
	AppliedAt int64

	OnApply func()
	OnCease func()
}
