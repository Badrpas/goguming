package mixins

import "fmt"

type KDA struct {
	KillCount, DeathCount, AttacksConnectedCount uint32
}

func (kda *KDA) ToString() string {
	return fmt.Sprintf("%d/%d/%d", kda.KillCount, kda.DeathCount, kda.AttacksConnectedCount)
}
