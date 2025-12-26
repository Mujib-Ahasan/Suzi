package attacks

import (
	"fmt"

	rs "github.com/Mujib-Ahasan/Suzi/common"
)

func Run(opts Options, mail bool) []rs.PlotC {
	var attackAll []rs.PlotC
	if mail {
		attackAll = append(attackAll, rs.PlotC{Results: BasicAttack(opts)})
		attackAll = append(attackAll, rs.PlotC{Results: BurstAttack(opts)})
		attackAll = append(attackAll, rs.PlotC{Results: RampUpAttack(opts, 1, 15)})
		attackAll = append(attackAll, rs.PlotC{Results: RandomLoadAttack(opts)})
		return attackAll
	}

	switch opts.Type {
	case "all":
		attackAll = append(attackAll, rs.PlotC{Results: BasicAttack(opts)})
		attackAll = append(attackAll, rs.PlotC{Results: BurstAttack(opts)})
		attackAll = append(attackAll, rs.PlotC{Results: RampUpAttack(opts, 1, 15)})
		attackAll = append(attackAll, rs.PlotC{Results: RandomLoadAttack(opts)})
	case "basic":
		attackAll = append(attackAll, rs.PlotC{Results: BasicAttack(opts)})
	case "burst":
		attackAll = append(attackAll, rs.PlotC{Results: BurstAttack(opts)})
	case "rampup":
		attackAll = append(attackAll, rs.PlotC{Results: RampUpAttack(opts, 1, 15)})
	case "random":
		attackAll = append(attackAll, rs.PlotC{Results: RandomLoadAttack(opts)})
	default:
		return []rs.PlotC{
			{Results: rs.PResultIn{
				Err: fmt.Errorf("unknown attack type: %s", opts.Type),
			},
			},
		}
	}
	return attackAll
}
