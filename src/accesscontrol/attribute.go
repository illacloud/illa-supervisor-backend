package accesscontrol

import "github.com/illacloud/illa-supervisor-backend/src/model"

// default
const ANONYMOUS_AUTH_TOKEN = ""
const DEFAULT_TEAM_ID = 0
const DEFAULT_UNIT_ID = 0

// user status in team
const STATUS_OK = 1
const STATUS_PENDING = 2
const STATUS_SUSPEND = 3

// Attirbute Unit List
// Attirbute Unit List
const (
	UNIT_TYPE_TEAM                      = 1  // cloud team
	UNIT_TYPE_TEAM_MEMBER               = 2  // cloud team member
	UNIT_TYPE_USER                      = 3  // cloud user
	UNIT_TYPE_INVITE                    = 4  // cloud invite
	UNIT_TYPE_DOMAIN                    = 5  // cloud domain
	UNIT_TYPE_BILLING                   = 6  // cloud billing
	UNIT_TYPE_BUILDER_DASHBOARD         = 7  // builder dabshboard
	UNIT_TYPE_APP                       = 8  // builder app
	UNIT_TYPE_COMPONENTS                = 9  // builder components
	UNIT_TYPE_RESOURCE                  = 10 // resource resource
	UNIT_TYPE_ACTION                    = 11 // resource action
	UNIT_TYPE_TRANSFORMER               = 12 // resource transformer
	UNIT_TYPE_JOB                       = 13 // hub job
	UNIT_TYPE_TREE_STATES               = 14 // components tree states
	UNIT_TYPE_KV_STATES                 = 15 // components k-v states
	UNIT_TYPE_SET_STATES                = 16 // components set states
	UNIT_TYPE_PROMOTE_CODES             = 17 // promote codes
	UNIT_TYPE_PROMOTE_CODE_USAGES       = 18 // promote codes usage table
	UNIT_TYPE_ROLES                     = 19 // team roles table
	UNIT_TYPE_USER_ROLE_RELATIONS       = 20 // user role relation table
	UNIT_TYPE_UNIT_ROLE_RELATIONS       = 21 // unit role relation table
	UNIT_TYPE_COMPENSATING_TRANSACTIONS = 22 // compensating transactions
	UNIT_TYPE_TRANSACTION_SERIALS       = 23 // transaction serials
	UNIT_TYPE_CAPACITIES                = 24 // capacity
	UNIT_TYPE_DRIVE                     = 25 // drive
	UNIT_TYPE_PERIPHERAL_SERVICE        = 26 // Peripheral service, including sql generate, STMP etc.
)

// global invite permission config
// owner & admin -> can invite admin, editor, viewer
// editor 	     -> can invite editor, viewer
// viewer 	     -> can invite viewer
// map[nowUserRole]map[atrgetUserRole]attribute

// this config map target role to target invite role attribute
// e.g. you want invite model.USER_ROLE_ADMIN, so it's mapped attribute is ACTION_ACCESS_INVITE_ADMIN
var InviteRoleAttributeMap = map[int]int{
	model.USER_ROLE_OWNER: ACTION_ACCESS_INVITE_OWNER, model.USER_ROLE_ADMIN: ACTION_ACCESS_INVITE_ADMIN, model.USER_ROLE_EDITOR: ACTION_ACCESS_INVITE_EDITOR, model.USER_ROLE_VIEWER: ACTION_ACCESS_INVITE_VIEWER,
}

// this config map target role to target manage user role attribute
// e.g. you want modify a user to role model.USER_ROLE_EDITOR, so it's mapped attribute is ACTION_MANAGE_ROLE_TO_EDITOR
var ModifyRoleFromAttributeMap = map[int]int{
	model.USER_ROLE_OWNER: ACTION_MANAGE_ROLE_FROM_OWNER, model.USER_ROLE_ADMIN: ACTION_MANAGE_ROLE_FROM_ADMIN, model.USER_ROLE_EDITOR: ACTION_MANAGE_ROLE_FROM_EDITOR, model.USER_ROLE_VIEWER: ACTION_MANAGE_ROLE_FROM_VIEWER,
}
var MadifyRoleToAttributeMap = map[int]int{
	model.USER_ROLE_OWNER: ACTION_MANAGE_ROLE_TO_OWNER, model.USER_ROLE_ADMIN: ACTION_MANAGE_ROLE_TO_ADMIN, model.USER_ROLE_EDITOR: ACTION_MANAGE_ROLE_TO_EDITOR, model.USER_ROLE_VIEWER: ACTION_MANAGE_ROLE_TO_VIEWER,
}

