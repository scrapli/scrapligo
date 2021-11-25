package netconf

import (
	"encoding/xml"
	"errors"
)

// ErrUnknownFilterType error for when user provides an unknown filter type.
var ErrUnknownFilterType = errors.New("unknown filter type")

// ErrDefaultsType error for when user provides an unknown default type.
var ErrDefaultsType = errors.New("unknown defaults type")

const (
	// FilterSubtreeType constant for filter type of subtree.
	FilterSubtreeType = "subtree"
	// FilterXpathType constant for filter type of xpath.
	FilterXpathType = "xpath"
	// DefaultsReportAllType constant for default type of "report-all".
	DefaultsReportAllType = "report-all"
	// DefaultsTrimType constant for default type of "trim".
	DefaultsTrimType = "trim"
	// DefaultsExplicitType constant for default type of "explicit".
	DefaultsExplicitType = "explicit"
	// DefaultsReportAllTaggedType constant for default type of "report-all-tagged".
	DefaultsReportAllTaggedType = "report-all-tagged"
	// DefaultsNamespace constant for defaults namespace.
	DefaultsNamespace = "urn:ietf:params:xml:ns:yang:ietf-netconf-with-defaults"
)

// Message struct representing the base rpc message payload.
type Message struct {
	XMLName   xml.Name    `xml:"rpc"`
	Namespace string      `xml:"xmlns,attr"`
	MessageID int         `xml:"message-id,attr"`
	Payload   interface{} `xml:",innerxml"`
}

// BuildPayload build an XML payload to send to netconf device.
func (d *Driver) BuildPayload(payload interface{}) *Message {
	baseElem := &Message{
		XMLName:   xml.Name{},
		Namespace: "urn:ietf:params:xml:ns:netconf:base:1.0",
		MessageID: d.messageID,
		Payload:   payload,
	}

	d.messageID++

	return baseElem
}

// BuildRPCElem creates an element for a rpc operation.
func (d *Driver) BuildRPCElem(
	filter string,
) (*Message, error) {
	netconfInput := d.BuildPayload(filter)

	return netconfInput, nil
}

// datastore source

// SourceElement struct representing the individual source message element.
type SourceElement struct {
	XMLName xml.Name
}

// Source struct representing the parent source message element.
type Source struct {
	XMLName xml.Name       `xml:"source"`
	Source  *SourceElement `xml:""`
}

// BuildSourceElem creates the "source" (for get-config for example) element of a netconf payload.
func (d *Driver) BuildSourceElem(source string) *Source {
	sourceElem := &Source{
		XMLName: xml.Name{},
		Source:  &SourceElement{XMLName: xml.Name{Local: source}},
	}

	return sourceElem
}

// datastore target

// TargetElement struct representing the individual target message element.
type TargetElement struct {
	XMLName xml.Name
}

// Target struct representing the parent target message element.
type Target struct {
	XMLName xml.Name       `xml:"target"`
	Source  *TargetElement `xml:""`
}

// BuildTargetElem creates the "target" (for edit-config for example) element of a netconf payload.
func (d *Driver) BuildTargetElem(target string) *Target {
	targetElem := &Target{
		XMLName: xml.Name{},
		Source:  &TargetElement{XMLName: xml.Name{Local: target}},
	}

	return targetElem
}

// filter

// Filter struct representing the parent filter message element.
type Filter struct {
	XMLName xml.Name `xml:"filter"`
	Type    string   `xml:"type,attr"`
	Select  string   `xml:"select,attr,omitempty"`
	Payload string   `xml:",innerxml"`
}

// BuildFilterElem creates the "filter" element of a netconf payload.
func (d *Driver) BuildFilterElem(filter, filterType string) (*Filter, error) {
	if filter == "" || filterType == "" {
		return nil, nil
	}

	if filterType == FilterSubtreeType {
		return &Filter{
			XMLName: xml.Name{},
			Type:    filterType,
			Select:  "",
			Payload: filter,
		}, nil
	} else if filterType == FilterXpathType {
		return &Filter{
			XMLName: xml.Name{},
			Type:    filterType,
			Select:  filter,
		}, nil
	}

	return nil, ErrUnknownFilterType
}

// with-defaults

// DefaultType struct representing the parent default/with-defaults message element.
type DefaultType struct {
	XMLName   xml.Name `xml:"with-defaults"`
	Namespace string   `xml:"xmlns,attr"`
	Type      string   `xml:",innerxml"`
}

// BuildDefaultsElem creates the "default" element of a netconf payload.
func (d *Driver) BuildDefaultsElem(defaultsType string) (*DefaultType, error) {
	if defaultsType == "" {
		return nil, nil
	}

	if defaultsType != DefaultsReportAllType &&
		defaultsType != DefaultsTrimType &&
		defaultsType != DefaultsExplicitType &&
		defaultsType != DefaultsReportAllTaggedType {
		return nil, ErrDefaultsType
	}

	return &DefaultType{
		XMLName:   xml.Name{},
		Namespace: DefaultsNamespace,
		Type:      defaultsType,
	}, nil
}

// get

// Get struct representing the get message element.
type Get struct {
	XMLName xml.Name `xml:"get"`
	Source  *Source  `xml:""`
	Filter  *Filter  `xml:""`
}

