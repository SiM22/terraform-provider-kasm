package provider

import (
    "log"
    "os"
    "terraform-provider-kasm/internal/resources/user"
)

var debugMode = false

func init() {
    debugMode = os.Getenv("KASM_DEBUG") != ""
}

func debugLog(format string, v ...interface{}) {
    if debugMode {
        log.Printf("[DEBUG] "+format, v...)
    }
}

// For pretty printing API responses
func debugLogObject(prefix string, obj interface{}) {
    if debugMode {
        log.Printf("[DEBUG] %s: %+v", prefix, obj)
    }
}

func debugCompareStates(prefix string, old, new user.UserResourceModel) {
    if !debugMode {
        return
    }

    debugLog("%s - Comparing states:", prefix)
    debugLog("ID: %v -> %v", old.ID, new.ID)
    debugLog("Username: %v -> %v", old.Username, new.Username)
    debugLog("FirstName: %v -> %v", old.FirstName, new.FirstName)
    debugLog("LastName: %v -> %v", old.LastName, new.LastName)
    debugLog("Organization: %v -> %v", old.Organization, new.Organization)
    debugLog("Phone: %v -> %v", old.Phone, new.Phone)
    debugLog("Locked: %v -> %v", old.Locked, new.Locked)
    debugLog("Disabled: %v -> %v", old.Disabled, new.Disabled)
    debugLog("Groups: %v -> %v", old.Groups, new.Groups)
}

func debugLogState(prefix string, state user.UserResourceModel) {
    if !debugMode {
        return
    }
    debugLog("%s STATE ========================", prefix)
    debugLog("ID: %v", state.ID.ValueString())
    debugLog("Groups IsNull: %v", state.Groups.IsNull())
    debugLog("Groups IsUnknown: %v", state.Groups.IsUnknown())
    debugLog("Locked: %v (IsNull: %v)", state.Locked.ValueBool(), state.Locked.IsNull())
    debugLog("Disabled: %v (IsNull: %v)", state.Disabled.ValueBool(), state.Disabled.IsNull())
    debugLog("=====================================")
}

func debugLogPlan(prefix string, plan, state user.UserResourceModel) {
    if !debugMode {
        return
    }
    debugLog("%s PLAN =========================", prefix)
    debugLog("Plan Groups -> State Groups:")
    debugLog("IsNull: %v -> %v", plan.Groups.IsNull(), state.Groups.IsNull())
    debugLog("IsUnknown: %v -> %v", plan.Groups.IsUnknown(), state.Groups.IsUnknown())
    debugLog("Plan Locked -> State Locked:")
    debugLog("IsNull: %v -> %v", plan.Locked.IsNull(), state.Locked.IsNull())
    debugLog("Value: %v -> %v", plan.Locked.ValueBool(), state.Locked.ValueBool())
    debugLog("Plan Disabled -> State Disabled:")
    debugLog("IsNull: %v -> %v", plan.Disabled.IsNull(), state.Disabled.IsNull())
    debugLog("Value: %v -> %v", plan.Disabled.ValueBool(), state.Disabled.ValueBool())
    debugLog("=====================================")
}