const (
	ATTRIBUTE_CATEGORY_ACCESS  = 1
	ATTRIBUTE_CATEGORY_DELETE  = 2
	ATTRIBUTE_CATEGORY_MANAGE  = 3
	ATTRIBUTE_CATEGORY_SPECIAL = 4
)

// Attribute List
// action access
const (
	// Basic Attribute
	ACTION_ACCESS_VIEW = iota + 1 // 访问 Attribute
	// Invite Attribute
	ACTION_ACCESS_INVITE_BY_LINK  // invite team member by link
	ACTION_ACCESS_INVITE_BY_EMAIL // invite team member by email
	ACTION_ACCESS_INVITE_OWNER    // can invite team member as an owner
	ACTION_ACCESS_INVITE_ADMIN    // can invite team member as an admin
	ACTION_ACCESS_INVITE_EDITOR   // can invite team member as an editor
	ACTION_ACCESS_INVITE_VIEWER   // can invite team member as a viewer
)

// action manage
const (
	// Team Attribute
	ACTION_MANAGE_TEAM_NAME          = iota + 1 // rename Team Attribute
	ACTION_MANAGE_TEAM_ICON                     // update icon
	ACTION_MANAGE_TEAM_CONFIG                   // update team config
	ACTION_MANAGE_UPDATE_TEAM_DOMAIN            // update team domain

	// Team Member Attribute
	ACTION_MANAGE_REMOVE_MEMBER    // remove member from a team
	ACTION_MANAGE_ROLE             // manage role of team member
	ACTION_MANAGE_ROLE_FROM_OWNER  // modify team member role from owner ..
	ACTION_MANAGE_ROLE_FROM_ADMIN  // modify team member role from admin ..
	ACTION_MANAGE_ROLE_FROM_EDITOR // modify team member role from editor ..
	ACTION_MANAGE_ROLE_FROM_VIEWER // modify team member role from viewer ..
	ACTION_MANAGE_ROLE_TO_OWNER    // modify team member role to owner
	ACTION_MANAGE_ROLE_TO_ADMIN    // modify team member role to admin
	ACTION_MANAGE_ROLE_TO_EDITOR   // modify team member role to editor
	ACTION_MANAGE_ROLE_TO_VIEWER   // modify team member role to viewer

	// User Attribute
	ACTION_MANAGE_RENAME_USER        // rename
	ACTION_MANAGE_UPDATE_USER_AVATAR // update avatar

	// Invite Attribute
	ACTION_MANAGE_CONFIG_INVITE // config invite
	ACTION_MANAGE_INVITE_LINK   // config invite link, open, close and renew

	// Domain Attribute
	ACTION_MANAGE_TEAM_DOMAIN // update team domain
	ACTION_MANAGE_APP_DOMAIN  // update app domain

	// Billing Attribute
	ACTION_MANAGE_PAYMENT      // manage payment, including create, update, cancel subscribe and purchase
	ACTION_MANAGE_PAYMENT_INFO // manage team payment info, including get portal info band billing info

	// Dashboard Attribute
	ACTION_MANAGE_DASHBOARD_BROADCAST

	// App Attribute
	ACTION_MANAGE_CREATE_APP // create APP
	ACTION_MANAGE_EDIT_APP   // edit APP

	// Resource Attribute
	ACTION_MANAGE_CREATE_RESOURCE // create resource
	ACTION_MANAGE_EDIT_RESOURCE   // edit resource

	// Action Attribute
	ACTION_MANAGE_CREATE_ACTION  // create action
	ACTION_MANAGE_EDIT_ACTION    // edit action
	ACTION_MANAGE_PREVIEW_ACTION // preview action
	ACTION_MANAGE_RUN_ACTION     // run action

	// Drive Attribute
	ACTION_MANAGE_CREATE_FILE      // create file
	ACTION_MANAGE_EDIT_FILE        // edit file
	ACTION_MANAGE_CREATE_SHARELINK // create sharelink
)

