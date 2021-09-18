package network

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/scrapli/scrapligo/logging"

	"github.com/scrapli/scrapligo/channel"
)

// ErrFailedToAcquirePriv raised when unable to acquire requested privilege level.
var ErrFailedToAcquirePriv = errors.New("failed to acquire requested privilege level")

func (d *Driver) buildPrivGraph() {
	d.privGraph = map[string]map[string]bool{}

	for _, privLevel := range d.PrivilegeLevels {
		privLevel.PatternRe = regexp.MustCompile(privLevel.Pattern)
		d.privGraph[privLevel.Name] = map[string]bool{}

		if privLevel.PreviousPriv != "" {
			d.privGraph[privLevel.Name][privLevel.PreviousPriv] = true
		}
	}

	for higherPrivLevel, privLevelList := range d.privGraph {
		for privLevel := range privLevelList {
			d.privGraph[privLevel][higherPrivLevel] = true
		}
	}
}

// UpdatePrivilegeLevels convenience method to build priv graph and create a new joined comms prompt
// pattern.
func (d *Driver) UpdatePrivilegeLevels() {
	d.buildPrivGraph()
	d.generateJoinedCommsPromptPattern()
}

func (d *Driver) escalate(escalatePriv string) error {
	var err error

	if !d.PrivilegeLevels[escalatePriv].EscalateAuth {
		_, err = d.Channel.SendInput(d.PrivilegeLevels[escalatePriv].Escalate, false, false, -1)
	} else {
		events := []*channel.SendInteractiveEvent{
			{
				ChannelInput:    d.PrivilegeLevels[escalatePriv].Escalate,
				ChannelResponse: d.PrivilegeLevels[escalatePriv].EscalatePrompt,
				HideInput:       false,
			},
			{
				ChannelInput:    d.AuthSecondary,
				ChannelResponse: d.PrivilegeLevels[escalatePriv].Pattern,
				HideInput:       true,
			},
		}
		_, err = d.Channel.SendInteractive(events, []string{d.PrivilegeLevels[escalatePriv].Pattern}, -1)
	}

	return err
}

func (d *Driver) deescalate(currentPriv string) error {
	_, err := d.Channel.SendInput(d.PrivilegeLevels[currentPriv].Deescalate, false, false, -1)
	return err
}

func (d *Driver) determineCurrentPriv(currentPrompt string) ([]string, error) {
	var matchingPrivLevels []string

PrivLevel:
	for privName, privData := range d.PrivilegeLevels {
		for _, notContains := range privData.PatternNotContains {
			if strings.Contains(currentPrompt, notContains) {
				continue PrivLevel
			}
		}

		promptMatch := privData.PatternRe.MatchString(currentPrompt)

		if promptMatch {
			matchingPrivLevels = append(matchingPrivLevels, privName)
		}
	}

	if len(matchingPrivLevels) == 0 {
		logging.LogError(
			d.FormatLogMessage(
				"error",
				fmt.Sprintf(
					"could not determine privilege level from provided prompt: %s\n",
					currentPrompt,
				),
			),
		)

		return []string{}, ErrCouldNotDeterminePriv
	}

	logging.LogDebug(
		d.FormatLogMessage(
			"error",
			fmt.Sprintf("determined current privilege level is one of: %s\n", matchingPrivLevels),
		),
	)

	return matchingPrivLevels, nil
}

func strSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func (d *Driver) buildPrivChangeMap(
	currentPriv, desiredPriv string,
	privChangeMap *[]string,
) []string {
	var workingPrivChangeMap []string

	if privChangeMap != nil {
		workingPrivChangeMap = *privChangeMap
	}

	workingPrivChangeMap = append(workingPrivChangeMap, currentPriv)

	if currentPriv == desiredPriv {
		return workingPrivChangeMap
	}

	for privName := range d.privGraph[currentPriv] {
		if !strSliceContains(workingPrivChangeMap, privName) {
			updatedPrivChangeMap := d.buildPrivChangeMap(
				privName,
				desiredPriv,
				&workingPrivChangeMap,
			)
			if len(updatedPrivChangeMap) > 0 {
				return updatedPrivChangeMap
			}
		}
	}

	return []string{}
}

func (d *Driver) processAcquirePriv(
	desiredPriv, currentPrompt string,
) (privilegeAction, string, error) {
	logging.LogDebug(
		d.FormatLogMessage(
			"info",
			fmt.Sprintf("attempting to acquire '%s' privilege level", desiredPriv),
		),
	)

	possibleCurrentPrivs, err := d.determineCurrentPriv(currentPrompt)

	if err != nil {
		return noAction, "", err
	}

	var currentPriv string

	switch {
	case strSliceContains(possibleCurrentPrivs, d.CurrentPriv):
		currentPriv = d.PrivilegeLevels[d.CurrentPriv].Name
	case strSliceContains(possibleCurrentPrivs, desiredPriv):
		currentPriv = d.PrivilegeLevels[desiredPriv].Name
	default:
		currentPriv = possibleCurrentPrivs[0]
	}

	if currentPriv == desiredPriv {
		logging.LogDebug(
			d.FormatLogMessage(
				"debug",
				"determined current privilege level is target privilege level, no action needed",
			),
		)

		d.CurrentPriv = desiredPriv

		return noAction, currentPriv, nil
	}

	mapToDesiredPriv := d.buildPrivChangeMap(currentPriv, desiredPriv, nil)

	// at this point we basically dont *know* the privilege leve we are at (or we wont/cant after
	// we do an escalation or deescalation, so we reset to the dummy priv level
	d.CurrentPriv = "UNKNOWN"

	if d.PrivilegeLevels[mapToDesiredPriv[1]].PreviousPriv != currentPriv {
		logging.LogDebug(d.FormatLogMessage("debug", "determined privilege deescalation necessary"))

		return deescalateAction, currentPriv, nil
	}

	logging.LogDebug(d.FormatLogMessage("debug", "determined privilege escalation necessary"))

	return escalateAction, d.PrivilegeLevels[mapToDesiredPriv[1]].Name, nil
}

// AcquirePriv acquire a target privilege level.
func (d *Driver) AcquirePriv(desiredPriv string) error {
	logging.LogDebug(
		d.FormatLogMessage(
			"debug",
			fmt.Sprintf("attempting to acquire privilege level: %s\n", desiredPriv),
		),
	)

	if _, ok := d.PrivilegeLevels[desiredPriv]; !ok {
		return ErrInvalidDesiredPriv
	}

	privChangeCount := 0

	for {
		currentPrompt, err := d.GetPrompt()
		if err != nil {
			logging.LogError(
				d.FormatLogMessage(
					"error",
					fmt.Sprintf("failed fetching prompt, error: %s\n", err),
				),
			)

			return err
		}

		privAction, targetPriv, err := d.processAcquirePriv(
			desiredPriv,
			currentPrompt,
		)

		switch {
		case err != nil:
			return err
		case privAction == noAction:
			return nil
		case privAction == escalateAction:
			err = d.escalate(targetPriv)
			if err != nil {
				return err
			}
		case privAction == deescalateAction:
			err = d.deescalate(targetPriv)
			if err != nil {
				return err
			}
		}

		privChangeCount++

		if privChangeCount > len(d.PrivilegeLevels)*2 {
			return ErrFailedToAcquirePriv
		}
	}
}
