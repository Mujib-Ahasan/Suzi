package attacks

import (
	"fmt"
	"log/slog"

	rs "github.com/Mujib-Ahasan/Suzi/common"
)

func Run(opts Options) []rs.PlotC {
	var attackAll []rs.PlotC
	switch opts.Type {
	case "all":
		slog.Debug("All Attack at once! Starting... \n")
		attackAll = append(attackAll, rs.PlotC{Attack: "basic", Results: BasicAttack(opts)})
		attackAll = append(attackAll, rs.PlotC{Attack: "burst", Results: BurstAttack(opts)})
		attackAll = append(attackAll, rs.PlotC{Attack: "rampup", Results: RampUpAttack(opts, 1, 15)})
		attackAll = append(attackAll, rs.PlotC{Attack: "random", Results: RandomLoadAttack(opts)})
	case "basic":
		slog.Debug("Basic Attack! Starting... \n")
		attackAll = append(attackAll, rs.PlotC{Attack: "basic", Results: BasicAttack(opts)})
	case "burst":
		slog.Debug("Burst Attack! Starting... \n")
		attackAll = append(attackAll, rs.PlotC{Attack: "burst", Results: BurstAttack(opts)})
	case "rampup":
		slog.Debug("Rampup Attack! Starting... \n")
		attackAll = append(attackAll, rs.PlotC{Attack: "rampup", Results: RampUpAttack(opts, 1, 15)})
	case "random":
		slog.Debug("Random Attack! Starting... \n")
		attackAll = append(attackAll, rs.PlotC{Attack: "random", Results: RandomLoadAttack(opts)})
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