// action delete
const (
	// Basic Attribute
	ACTION_DELETE = iota + 1 // delete Attribute

	// Domain Attribute
	ACTION_DELETE_TEAM_DOMAIN // delete Team Domain
	ACTION_DELETE_APP_DOMAIN  // delete App Domain

)

// action manage special (only owner and admin can access by default)
const (
	// Team Attribute
	ACTION_SPECIAL_EDITOR_AND_VIEWER_CAN_INVITE_BY_LINK_SW = iota + 1 // the "editor and viewer can invite" switch
	// Team Member Attribute
	ACTION_SPECIAL_TRANSFER_OWNER // transfer team owner to others
	// Invite Attribute
	ACTION_SPECIAL_INVITE_LINK_RENEW // renew the invite link
	// APP Attribute
	ACTION_SPECIAL_RELEASE_APP // release APP
	// SQL Generate
	ACTION_SPECIAL_GENERATE_SQL // generate sql
	// APP Snapshot
	ACTOIN_SPECIAL_TAKE_SNAPSHOT
	ACTOIN_SPECIAL_RECOVER_SNAPSHOT
)

// Attribute Config List
// Only define avaliable attribute here
// map[AttributeCategory][role][unitType][Attribute]status
var AttributeConfigList = map[int]map[int]map[int]map[int]bool{
	ATTRIBUTE_CATEGORY_ACCESS: {
		model.USER_ROLE_ANONYMOUS: {
			UNIT_TYPE_APP:    {ACTION_ACCESS_VIEW: true}, // only should for public app
			UNIT_TYPE_ACTION: {ACTION_ACCESS_VIEW: true}, // only should for public action
		},
		model.USER_ROLE_OWNER: {
			UNIT_TYPE_TEAM:              {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_USER:              {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_INVITE:            {ACTION_ACCESS_VIEW: true, ACTION_ACCESS_INVITE_BY_LINK: true, ACTION_ACCESS_INVITE_BY_EMAIL: true, ACTION_ACCESS_INVITE_ADMIN: true, ACTION_ACCESS_INVITE_EDITOR: true, ACTION_ACCESS_INVITE_VIEWER: true},
			UNIT_TYPE_DOMAIN:            {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_BILLING:           {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_APP:               {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_COMPONENTS:        {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_RESOURCE:          {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_ACTION:            {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_TRANSFORMER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_JOB:               {ACTION_ACCESS_VIEW: true},
		},
		model.USER_ROLE_ADMIN: {
			UNIT_TYPE_TEAM:              {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_USER:              {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_INVITE:            {ACTION_ACCESS_VIEW: true, ACTION_ACCESS_INVITE_BY_LINK: true, ACTION_ACCESS_INVITE_BY_EMAIL: true, ACTION_ACCESS_INVITE_ADMIN: true, ACTION_ACCESS_INVITE_EDITOR: true, ACTION_ACCESS_INVITE_VIEWER: true},
			UNIT_TYPE_DOMAIN:            {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_APP:               {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_COMPONENTS:        {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_RESOURCE:          {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_ACTION:            {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_TRANSFORMER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_JOB:               {ACTION_ACCESS_VIEW: true},
		},
		model.USER_ROLE_EDITOR: {
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_USER:              {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_INVITE:            {ACTION_ACCESS_VIEW: true, ACTION_ACCESS_INVITE_BY_LINK: true, ACTION_ACCESS_INVITE_BY_EMAIL: true, ACTION_ACCESS_INVITE_EDITOR: true, ACTION_ACCESS_INVITE_VIEWER: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_APP:               {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_COMPONENTS:        {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_RESOURCE:          {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_ACTION:            {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_TRANSFORMER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_JOB:               {ACTION_ACCESS_VIEW: true},
		},
		model.USER_ROLE_VIEWER: {
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_USER:              {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_INVITE:            {ACTION_ACCESS_VIEW: true, ACTION_ACCESS_INVITE_BY_LINK: true, ACTION_ACCESS_INVITE_BY_EMAIL: true, ACTION_ACCESS_INVITE_VIEWER: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_APP:               {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_COMPONENTS:        {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_RESOURCE:          {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_ACTION:            {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_TRANSFORMER:       {ACTION_ACCESS_VIEW: true},
			UNIT_TYPE_JOB:               {ACTION_ACCESS_VIEW: true},
		},
	},
	ATTRIBUTE_CATEGORY_DELETE: {
		model.USER_ROLE_OWNER: {
			UNIT_TYPE_TEAM:              {ACTION_DELETE: true},
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_DELETE: true},
			UNIT_TYPE_USER:              {ACTION_DELETE: true},
			UNIT_TYPE_INVITE:            {ACTION_DELETE: true},
			UNIT_TYPE_DOMAIN:            {ACTION_DELETE: true, ACTION_DELETE_TEAM_DOMAIN: true, ACTION_DELETE_APP_DOMAIN: true},
			UNIT_TYPE_BILLING:           {ACTION_DELETE: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_DELETE: true},
			UNIT_TYPE_APP:               {ACTION_DELETE: true},
			UNIT_TYPE_COMPONENTS:        {ACTION_DELETE: true},
			UNIT_TYPE_RESOURCE:          {ACTION_DELETE: true},
			UNIT_TYPE_ACTION:            {ACTION_DELETE: true},
			UNIT_TYPE_TRANSFORMER:       {ACTION_DELETE: true},
			UNIT_TYPE_JOB:               {ACTION_DELETE: true},
		},
		model.USER_ROLE_ADMIN: {
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_DELETE: true},
			UNIT_TYPE_USER:              {ACTION_DELETE: true},
			UNIT_TYPE_INVITE:            {ACTION_DELETE: true},
			UNIT_TYPE_DOMAIN:            {ACTION_DELETE: true, ACTION_DELETE_TEAM_DOMAIN: true, ACTION_DELETE_APP_DOMAIN: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_DELETE: true},
			UNIT_TYPE_APP:               {ACTION_DELETE: true},
			UNIT_TYPE_COMPONENTS:        {ACTION_DELETE: true},
			UNIT_TYPE_RESOURCE:          {ACTION_DELETE: true},
			UNIT_TYPE_ACTION:            {ACTION_DELETE: true},
			UNIT_TYPE_TRANSFORMER:       {ACTION_DELETE: true},
			UNIT_TYPE_JOB:               {ACTION_DELETE: true},
		},
		model.USER_ROLE_EDITOR: {
			UNIT_TYPE_TEAM_MEMBER: {ACTION_DELETE: true},
			UNIT_TYPE_USER:        {ACTION_DELETE: true},
			UNIT_TYPE_INVITE:      {ACTION_DELETE: true},
			UNIT_TYPE_APP:         {ACTION_DELETE: true},
			UNIT_TYPE_COMPONENTS:  {ACTION_DELETE: true},
			UNIT_TYPE_RESOURCE:    {ACTION_DELETE: true},
			UNIT_TYPE_ACTION:      {ACTION_DELETE: true},
			UNIT_TYPE_TRANSFORMER: {ACTION_DELETE: true},
			UNIT_TYPE_JOB:         {ACTION_DELETE: true},
		},
		model.USER_ROLE_VIEWER: {
			UNIT_TYPE_TEAM_MEMBER: {ACTION_DELETE: true},
			UNIT_TYPE_USER:        {ACTION_DELETE: true},
		},
	},
	ATTRIBUTE_CATEGORY_MANAGE: {
		model.USER_ROLE_ANONYMOUS: {
			UNIT_TYPE_APP: {ACTION_MANAGE_RUN_ACTION: true},
		},
		model.USER_ROLE_OWNER: {
			UNIT_TYPE_TEAM:              {ACTION_MANAGE_TEAM_NAME: true, ACTION_MANAGE_TEAM_ICON: true, ACTION_MANAGE_TEAM_CONFIG: true, ACTION_MANAGE_UPDATE_TEAM_DOMAIN: true},
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_MANAGE_REMOVE_MEMBER: true, ACTION_MANAGE_ROLE: true, ACTION_MANAGE_ROLE_FROM_OWNER: true, ACTION_MANAGE_ROLE_FROM_ADMIN: true, ACTION_MANAGE_ROLE_FROM_EDITOR: true, ACTION_MANAGE_ROLE_FROM_VIEWER: true, ACTION_MANAGE_ROLE_TO_OWNER: true, ACTION_MANAGE_ROLE_TO_ADMIN: true, ACTION_MANAGE_ROLE_TO_EDITOR: true, ACTION_MANAGE_ROLE_TO_VIEWER: true},
			UNIT_TYPE_USER:              {ACTION_MANAGE_RENAME_USER: true, ACTION_MANAGE_UPDATE_USER_AVATAR: true},
			UNIT_TYPE_INVITE:            {ACTION_MANAGE_CONFIG_INVITE: true, ACTION_MANAGE_INVITE_LINK: true},
			UNIT_TYPE_DOMAIN:            {ACTION_MANAGE_TEAM_DOMAIN: true, ACTION_MANAGE_APP_DOMAIN: true},
			UNIT_TYPE_BILLING:           {ACTION_MANAGE_PAYMENT_INFO: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_MANAGE_DASHBOARD_BROADCAST: true},
			UNIT_TYPE_APP:               {ACTION_MANAGE_CREATE_APP: true, ACTION_MANAGE_EDIT_APP: true},
			UNIT_TYPE_COMPONENTS:        {},
			UNIT_TYPE_RESOURCE:          {ACTION_MANAGE_CREATE_RESOURCE: true, ACTION_MANAGE_EDIT_RESOURCE: true},
			UNIT_TYPE_ACTION:            {ACTION_MANAGE_CREATE_ACTION: true, ACTION_MANAGE_EDIT_ACTION: true, ACTION_MANAGE_PREVIEW_ACTION: true, ACTION_MANAGE_RUN_ACTION: true},
			UNIT_TYPE_TRANSFORMER:       {},
			UNIT_TYPE_JOB:               {},
		},
		model.USER_ROLE_ADMIN: {
			UNIT_TYPE_TEAM:              {ACTION_MANAGE_TEAM_NAME: true, ACTION_MANAGE_TEAM_ICON: true, ACTION_MANAGE_UPDATE_TEAM_DOMAIN: true, ACTION_MANAGE_TEAM_CONFIG: true},
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_MANAGE_REMOVE_MEMBER: true, ACTION_MANAGE_ROLE: true, ACTION_MANAGE_ROLE_FROM_ADMIN: true, ACTION_MANAGE_ROLE_FROM_EDITOR: true, ACTION_MANAGE_ROLE_FROM_VIEWER: true, ACTION_MANAGE_ROLE_TO_ADMIN: true, ACTION_MANAGE_ROLE_TO_EDITOR: true, ACTION_MANAGE_ROLE_TO_VIEWER: true},
			UNIT_TYPE_USER:              {ACTION_MANAGE_RENAME_USER: true, ACTION_MANAGE_UPDATE_USER_AVATAR: true},
			UNIT_TYPE_INVITE:            {ACTION_MANAGE_CONFIG_INVITE: true, ACTION_MANAGE_INVITE_LINK: true},
			UNIT_TYPE_DOMAIN:            {ACTION_MANAGE_TEAM_DOMAIN: true, ACTION_MANAGE_APP_DOMAIN: true},
			UNIT_TYPE_BILLING:           {ACTION_MANAGE_PAYMENT_INFO: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_MANAGE_DASHBOARD_BROADCAST: true},
			UNIT_TYPE_APP:               {ACTION_MANAGE_CREATE_APP: true, ACTION_MANAGE_EDIT_APP: true},
			UNIT_TYPE_COMPONENTS:        {},
			UNIT_TYPE_RESOURCE:          {ACTION_MANAGE_CREATE_RESOURCE: true, ACTION_MANAGE_EDIT_RESOURCE: true},
			UNIT_TYPE_ACTION:            {ACTION_MANAGE_CREATE_ACTION: true, ACTION_MANAGE_EDIT_ACTION: true, ACTION_MANAGE_PREVIEW_ACTION: true, ACTION_MANAGE_RUN_ACTION: true},
			UNIT_TYPE_TRANSFORMER:       {},
			UNIT_TYPE_JOB:               {},
		},
		model.USER_ROLE_EDITOR: {
			UNIT_TYPE_TEAM_MEMBER:       {ACTION_MANAGE_REMOVE_MEMBER: true, ACTION_MANAGE_ROLE: true, ACTION_MANAGE_ROLE_FROM_EDITOR: true, ACTION_MANAGE_ROLE_FROM_VIEWER: true, ACTION_MANAGE_ROLE_TO_EDITOR: true, ACTION_MANAGE_ROLE_TO_VIEWER: true},
			UNIT_TYPE_USER:              {ACTION_MANAGE_RENAME_USER: true, ACTION_MANAGE_UPDATE_USER_AVATAR: true},
			UNIT_TYPE_BUILDER_DASHBOARD: {ACTION_MANAGE_DASHBOARD_BROADCAST: true},
			UNIT_TYPE_APP:               {ACTION_MANAGE_CREATE_APP: true, ACTION_MANAGE_EDIT_APP: true},
			UNIT_TYPE_RESOURCE:          {ACTION_MANAGE_CREATE_RESOURCE: true, ACTION_MANAGE_EDIT_RESOURCE: true},
			UNIT_TYPE_ACTION:            {ACTION_MANAGE_CREATE_ACTION: true, ACTION_MANAGE_EDIT_ACTION: true, ACTION_MANAGE_PREVIEW_ACTION: true, ACTION_MANAGE_RUN_ACTION: true},
			UNIT_TYPE_TRANSFORMER:       {},
			UNIT_TYPE_JOB:               {},
		},
		model.USER_ROLE_VIEWER: {
			UNIT_TYPE_TEAM_MEMBER: {ACTION_MANAGE_REMOVE_MEMBER: true, ACTION_MANAGE_ROLE: true, ACTION_MANAGE_ROLE_FROM_VIEWER: true, ACTION_MANAGE_ROLE_TO_VIEWER: true},
			UNIT_TYPE_USER:        {ACTION_MANAGE_RENAME_USER: true, ACTION_MANAGE_UPDATE_USER_AVATAR: true},
			UNIT_TYPE_ACTION:      {ACTION_MANAGE_RUN_ACTION: true},
			UNIT_TYPE_JOB:         {},
		},
	},
	ATTRIBUTE_CATEGORY_SPECIAL: {
		model.USER_ROLE_OWNER: {
			UNIT_TYPE_TEAM:        {ACTION_SPECIAL_EDITOR_AND_VIEWER_CAN_INVITE_BY_LINK_SW: true},
			UNIT_TYPE_TEAM_MEMBER: {ACTION_SPECIAL_TRANSFER_OWNER: true},
			UNIT_TYPE_INVITE:      {ACTION_SPECIAL_INVITE_LINK_RENEW: true},
			UNIT_TYPE_APP:         {ACTION_SPECIAL_RELEASE_APP: true},
		},
		model.USER_ROLE_ADMIN: {
			UNIT_TYPE_TEAM:   {ACTION_SPECIAL_EDITOR_AND_VIEWER_CAN_INVITE_BY_LINK_SW: true},
			UNIT_TYPE_INVITE: {ACTION_SPECIAL_INVITE_LINK_RENEW: true},
			UNIT_TYPE_APP:    {ACTION_SPECIAL_RELEASE_APP: true},
		},
		model.USER_ROLE_EDITOR: {
			UNIT_TYPE_APP: {ACTION_SPECIAL_RELEASE_APP: true},
		},
		model.USER_ROLE_VIEWER: {},
	},
}

type Attribute struct {
	Access  map[int]bool
	Delete  map[int]bool
	Manage  map[int]bool
	Special map[int]bool
}

func NewAttribute(userRole int, unitType int) *Attribute {
	attr := &Attribute{
		Access:  AttributeConfigList[ATTRIBUTE_CATEGORY_ACCESS][userRole][unitType],
		Delete:  AttributeConfigList[ATTRIBUTE_CATEGORY_DELETE][userRole][unitType],
		Manage:  AttributeConfigList[ATTRIBUTE_CATEGORY_MANAGE][userRole][unitType],
		Special: AttributeConfigList[ATTRIBUTE_CATEGORY_SPECIAL][userRole][unitType],
	}
	return attr
}

type AttributeGroup struct {
	UserRole  int
	UnitType  int
	UnitID    int
	Attribute *Attribute
}

func (attrg *AttributeGroup) SetUserRole(userRole int) {
	attrg.UserRole = userRole
}

func (attrg *AttributeGroup) SetUnitType(unitType int) {
	attrg.UnitType = unitType
}

func (attrg *AttributeGroup) SetUnitID(unitID int) {
	attrg.UnitID = unitID
}

func (attrg *AttributeGroup) CanAccess(attribute int) bool {
	r, match := attrg.Attribute.Access[attribute]
	if !match {
		return false
	}
	return r
}

func (attrg *AttributeGroup) CanDelete(attribute int) bool {
	r, match := attrg.Attribute.Delete[attribute]
	if !match {
		return false
	}
	return r
}

func (attrg *AttributeGroup) CanManage(attribute int) bool {
	r, match := attrg.Attribute.Manage[attribute]
	if !match {
		return false
	}
	return r
}

func (attrg *AttributeGroup) CanManageSpecial(attribute int) bool {
	r, match := attrg.Attribute.Special[attribute]
	if !match {
		return false
	}
	return r
}

func (attrg *AttributeGroup) CanModify(attribute, fromID, toID int) bool {
	// @todo: extend this method, now only support modify user role check.
	if attribute == ACTION_MANAGE_ROLE {
		return attrg.CanModifyRoleFromTo(fromID, toID)
	}
	return false
}

func (attrg *AttributeGroup) CanInvite(userRole int) bool {
	// convert to attribute
	attribute, hit := InviteRoleAttributeMap[userRole]
	if !hit {
		return false
	}
	// check attirbute
	r, match := attrg.Attribute.Access[attribute]
	if !match {
		return false
	}
	return r
}

func (attrg *AttributeGroup) CanModifyRoleFromTo(fromRole, toRole int) bool {
	// convert to attribute
	fromRoleAttribute, fromHit := ModifyRoleFromAttributeMap[fromRole]
	toRoleAttribute, toHit := MadifyRoleToAttributeMap[toRole]
	if !fromHit || !toHit {
		return false
	}
	// check attirbute
	fromResult, fromMatch := attrg.Attribute.Manage[fromRoleAttribute]
	toResult, toMatch := attrg.Attribute.Manage[toRoleAttribute]
	if !fromMatch || !toMatch {
		return false
	}
	return fromResult && toResult
}

func (attrg *AttributeGroup) DoesNowUserAreEditorOrViewer() bool {
	if attrg.UserRole == model.USER_ROLE_EDITOR || attrg.UserRole == model.USER_ROLE_VIEWER {
		return true
	}
	return false
}

func NewAttributeGroup(userRole int, unitType int) *AttributeGroup {
	attr := NewAttribute(userRole, unitType)
	attrg := &AttributeGroup{
		UserRole:  userRole,
		UnitType:  unitType,
		UnitID:    0, // 0 for placeholder, this feature has not implemented.
		Attribute: attr,
	}
	return attrg
}
