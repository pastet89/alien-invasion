package simulator

import "github.com/pastet89/alien-invasion/utils"

type configVars struct {
    aliensCount int
    maxTravelsPerAlien int
    fightFormatMsg string
    alienNamePrefix string
    allAliensAreTrappedMessage string
}

func (s *Simulator) setMainConfigVars() {
	s.configVars = configVars{}
	config := utils.GetConfig()
	s.configVars.maxTravelsPerAlien = utils.GetIntConfigVar(config, "game/max_travels_per_alien")
	s.configVars.fightFormatMsg = utils.GetStringConfigVar(config, "game/fight_format_message")
	s.configVars.alienNamePrefix = utils.GetStringConfigVar(config, "game/alien_name_prefix")
	s.configVars.allAliensAreTrappedMessage = utils.GetStringConfigVar(config, "game/all_aliens_are_trapped_message")
}