// BuildGetElem creates a get element for a get operation.
func (d *Driver) BuildGetElem(
	filter, filterType string,
) (*Message, error) {
	filterElem, err := d.BuildFilterElem(filter, filterType)
	if err != nil {
		return nil, err
	}

	getElem := &Get{
		XMLName: xml.Name{},
		Filter:  filterElem,
	}

	netconfInput := d.BuildPayload(getElem)

	return netconfInput, nil
}

// get-config

// GetConfig struct representing the get-config message element.
type GetConfig struct {
	XMLName  xml.Name     `xml:"get-config"`
	Source   *Source      `xml:""`
	Filter   *Filter      `xml:""`
	Defaults *DefaultType `xml:""`
}

// BuildGetConfigElem creates a get-config element for a get operation.
func (d *Driver) BuildGetConfigElem(
	source, filter, filterType, defaultType string,
) (*Message, error) {
	filterElem, err := d.BuildFilterElem(filter, filterType)
	if err != nil {
		return nil, err
	}

	defaultsElem, err := d.BuildDefaultsElem(defaultType)
	if err != nil {
		return nil, err
	}

	getConfigElem := &GetConfig{
		XMLName:  xml.Name{},
		Source:   d.BuildSourceElem(source),
		Filter:   filterElem,
		Defaults: defaultsElem,
	}

	netconfInput := d.BuildPayload(getConfigElem)

	return netconfInput, nil
}

// edit-config

// EditConfig struct representing the edit-config message element.
type EditConfig struct {
	XMLName xml.Name `xml:"edit-config"`
	Target  *Target  `xml:""`
	Payload string   `xml:",innerxml"`
}

// BuildEditConfigElem creates a edit-config element for a get operation.
func (d *Driver) BuildEditConfigElem(
	config, target string,
) *Message {
	editConfigElem := &EditConfig{
		XMLName: xml.Name{},
		Target:  d.BuildTargetElem(target),
		Payload: config,
	}

	netconfInput := d.BuildPayload(editConfigElem)

	return netconfInput
}

// delete-config

// DeleteConfig struct representing the delete-config message element.
type DeleteConfig struct {
	XMLName xml.Name `xml:"delete-config"`
	Target  *Target  `xml:""`
}

// BuildDeleteConfigElem creates a delete-config element for a get operation.
func (d *Driver) BuildDeleteConfigElem(
	target string,
) *Message {
	deleteConfigElem := &DeleteConfig{
		XMLName: xml.Name{},
		Target:  d.BuildTargetElem(target),
	}

	netconfInput := d.BuildPayload(deleteConfigElem)

	return netconfInput
}

// copy-config

// CopyConfig struct representing the copy-config message element.
type CopyConfig struct {
	XMLName xml.Name `xml:"copy-config"`
	Source  *Source  `xml:""`
	Target  *Target  `xml:""`
}

// BuildCopyConfigElem creates a copy-config element for a copy-config operation.
func (d *Driver) BuildCopyConfigElem(
	source,
	target string,
) *Message {
	copyConfigElem := &CopyConfig{
		XMLName: xml.Name{},
		Source:  d.BuildSourceElem(source),
		Target:  d.BuildTargetElem(target),
	}

	netconfInput := d.BuildPayload(copyConfigElem)

	return netconfInput
}

// commit

// Commit struct representing the commit message element.
type Commit struct {
	XMLName xml.Name `xml:"commit"`
}

// BuildCommitElem creates a commit element for a get operation.
func (d *Driver) BuildCommitElem() *Message {
	commitElem := &Commit{
		XMLName: xml.Name{},
	}

	netconfInput := d.BuildPayload(commitElem)

	return netconfInput
}

// discard

// Discard struct representing the discard message element.
type Discard struct {
	XMLName xml.Name `xml:"discard-changes"`
}

// BuildDiscardElem creates a discard element for a get operation.
func (d *Driver) BuildDiscardElem() *Message {
	discardElem := &Discard{
		XMLName: xml.Name{},
	}

	netconfInput := d.BuildPayload(discardElem)

	return netconfInput
}

// lock

// Lock struct representing the lock message element.
type Lock struct {
	XMLName xml.Name `xml:"lock"`
	Target  *Target  `xml:""`
}

// BuildLockElem creates a lock element for a get operation.
func (d *Driver) BuildLockElem(target string) *Message {
	lockElem := &Lock{
		XMLName: xml.Name{},
		Target:  d.BuildTargetElem(target),
	}

	netconfInput := d.BuildPayload(lockElem)

	return netconfInput
}

// unlock

// Unlock struct representing the unlock message element.
type Unlock struct {
	XMLName xml.Name `xml:"unlock"`
	Target  *Target  `xml:""`
}

// BuildUnlockElem creates a unlock element for a get operation.
func (d *Driver) BuildUnlockElem(target string) *Message {
	unlockElem := &Unlock{
		XMLName: xml.Name{},
		Target:  d.BuildTargetElem(target),
	}

	netconfInput := d.BuildPayload(unlockElem)

	return netconfInput
}

// validate

// Validate struct representing the validate message element.
type Validate struct {
	XMLName xml.Name `xml:"validate"`
	Source  *Source  `xml:""`
}

// BuildValidateElem creates a validate element for a get operation.
func (d *Driver) BuildValidateElem(source string) *Message {
	validateElem := &Validate{
		XMLName: xml.Name{},
		Source:  d.BuildSourceElem(source),
	}

	netconfInput := d.BuildPayload(validateElem)

	return netconfInput
}
