package session

import (
	"regexp"
	"strings"

	"github.com/evilsocket/bettercap-ng/core"
)

const (
	PromptVariable = "$"
	DefaultPrompt  = "{by}{fw}{iface.cidr} {fb}> {iface.addr} {reset} {bold}» {reset}"
)

var PromptEffects = map[string]string{
	"{bold}":  core.BOLD,
	"{dim}":   core.DIM,
	"{r}":     core.RED,
	"{g}":     core.GREEN,
	"{b}":     core.BLUE,
	"{y}":     core.YELLOW,
	"{fb}":    core.FG_BLACK,
	"{fw}":    core.FG_WHITE,
	"{bdg}":   core.BG_DGRAY,
	"{br}":    core.BG_RED,
	"{bg}":    core.BG_GREEN,
	"{by}":    core.BG_YELLOW,
	"{blb}":   core.BG_LBLUE, // Ziggy this is for you <3
	"{reset}": core.RESET,
}

var PromptCallbacks = map[string]func(s *Session) string{
	"{iface.cidr}": func(s *Session) string {
		return s.Interface.CIDR()
	},
	"{iface.addr}": func(s *Session) string {
		return s.Interface.IpAddress
	},
	"{iface.mac}": func(s *Session) string {
		return s.Interface.HwAddress
	},
	"{gw.addr}": func(s *Session) string {
		return s.Gateway.IpAddress
	},
	"{gw.mac}": func(s *Session) string {
		return s.Gateway.HwAddress
	},
}

var envRe = regexp.MustCompile("{env\\.(.+)}")

type Prompt struct {
}

func NewPrompt() Prompt {
	return Prompt{}
}

func (p Prompt) Render(s *Session) string {
	found, prompt := s.Env.Get(PromptVariable)
	if found == false {
		prompt = DefaultPrompt
	}

	for tok, effect := range PromptEffects {
		prompt = strings.Replace(prompt, tok, effect, -1)
	}

	for tok, cb := range PromptCallbacks {
		prompt = strings.Replace(prompt, tok, cb(s), -1)
	}

	m := envRe.FindStringSubmatch(prompt)
	if len(m) == 2 {
		name := m[1]
		_, value := s.Env.Get(name)
		prompt = strings.Replace(prompt, m[0], value, -1)
	}

	// make sure an user error does not screw all terminal
	if strings.HasPrefix(prompt, core.RESET) == false {
		prompt += core.RESET
	}

	return prompt
}
